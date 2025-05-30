# 🐬 MariaDB 테스트 컨테이너

CSV 파일로부터 대용량 테스트 데이터를 빠르게 로드하는 MariaDB 컨테이너입니다.

## 📋 구성

### 테이블 구조

- **users**: 사용자 정보 (이름, 이메일, 나이)
- **categories**: 상품 카테고리 (10개 기본 카테고리)
- **products**: 상품 정보 (이름, 가격, 카테고리, 재고)
- **orders**: 주문 정보 (사용자-상품 관계, 수량, 상태, 날짜)

### 파일 구조

```
sample/mariadb/
├── Dockerfile              # MariaDB 컨테이너 정의
├── docker-compose.yml      # 컨테이너 실행 설정
├── init.sql               # 테이블 생성 스크립트
├── load_csv.sql           # CSV 데이터 로드 스크립트
└── README.md              # 이 파일
```

## 🚀 사용 방법

### 1. CSV 데이터 생성

먼저 `../common` 폴더에서 CSV 데이터를 생성합니다:

```bash
cd ../common

# uv 방식 (권장)
uv run generate_data.py

# 또는 pip 방식
pip install -r requirements.txt
python generate_data.py
```

### 2. MariaDB 컨테이너 실행

```bash
cd ../mariadb

# 컨테이너 빌드 및 시작
docker-compose up --build -d

# 로그 확인
docker-compose logs -f mariadb-test
```

### 3. 데이터 로드 완료 확인

```bash
# healthcheck 상태 확인
docker ps

# 상태가 "healthy"가 되면 로드 완료!
```

## 🔗 접속 정보

### 데이터베이스 연결

- **Host**: localhost
- **Port**: 3306
- **Database**: testdb
- **Username**: root
- **Password**: pwd

### 연결 예시

```bash
# 컨테이너 내부 접속
docker exec -it mariadb-test-container mysql -u root -ppwd

# 호스트에서 직접 접속
mysql -h localhost -P 3306 -u root -ppwd testdb
```

## 📊 데이터 확인

### 기본 통계 쿼리

```sql
-- 테이블별 데이터 개수
SELECT
    (SELECT COUNT(*) FROM users) as total_users,
    (SELECT COUNT(*) FROM products) as total_products,
    (SELECT COUNT(*) FROM orders) as total_orders,
    (SELECT COUNT(*) FROM categories) as total_categories;

-- 사용자별 주문 통계 (상위 10명)
SELECT
    u.name,
    COUNT(o.id) as order_count,
    SUM(o.total_amount) as total_spent
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id, u.name
ORDER BY total_spent DESC
LIMIT 10;

-- 카테고리별 주문 통계
SELECT
    c.name as category,
    COUNT(o.id) as order_count,
    AVG(p.price) as avg_price
FROM categories c
LEFT JOIN products p ON c.id = p.category_id
LEFT JOIN orders o ON p.id = o.product_id
GROUP BY c.id, c.name
ORDER BY order_count DESC;
```

## ⚙️ 설정 조정

### 환경변수

`docker-compose.yml`에서 설정 변경 가능:

```yaml
environment:
  - MYSQL_ROOT_PASSWORD=pwd # root 비밀번호
  - MYSQL_DATABASE=testdb # 기본 데이터베이스
```

### CSV 데이터 크기 조정

`../common/.env` 파일에서 데이터 개수 조정:

```bash
USERS_COUNT=1000000      # 사용자 수
PRODUCTS_COUNT=100000    # 상품 수
ORDERS_COUNT=1000000     # 주문 수
```

## 🔧 트러블슈팅

### 로드 실패시

```bash
# 컨테이너와 볼륨 완전 삭제 후 재시작
docker-compose down -v
docker-compose up --build -d
```

### CSV 파일 없음 에러

```bash
# common 폴더에서 CSV 생성 확인
ls -la ../common/csv_data/

# CSV 파일이 없으면 생성
cd ../common && uv run generate_data.py
```

### 권한 에러

```bash
# CSV 파일 권한 확인
chmod 644 ../common/csv_data/*.csv
```

### 메모리 부족

`init.sql`에서 더 작은 데이터로 테스트:

```bash
# common 폴더에서 소규모 데이터 생성
cd ../common
echo "USERS_COUNT=10000" > .env
echo "PRODUCTS_COUNT=1000" >> .env
echo "ORDERS_COUNT=50000" >> .env
uv run generate_data.py
```

## 📝 성능 최적화

### healthcheck 대기 시간

대용량 데이터의 경우 `start-period` 조정:

```yaml
healthcheck:
  start_period: 1800s # 30분 (기본값)
  # 더 큰 데이터: 3600s (1시간)
```

### 데이터베이스 설정

더 빠른 로드를 위해 `my.cnf` 추가 가능

## 🎯 활용 예시

- **성능 테스트**: 대용량 쿼리 성능 측정
- **애플리케이션 테스트**: 실제 데이터와 유사한 환경
- **학습 목적**: SQL 쿼리 연습
- **벤치마킹**: 다른 DB와 성능 비교
