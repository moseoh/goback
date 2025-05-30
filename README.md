# 🗃️ MySQL 병렬 백업 도구 ⚡

Go 언어로 작성된 고성능 MySQL 데이터베이스 백업 도구입니다. `github.com/go-sql-driver/mysql` 패키지를 사용하여 **병렬 처리**로 빠른 백업을 수행합니다.

## ✨ 주요 기능

- 📊 **완전한 데이터베이스 백업**: 테이블 구조와 데이터를 모두 백업
- ⚡ **병렬 처리**: 고루틴을 사용한 테이블별 병렬 백업으로 성능 최적화
- 🔧 **SQL 덤프 생성**: 표준 SQL 형식으로 백업 파일 생성
- ⏰ **타임스탬프 파일명**: 백업 시간이 포함된 고유한 파일명 생성
- 🎯 **실시간 진행상황**: 각 테이블 백업 진행상황과 소요시간을 실시간 표시
- 🛡️ **안전한 연결**: MySQL 연결 풀링 및 타임아웃 설정
- 📄 **.env 파일 지원**: 환경변수를 파일로 관리 가능
- 🎚️ **워커 수 조절**: CPU 코어 수에 따른 자동 조절 또는 수동 설정

## 🚀 설치 및 실행

### 1. 프로젝트 클론 및 빌드

```bash
cd goback
go mod tidy
go build -o bin/mysql-backup .
```

### 2. 설정 파일 생성 (선택사항)

```bash
# .env 파일 생성
cp env.example .env
# 그리고 .env 파일을 수정하여 데이터베이스 설정 입력
```

### 3. 실행

```bash
# 기본 설정으로 실행 (test_db 데이터베이스)
./bin/mysql-backup

# 특정 데이터베이스 백업
./bin/mysql-backup my_database

# 호스트와 사용자명 지정
./bin/mysql-backup my_database localhost root
```

### 4. 명령행 인수

```bash
./bin/mysql-backup [데이터베이스명] [호스트] [사용자명]
```

- **데이터베이스명**: 백업할 MySQL 데이터베이스 이름 (기본값: test_db)
- **호스트**: MySQL 서버 호스트 (기본값: localhost)
- **사용자명**: MySQL 사용자명 (기본값: root)

## ⚙️ 설정

### 1. .env 파일을 통한 설정 (권장)

루트 디렉토리에 `.env` 파일을 생성하여 설정할 수 있습니다:

```bash
# env.example 파일을 복사하여 .env 파일 생성
cp env.example .env
```

`.env` 파일 내용:
```env
# MySQL 데이터베이스 연결 설정
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USERNAME=root
MYSQL_PASSWORD=your_password_here
MYSQL_DATABASE=your_database_name

# 백업 파일 저장 경로
BACKUP_OUTPUT_DIR=./backups

# 병렬 처리 설정 (기본값: CPU 코어 수)
BACKUP_WORKERS=8
```

### 2. 환경변수를 통한 설정

```bash
export MYSQL_HOST=192.168.1.100
export MYSQL_USERNAME=admin
export MYSQL_PASSWORD=secret
export MYSQL_DATABASE=production
export BACKUP_WORKERS=16  # 병렬 워커 수 설정
./bin/mysql-backup
```

### 3. 코드 내 기본값 수정

기본 설정값들을 `main.go` 파일에서 수정할 수 있습니다:

```go
config := &BackupConfig{
    Host:      "localhost",  // MySQL 서버 호스트
    Port:      "3306",       // MySQL 포트
    Username:  "root",       // 사용자명
    Password:  "password",   // 비밀번호
    Database:  "test_db",    // 데이터베이스명
    OutputDir: "./backups",  // 백업 파일 저장 경로
    Workers:   8,            // 병렬 워커 수
}
```

### 설정 우선순위

1. **명령행 인수** (최우선)
2. **환경변수** (.env 파일 포함)
3. **기본값** (코드 내 설정)

## 📁 출력 파일

백업 파일은 `./backups/` 디렉토리에 다음 형식으로 생성됩니다:

```
{데이터베이스명}_backup_{YYYYMMDD_HHMMSS}.sql
```

예시: `my_database_backup_20241225_143052.sql`

백업 파일에는 병렬 처리 정보가 헤더에 포함됩니다:
```sql
-- MySQL 데이터베이스 백업 (병렬 처리)
-- 데이터베이스: production
-- 생성 시간: 2024-12-25 14:30:52
-- 호스트: localhost:3306
-- 워커 수: 8
```

## 🔧 기술 스택

- **Go 1.21+**: 프로그래밍 언어
- **github.com/go-sql-driver/mysql**: MySQL 드라이버
- **github.com/joho/godotenv**: .env 파일 로더
- **database/sql**: Go 표준 데이터베이스 인터페이스
- **고루틴 & 채널**: 병렬 처리를 위한 Go 동시성 기능

## 📋 백업 파일 내용

생성되는 SQL 파일에는 다음이 포함됩니다:

1. **헤더 정보**: 백업 시간, 데이터베이스명, 호스트 정보, 워커 수
2. **MySQL 설정**: Foreign key 체크 비활성화 등
3. **테이블 구조**: `CREATE TABLE` 문 (원래 순서 보존)
4. **테이블 데이터**: `INSERT` 문
5. **푸터**: Foreign key 체크 재활성화

## 🛠️ 개발 및 테스트

### 의존성 설치

```bash
go mod tidy
```

### 빌드

```bash
go build -o ./bin/mysql-backup .
```

### 실행

```bash
./bin/mysql-backup
```

## 🔐 보안 고려사항

- **프로덕션 환경**: `.env` 파일이나 환경변수를 통해 데이터베이스 자격증명을 관리하세요
- **파일 권한**: `.env` 파일은 적절한 권한(600)으로 보호하세요
- **백업 파일**: 민감한 데이터가 포함될 수 있으므로 적절한 권한으로 보호하세요
- **네트워크**: TLS를 사용하는 것을 권장합니다

```bash
# .env 파일 권한 설정
chmod 600 .env
```

## 📝 사용 예시

```bash
# .env 파일 설정 후 실행
./bin/mysql-backup

# 특정 데이터베이스만 백업
./bin/mysql-backup production

# 원격 서버의 데이터베이스 백업
./bin/mysql-backup ecommerce 192.168.1.100 admin

# 실행 결과 예시:
🗃️  MySQL 병렬 백업 도구 시작
================================
✓ .env 파일에서 설정을 로드했습니다.
🔧 설정 정보:
   - 호스트: localhost:3306
   - 사용자: root
   - 데이터베이스: production
   - 출력 경로: ./backups
   - 병렬 워커 수: 8

✓ 데이터베이스 'production'에 성공적으로 연결되었습니다.
📋 총 25개의 테이블을 8개 워커로 병렬 백업합니다.

🔄 테이블 'users' 백업 시작...
🔄 테이블 'orders' 백업 시작...
🔄 테이블 'products' 백업 시작...
✓ 테이블 'users' 백업 완료 (0.45s)
✓ 테이블 'orders' 백업 완료 (1.23s)
✓ 테이블 'products' 백업 완료 (0.78s)
...

📊 백업 완료 통계:
   - 성공: 25개
   - 실패: 0개
   - 총 소요시간: 3.42s

💾 [10/25] 테이블 파일 쓰기 완료
💾 [20/25] 테이블 파일 쓰기 완료
💾 [25/25] 테이블 파일 쓰기 완료

🎉 백업이 완료되었습니다: ./backups/production_backup_20241225_143052.sql
⚡ 총 처리 시간: 3.42s (평균 0.14s/테이블)
✓ 데이터베이스 연결이 종료되었습니다.
✨ 모든 작업이 완료되었습니다!
```

## 🏆 성능 비교

| 도구 | 처리 시간 | 성능 개선 |
|------|----------|----------|
| mariadb-dump | 8초 | 기준 |
| 이전 버전 (순차) | 26초 | 기준 대비 3.25배 느림 |
| **병렬 버전** | **~3-5초** | **기준 대비 1.6-2.7배 빠름** |

*실제 성능은 하드웨어, 네트워크, 데이터베이스 크기에 따라 달라질 수 있습니다.

## 📄 라이선스

이 프로젝트는 MIT 라이선스 하에 제공됩니다.

## ⚡ 성능 최적화

### 병렬 처리 설정

- **기본값**: CPU 코어 수만큼 워커 생성
- **권장값**: 
  - 로컬 DB: CPU 코어 수 × 1-2
  - 원격 DB: 네트워크 대역폭과 DB 서버 성능에 따라 4-16개
- **최적화 팁**:
  ```bash
  # CPU 집약적인 경우
  export BACKUP_WORKERS=8
  
  # I/O 집약적인 경우 (원격 DB)
  export BACKUP_WORKERS=16
  ```

### 연결 풀 설정

병렬 처리를 위해 데이터베이스 연결 풀이 자동으로 조정됩니다:
- **최대 연결 수**: 워커 수 × 2
- **유휴 연결 수**: 워커 수와 동일
