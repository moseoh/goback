# MariaDB 테스트 환경

MariaDB를 위한 대용량 테스트 데이터 생성 및 데이터베이스 구축 통합 솔루션입니다.

## 📁 프로젝트 구조

```
mariadb/
├── run.sh                    # 🚀 통합 실행 스크립트
├── docker-compose.yml        # Docker 컨테이너 설정
├── Dockerfile               # MariaDB 이미지 설정
├── init.sql                 # 데이터베이스 초기화 (자동 생성됨)
├── load_csv.sql            # CSV 데이터 로드 (자동 생성됨)
└── data_generator/         # 데이터 생성 도구
    ├── .env                # 🔧 데이터 생성 설정 (핵심!)
    ├── config.py           # 공통 설정 파일
    ├── generate_data.py    # CSV 데이터 생성
    ├── generate_sql.py     # SQL 스크립트 생성  
    ├── pyproject.toml      # Python 프로젝트 설정
    └── csv_data/           # 생성된 CSV 파일들 (자동 생성됨)
```

## 🚀 빠른 시작

### 원클릭 실행

```bash
./run.sh
```

이 스크립트가 자동으로 수행하는 작업:

1. **Python 환경 설정**: 가상환경 생성 및 의존성 설치
2. **SQL 파일 생성**: `.env` 설정 기반으로 `init.sql`, `load_csv.sql` 자동 생성
3. **CSV 데이터 생성**: 대용량 테스트 데이터 생성
4. **MariaDB 실행**: Docker 컨테이너 빌드 및 실행

### 수동 실행

각 단계를 개별적으로 실행하려면:

```bash
# 1. 데이터 생성
cd data_generator
python3 -m venv .venv
source .venv/bin/activate
pip install python-dotenv

# 2. SQL 파일 생성
python generate_sql.py

# 3. CSV 데이터 생성  
python generate_data.py
deactivate

# 4. MariaDB 실행
cd ..
docker-compose up --build -d
```

## ⚙️ 설정 커스터마이징

### .env 파일로 데이터 양 조정 (핵심!)

`data_generator/.env` 파일에서 모든 데이터 생성 설정을 관리합니다:

```bash
# 테스트 데이터 생성 설정
# 각 테이블별 생성할 데이터 개수를 설정합니다.

# 사용자 수 (기본값: 1,000,000)
USERS_COUNT=50000

# 상품 수 (기본값: 100,000)
PRODUCTS_COUNT=100000

# 주문 수 (기본값: 1,000,000)
ORDERS_TABLE_COUNT=10    # 주문 파일 분할 개수
ORDERS_COUNT=1000000     # 각 파일당 주문 수

# 출력 디렉토리 (기본값: csv_data)
OUTPUT_DIR=csv_data
```

### 추가 설정 (config.py)

테이블 이름, 주문 상태 등은 `data_generator/config.py`에서 관리:

```python
# 테이블 이름 변경
TABLES = {
    "users": "users",
    "products": "products", 
    "orders": "orders",
    "categories": "categories"
}

# 주문 상태 커스터마이징
ORDER_STATUSES = ['pending', 'confirmed', 'shipped', 'delivered', 'cancelled']
```

### 설정 적용

설정 변경 후 다시 실행하려면:

```bash
./run.sh
```

## 📊 데이터베이스 연결

실행 완료 후 다음 정보로 연결:

- **Host**: localhost  
- **Port**: 25001
- **Database**: testdb
- **Username**: root
- **Password**: qwe123

### 연결 예시

```bash
# MySQL CLI
mysql -h localhost -P 25001 -u root -pqwe123 testdb

# CSV 데이터 로드
docker exec -i mariadb-test-container mysql -u root -pqwe123 testdb < load_csv.sql
```

## 📈 생성되는 데이터

### 기본 설정 기준 (.env 파일):

- **users**: 50,000 레코드
- **products**: 100,000 레코드  
- **orders**: 10,000,000 레코드 (10개 파일로 분할, 각 파일당 1,000,000 레코드)
- **categories**: 10 레코드 (전자제품, 의류, 도서 등)

### 테이블 구조:

#### users
- id, name, email, age, created_at, updated_at

#### products  
- id, name, price, category_id, stock, description, created_at

#### orders
- id, user_id, product_id, quantity, total_amount, order_status, order_date

#### categories
- id, name, description, created_at

## 🔧 고급 설정

### 환경변수로 빠른 조정

`.env` 파일을 수정하여 즉시 적용:

```bash
# data_generator/.env
USERS_COUNT=500000           # 50만 사용자
PRODUCTS_COUNT=50000         # 5만 상품
ORDERS_TABLE_COUNT=5         # 5개 파일로 분할  
ORDERS_COUNT=200000          # 각 파일당 20만 주문
OUTPUT_DIR=custom_csv_data   # 커스텀 출력 폴더
```

### Docker 설정 변경

`docker-compose.yml`에서 포트나 비밀번호 변경:

```yaml
ports:
  - "25001:3306"  # 호스트 포트 변경
environment:
  - MYSQL_ROOT_PASSWORD=qwe123  # 비밀번호 변경
```

## 🗂️ 파일 크기 예상

현재 `.env` 설정 기준:
- **users.csv**: ~5MB (50,000 레코드)
- **products.csv**: ~20MB (100,000 레코드)
- **orders1-10.csv**: ~900MB (총 1,000만 레코드, 각 파일 ~90MB)
- **총 용량**: ~925MB

## 🚦 문제 해결

### 컨테이너 재시작
```bash
docker-compose down
docker-compose up --build -d
```

### 로그 확인
```bash
docker-compose logs mariadb-test
```

### 데이터 초기화
```bash
docker-compose down -v  # 볼륨까지 삭제
./run.sh                # 처음부터 다시 실행
```

### 메모리 부족 시
`.env` 파일에서 데이터 양을 줄이세요:

```bash
# 메모리 절약 설정
USERS_COUNT=10000
PRODUCTS_COUNT=5000  
ORDERS_TABLE_COUNT=3
ORDERS_COUNT=50000
```

## 🏗️ 확장 가능한 구조

이 프로젝트는 설정 중심 설계로 쉽게 확장 가능합니다:

1. **데이터 양 조정**: `.env` 파일만 수정
2. **새 테이블 추가**: `config.py`에 테이블 정의 추가
3. **새 데이터 타입**: `generate_data.py`에 생성 함수 추가
4. **새 SQL 스크립트**: `generate_sql.py`에 생성 로직 추가

모든 설정이 파일 기반으로 관리되므로 버전 관리와 팀 협업이 용이합니다! 🎯
