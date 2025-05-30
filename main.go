package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type BackupConfig struct {
	Host        string
	Port        string
	Username    string
	Password    string
	Database    string
	OutputDir   string
	Workers     int // 병렬 워커 수
	BatchSize   int // 배치 처리 크기
	MultiInsert int // 멀티 INSERT 문의 최대 행 수
}

type MySQLBackup struct {
	config *BackupConfig
	db     *sql.DB
}

type TableBackupResult struct {
	TableName string
	Error     error
	Index     int    // 원래 순서 보존용
	RowCount  int64  // 백업된 행 수
	TempFile  string // 임시 파일 경로
}

type TableInfo struct {
	Name             string
	EstimatedRows    int64
	IsLargeTable     bool
	OptimalMethod    string
	OrderColumn      string
	OrderColumnType  string
	HasAutoIncrement bool
	HasTimestamp     bool
}

func NewMySQLBackup(config *BackupConfig) *MySQLBackup {
	return &MySQLBackup{
		config: config,
	}
}

func (mb *MySQLBackup) Connect() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		mb.config.Username, mb.config.Password, mb.config.Host, mb.config.Port, mb.config.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("데이터베이스 연결 실패: %v", err)
	}

	// 연결 설정 (병렬 처리를 위해 연결 수 증가)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(mb.config.Workers * 2) // 워커 수의 2배로 설정
	db.SetMaxIdleConns(mb.config.Workers)

	// 연결 테스트
	if err := db.Ping(); err != nil {
		return fmt.Errorf("데이터베이스 핑 실패: %v", err)
	}

	mb.db = db
	fmt.Printf("✓ 데이터베이스 '%s'에 성공적으로 연결되었습니다.\n", mb.config.Database)
	return nil
}

func (mb *MySQLBackup) GetTables() ([]string, error) {
	query := "SHOW TABLES"
	rows, err := mb.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("테이블 목록 조회 실패: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("테이블 이름 스캔 실패: %v", err)
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

func (mb *MySQLBackup) analyzeTable(tableName string) (*TableInfo, error) {
	info := &TableInfo{Name: tableName}

	// 1. 테이블 크기 추정 (INFORMATION_SCHEMA 사용)
	sizeQuery := `
		SELECT COALESCE(TABLE_ROWS, 0) 
		FROM INFORMATION_SCHEMA.TABLES 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?`

	err := mb.db.QueryRow(sizeQuery, mb.config.Database, tableName).Scan(&info.EstimatedRows)
	if err != nil {
		info.EstimatedRows = 0 // 추정 실패시 0으로 설정
	}

	info.IsLargeTable = info.EstimatedRows > 10000

	// 2. 최적의 순서 컬럼 찾기 (우선순위: AUTO_INCREMENT > TIMESTAMP > 순차적 PK)
	orderColumn, columnType, method := mb.findBestOrderColumn(tableName)
	info.OrderColumn = orderColumn
	info.OrderColumnType = columnType
	info.OptimalMethod = method

	// 3. 특수 컬럼 존재 여부 확인
	info.HasAutoIncrement = strings.Contains(columnType, "auto_increment")
	info.HasTimestamp = strings.Contains(strings.ToLower(columnType), "timestamp") ||
		strings.Contains(strings.ToLower(columnType), "datetime")

	return info, nil
}

func (mb *MySQLBackup) findBestOrderColumn(tableName string) (string, string, string) {
	// 1순위: AUTO_INCREMENT 컬럼 찾기
	autoIncQuery := `
		SELECT COLUMN_NAME, COLUMN_TYPE
		FROM INFORMATION_SCHEMA.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? 
		AND EXTRA LIKE '%auto_increment%'
		LIMIT 1`

	var columnName, columnType string
	err := mb.db.QueryRow(autoIncQuery, mb.config.Database, tableName).Scan(&columnName, &columnType)
	if err == nil {
		return columnName, columnType, "auto_increment_cursor"
	}

	// 2순위: 정수형 Primary Key (UUID가 아닌)
	pkQuery := `
		SELECT c.COLUMN_NAME, c.COLUMN_TYPE
		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE k
		JOIN INFORMATION_SCHEMA.COLUMNS c ON k.COLUMN_NAME = c.COLUMN_NAME 
		WHERE k.TABLE_SCHEMA = ? AND k.TABLE_NAME = ? 
		AND k.CONSTRAINT_NAME = 'PRIMARY'
		AND c.TABLE_SCHEMA = ? AND c.TABLE_NAME = ?
		AND c.DATA_TYPE IN ('int', 'bigint', 'smallint', 'tinyint', 'mediumint')
		ORDER BY k.ORDINAL_POSITION
		LIMIT 1`

	err = mb.db.QueryRow(pkQuery, mb.config.Database, tableName, mb.config.Database, tableName).Scan(&columnName, &columnType)
	if err == nil {
		return columnName, columnType, "integer_pk_cursor"
	}

	// 3순위: TIMESTAMP/DATETIME 컬럼 (created_at, updated_at 등)
	timestampQuery := `
		SELECT COLUMN_NAME, COLUMN_TYPE
		FROM INFORMATION_SCHEMA.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? 
		AND (DATA_TYPE IN ('timestamp', 'datetime') 
		     OR COLUMN_NAME IN ('created_at', 'updated_at', 'date_created', 'date_modified'))
		ORDER BY 
			CASE 
				WHEN COLUMN_NAME = 'created_at' THEN 1
				WHEN COLUMN_NAME = 'date_created' THEN 2
				WHEN COLUMN_NAME = 'updated_at' THEN 3
				WHEN DATA_TYPE = 'timestamp' THEN 4
				WHEN DATA_TYPE = 'datetime' THEN 5
				ELSE 6
			END
		LIMIT 1`

	err = mb.db.QueryRow(timestampQuery, mb.config.Database, tableName).Scan(&columnName, &columnType)
	if err == nil {
		return columnName, columnType, "timestamp_cursor"
	}

	// 4순위: ROWID 사용 (MySQL 8.0+, InnoDB 테이블)
	// MySQL의 숨겨진 ROWID 활용
	return "_rowid", "bigint", "rowid_cursor"
}

func (mb *MySQLBackup) BackupTable(tableName string) (string, int64, error) {
	var sqlContent strings.Builder

	// 테이블 구조 백업
	createTableSQL, err := mb.getCreateTableSQL(tableName)
	if err != nil {
		return "", 0, fmt.Errorf("테이블 구조 조회 실패: %v", err)
	}

	sqlContent.WriteString(fmt.Sprintf("-- 테이블 %s 구조\n", tableName))
	sqlContent.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS `%s`;\n", tableName))
	sqlContent.WriteString(createTableSQL + ";\n\n")

	// 테이블 분석
	tableInfo, err := mb.analyzeTable(tableName)
	if err != nil {
		return "", 0, fmt.Errorf("테이블 분석 실패: %v", err)
	}

	// 최적 방법으로 데이터 백업
	var dataSQL string
	var rowCount int64

	if !tableInfo.IsLargeTable {
		// 소용량: 단순한 방법이 가장 빠름
		dataSQL, rowCount, err = mb.getTableDataSimple(tableName)
	} else {
		// 대용량: 최적 방법 선택
		switch tableInfo.OptimalMethod {
		case "auto_increment_cursor", "integer_pk_cursor":
			dataSQL, rowCount, err = mb.getTableDataCursorBased(tableName, tableInfo.OrderColumn, "순차 커서")
		case "timestamp_cursor":
			dataSQL, rowCount, err = mb.getTableDataCursorBased(tableName, tableInfo.OrderColumn, "시간 커서")
		case "rowid_cursor":
			dataSQL, rowCount, err = mb.getTableDataRowIdBased(tableName)
		default:
			dataSQL, rowCount, err = mb.getTableDataStreaming(tableName)
		}
	}

	if err != nil {
		return "", 0, fmt.Errorf("테이블 데이터 조회 실패: %v", err)
	}

	if dataSQL != "" {
		sqlContent.WriteString(fmt.Sprintf("-- 테이블 %s 데이터 (%d 행, %s)\n",
			tableName, rowCount, tableInfo.OptimalMethod))
		sqlContent.WriteString(dataSQL + "\n")
	}

	return sqlContent.String(), rowCount, nil
}

func (mb *MySQLBackup) getCreateTableSQL(tableName string) (string, error) {
	query := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
	var table, createSQL string
	err := mb.db.QueryRow(query).Scan(&table, &createSQL)
	if err != nil {
		return "", err
	}
	return createSQL, nil
}

// 소용량 테이블: 기존 방식 (단순하고 빠름)
func (mb *MySQLBackup) getTableDataSimple(tableName string) (string, int64, error) {
	query := fmt.Sprintf("SELECT * FROM `%s`", tableName)
	rows, err := mb.db.Query(query)
	if err != nil {
		return "", 0, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", 0, err
	}

	var insertStatements []string
	var rowCount int64

	columnNames := make([]string, len(columns))
	for i, col := range columns {
		columnNames[i] = fmt.Sprintf("`%s`", col)
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return "", 0, err
		}

		var valueStrings []string
		for _, value := range values {
			if value == nil {
				valueStrings = append(valueStrings, "NULL")
			} else {
				switch v := value.(type) {
				case []byte:
					escaped := strings.ReplaceAll(string(v), "'", "\\'")
					escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
					valueStrings = append(valueStrings, fmt.Sprintf("'%s'", escaped))
				case string:
					escaped := strings.ReplaceAll(v, "'", "\\'")
					escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
					valueStrings = append(valueStrings, fmt.Sprintf("'%s'", escaped))
				case time.Time:
					valueStrings = append(valueStrings, fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05")))
				default:
					valueStrings = append(valueStrings, fmt.Sprintf("'%v'", v))
				}
			}
		}

		insertSQL := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s);",
			tableName,
			strings.Join(columnNames, ", "),
			strings.Join(valueStrings, ", "))
		insertStatements = append(insertStatements, insertSQL)
		rowCount++
	}

	return strings.Join(insertStatements, "\n"), rowCount, nil
}

// 커서 기반 페이징 (AUTO_INCREMENT, 정수 PK, TIMESTAMP 등)
func (mb *MySQLBackup) getTableDataCursorBased(tableName, orderColumn, method string) (string, int64, error) {
	var allInserts []string
	var rowCount int64
	var lastValue interface{}

	for {
		var query string
		var rows *sql.Rows
		var err error

		if lastValue == nil {
			// 첫 번째 배치
			query = fmt.Sprintf("SELECT * FROM `%s` ORDER BY `%s` LIMIT %d",
				tableName, orderColumn, mb.config.BatchSize)
			rows, err = mb.db.Query(query)
		} else {
			// 다음 배치들
			query = fmt.Sprintf("SELECT * FROM `%s` WHERE `%s` > ? ORDER BY `%s` LIMIT %d",
				tableName, orderColumn, orderColumn, mb.config.BatchSize)
			rows, err = mb.db.Query(query, lastValue)
		}

		if err != nil {
			return "", 0, err
		}

		columns, err := rows.Columns()
		if err != nil {
			rows.Close()
			return "", 0, err
		}

		// 순서 컬럼의 인덱스 찾기
		orderIndex := -1
		for i, col := range columns {
			if col == orderColumn {
				orderIndex = i
				break
			}
		}

		batchInserts, batchCount, newLastValue, err := mb.processCursorRows(rows, tableName, columns, orderIndex)
		rows.Close()

		if err != nil {
			return "", 0, err
		}

		if batchCount == 0 {
			break // 더 이상 데이터가 없음
		}

		allInserts = append(allInserts, batchInserts...)
		rowCount += batchCount
		lastValue = newLastValue

		if batchCount < int64(mb.config.BatchSize) {
			break // 마지막 배치
		}
	}

	return strings.Join(allInserts, "\n"), rowCount, nil
}

// ROWID 기반 처리 (MySQL 8.0+)
func (mb *MySQLBackup) getTableDataRowIdBased(tableName string) (string, int64, error) {
	// 로깅 제거
	// fmt.Printf("   📊 테이블 '%s': ROWID 방식으로 처리\n", tableName)

	// ROWID가 지원되는지 확인
	testQuery := fmt.Sprintf("SELECT _rowid FROM `%s` LIMIT 1", tableName)
	_, err := mb.db.Query(testQuery)
	if err != nil {
		// ROWID 지원하지 않으면 스트리밍으로 폴백
		// fmt.Printf("   ⚠️ ROWID 미지원, 스트리밍 방식으로 전환\n")
		return mb.getTableDataStreaming(tableName)
	}

	return mb.getTableDataCursorBased(tableName, "_rowid", "ROWID 커서")
}

// 대용량 테이블 스트리밍 (최후의 수단)
func (mb *MySQLBackup) getTableDataStreaming(tableName string) (string, int64, error) {
	// 로깅 제거
	// fmt.Printf("   📊 테이블 '%s': 스트리밍 방식으로 처리\n", tableName)

	query := fmt.Sprintf("SELECT * FROM `%s`", tableName)
	rows, err := mb.db.Query(query)
	if err != nil {
		return "", 0, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", 0, err
	}

	var allInserts []string
	var rowCount int64

	columnNames := make([]string, len(columns))
	for i, col := range columns {
		columnNames[i] = fmt.Sprintf("`%s`", col)
	}

	var currentBatch []string

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return "", 0, err
		}

		var valueStrings []string
		for _, value := range values {
			if value == nil {
				valueStrings = append(valueStrings, "NULL")
			} else {
				switch v := value.(type) {
				case []byte:
					escaped := strings.ReplaceAll(string(v), "'", "\\'")
					escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
					valueStrings = append(valueStrings, fmt.Sprintf("'%s'", escaped))
				case string:
					escaped := strings.ReplaceAll(v, "'", "\\'")
					escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
					valueStrings = append(valueStrings, fmt.Sprintf("'%s'", escaped))
				case time.Time:
					valueStrings = append(valueStrings, fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05")))
				default:
					valueStrings = append(valueStrings, fmt.Sprintf("'%v'", v))
				}
			}
		}

		valueGroup := fmt.Sprintf("(%s)", strings.Join(valueStrings, ", "))
		currentBatch = append(currentBatch, valueGroup)
		rowCount++

		// 배치가 찼으면 INSERT 문 생성
		if len(currentBatch) >= mb.config.MultiInsert {
			insertSQL := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s;",
				tableName,
				strings.Join(columnNames, ", "),
				strings.Join(currentBatch, ", "))
			allInserts = append(allInserts, insertSQL)
			currentBatch = currentBatch[:0] // 슬라이스 재사용
		}
	}

	// 남은 배치 처리
	if len(currentBatch) > 0 {
		insertSQL := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s;",
			tableName,
			strings.Join(columnNames, ", "),
			strings.Join(currentBatch, ", "))
		allInserts = append(allInserts, insertSQL)
	}

	return strings.Join(allInserts, "\n"), rowCount, nil
}

func (mb *MySQLBackup) processCursorRows(rows *sql.Rows, tableName string, columns []string, orderIndex int) ([]string, int64, interface{}, error) {
	var insertStatements []string
	var rowCount int64
	var lastValue interface{}

	columnNames := make([]string, len(columns))
	for i, col := range columns {
		columnNames[i] = fmt.Sprintf("`%s`", col)
	}

	var currentBatch []string

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, 0, nil, err
		}

		// 순서 컬럼 값 저장
		if orderIndex >= 0 {
			lastValue = values[orderIndex]
		}

		var valueStrings []string
		for _, value := range values {
			if value == nil {
				valueStrings = append(valueStrings, "NULL")
			} else {
				switch v := value.(type) {
				case []byte:
					escaped := strings.ReplaceAll(string(v), "'", "\\'")
					escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
					valueStrings = append(valueStrings, fmt.Sprintf("'%s'", escaped))
				case string:
					escaped := strings.ReplaceAll(v, "'", "\\'")
					escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
					valueStrings = append(valueStrings, fmt.Sprintf("'%s'", escaped))
				case time.Time:
					valueStrings = append(valueStrings, fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05")))
				default:
					valueStrings = append(valueStrings, fmt.Sprintf("'%v'", v))
				}
			}
		}

		valueGroup := fmt.Sprintf("(%s)", strings.Join(valueStrings, ", "))
		currentBatch = append(currentBatch, valueGroup)
		rowCount++

		// 배치가 찼으면 INSERT 문 생성
		if len(currentBatch) >= mb.config.MultiInsert {
			insertSQL := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s;",
				tableName,
				strings.Join(columnNames, ", "),
				strings.Join(currentBatch, ", "))
			insertStatements = append(insertStatements, insertSQL)
			currentBatch = currentBatch[:0]
		}
	}

	// 남은 배치 처리
	if len(currentBatch) > 0 {
		insertSQL := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s;",
			tableName,
			strings.Join(columnNames, ", "),
			strings.Join(currentBatch, ", "))
		insertStatements = append(insertStatements, insertSQL)
	}

	return insertStatements, rowCount, lastValue, nil
}

func (mb *MySQLBackup) backupTableWorker(tableName string, index int, resultChan chan<- TableBackupResult, progressChan chan<- string) {
	start := time.Now()
	progressChan <- fmt.Sprintf("🔄 테이블 '%s' 백업 시작...", tableName)

	sql, rowCount, err := mb.BackupTable(tableName)
	duration := time.Since(start)

	if err != nil {
		progressChan <- fmt.Sprintf("❌ 테이블 '%s' 백업 실패 (%.2fs): %v", tableName, duration.Seconds(), err)
	} else {
		progressChan <- fmt.Sprintf("✓ 테이블 '%s' 백업 완료 (%.2fs, %d행)", tableName, duration.Seconds(), rowCount)
	}

	resultChan <- TableBackupResult{
		TableName: tableName,
		Error:     err,
		Index:     index,
		RowCount:  rowCount,
		TempFile:  sql, // SQL 내용을 TempFile 필드에 임시 저장
	}
}

func (mb *MySQLBackup) BackupDatabase() error {
	start := time.Now()

	// 출력 디렉토리 생성
	if err := os.MkdirAll(mb.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("출력 디렉토리 생성 실패: %v", err)
	}

	// 파일명 생성 (타임스탬프 포함)
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_backup_%s.sql", mb.config.Database, timestamp)
	filepath := filepath.Join(mb.config.OutputDir, filename)

	// 백업 파일 생성
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("백업 파일 생성 실패: %v", err)
	}
	defer file.Close()

	// 버퍼링된 writer 사용 (성능 향상)
	writer := bufio.NewWriterSize(file, 1024*1024) // 1MB 버퍼
	defer writer.Flush()

	// 헤더 작성
	header := fmt.Sprintf(`-- MySQL 데이터베이스 백업 (적응형 지능 최적화)
-- 데이터베이스: %s
-- 생성 시간: %s
-- 호스트: %s:%s
-- 워커 수: %d
-- 배치 크기: %d
-- 멀티 INSERT 크기: %d

SET FOREIGN_KEY_CHECKS=0;
SET SQL_MODE="NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

`, mb.config.Database, time.Now().Format("2006-01-02 15:04:05"),
		mb.config.Host, mb.config.Port, mb.config.Workers, mb.config.BatchSize, mb.config.MultiInsert)

	if _, err := writer.WriteString(header); err != nil {
		return fmt.Errorf("헤더 작성 실패: %v", err)
	}

	// 테이블 목록 조회
	tables, err := mb.GetTables()
	if err != nil {
		return err
	}

	// 실제 사용될 워커 수 (테이블 수와 설정된 워커 수 중 작은 값)
	actualWorkers := mb.config.Workers
	if len(tables) < actualWorkers {
		actualWorkers = len(tables)
	}

	fmt.Printf("📋 총 %d개의 테이블을 %d개 워커로 병렬 백업합니다.\n", len(tables), actualWorkers)

	// 채널 생성
	resultChan := make(chan TableBackupResult, len(tables))
	progressChan := make(chan string, len(tables)*2)

	// 워크그룹 생성
	var wg sync.WaitGroup

	// 진행상황 출력 고루틴
	go func() {
		for msg := range progressChan {
			fmt.Println(msg)
		}
	}()

	// 워커 풀을 사용하여 테이블 백업 (고루틴 수 제한)
	semaphore := make(chan struct{}, actualWorkers)

	for i, tableName := range tables {
		wg.Add(1)
		go func(tableName string, index int) {
			defer wg.Done()
			semaphore <- struct{}{} // 워커 슬롯 획득
			mb.backupTableWorker(tableName, index, resultChan, progressChan)
			<-semaphore // 워커 슬롯 반환
		}(tableName, i)
	}

	// 모든 워커 완료 대기
	go func() {
		wg.Wait()
		close(resultChan)
		close(progressChan)
	}()

	// 결과 수집 (원래 순서 보존)
	results := make([]TableBackupResult, len(tables))
	completedCount := 0
	failedCount := 0
	totalRows := int64(0)

	for result := range resultChan {
		results[result.Index] = result
		if result.Error != nil {
			failedCount++
		} else {
			completedCount++
			totalRows += result.RowCount
		}
	}

	fmt.Printf("\n📊 백업 완료 통계:\n")
	fmt.Printf("   - 성공: %d개\n", completedCount)
	fmt.Printf("   - 실패: %d개\n", failedCount)
	fmt.Printf("   - 총 행 수: %d행\n", totalRows)
	fmt.Printf("   - 총 소요시간: %.2fs\n\n", time.Since(start).Seconds())

	// 임시 파일들을 순서대로 합치기
	fmt.Printf("📄 임시 파일들을 합치는 중...\n")
	for i, result := range results {
		if result.Error != nil {
			log.Printf("⚠️ 테이블 '%s' 백업 실패: %v", result.TableName, result.Error)
			continue
		}

		// 최종 파일에 쓰기
		if _, err := writer.WriteString(result.TempFile + "\n"); err != nil {
			return fmt.Errorf("최종 파일 쓰기 실패: %v", err)
		}

		// 파일 합치기 진행상황 출력
		if (i+1)%10 == 0 || i == len(results)-1 {
			fmt.Printf("📄 [%d/%d] 임시 파일 합치기 완료\n", i+1, len(results))
		}
	}

	// 푸터 작성
	footer := "\nSET FOREIGN_KEY_CHECKS=1;\n"
	if _, err := writer.WriteString(footer); err != nil {
		return fmt.Errorf("푸터 작성 실패: %v", err)
	}

	totalDuration := time.Since(start)
	fmt.Printf("🎉 백업이 완료되었습니다: %s\n", filepath)
	fmt.Printf("⚡ 총 처리 시간: %.2fs (평균 %.2fs/테이블, %.0f행/초)\n",
		totalDuration.Seconds(),
		totalDuration.Seconds()/float64(len(tables)),
		float64(totalRows)/totalDuration.Seconds())

	return nil
}

func (mb *MySQLBackup) Close() {
	if mb.db != nil {
		mb.db.Close()
		fmt.Println("✓ 데이터베이스 연결이 종료되었습니다.")
	}
}

func main() {
	fmt.Println("🗃️  MySQL 적응형 지능 백업 도구 시작")
	fmt.Println("========================================")

	// 환경변수에서 설정 읽기 (우선순위: 환경변수 > 기본값)
	config := LoadConfigFromEnv()

	// 워커 수 설정 (기본값: CPU 코어 수)
	if config.Workers == 0 {
		config.Workers = runtime.NumCPU()
	}

	// 명령행 인수로 설정 덮어쓰기 (우선순위: 명령행 > 환경변수 > 기본값)
	if len(os.Args) > 1 {
		config.Database = os.Args[1]
	}
	if len(os.Args) > 2 {
		config.Host = os.Args[2]
	}
	if len(os.Args) > 3 {
		config.Username = os.Args[3]
	}

	fmt.Printf("🔧 설정 정보:\n")
	fmt.Printf("   - 호스트: %s:%s\n", config.Host, config.Port)
	fmt.Printf("   - 사용자: %s\n", config.Username)
	fmt.Printf("   - 데이터베이스: %s\n", config.Database)
	fmt.Printf("   - 출력 경로: %s\n", config.OutputDir)
	fmt.Printf("   - 병렬 워커 수: %d\n", config.Workers)
	fmt.Printf("   - 배치 크기: %d\n", config.BatchSize)
	fmt.Printf("   - 멀티 INSERT 크기: %d\n", config.MultiInsert)
	fmt.Println()

	backup := NewMySQLBackup(config)

	// 데이터베이스 연결
	if err := backup.Connect(); err != nil {
		log.Fatal(err)
	}
	defer backup.Close()

	// 백업 실행
	if err := backup.BackupDatabase(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("✨ 모든 작업이 완료되었습니다!")
}
