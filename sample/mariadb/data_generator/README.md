# 🔧 테스트 데이터 생성기

대용량 테스트 데이터를 빠르게 생성하는 Python 스크립트입니다.

## 📋 기능

- **사용자 데이터**: 이름, 이메일, 나이
- **상품 데이터**: 상품명, 가격, 카테고리, 재고
- **주문 데이터**: 사용자-상품 관계, 주문 상태, 날짜
- **환경변수 설정**: 생성할 데이터 개수를 유연하게 조정

## 🚀 빠른 시작

### 1. uv 설치 (한 번만)

```bash
# macOS/Linux
curl -LsSf https://astral.sh/uv/install.sh | sh

# Windows
powershell -c "irm https://astral.sh/uv/install.ps1 | iex"

# 또는 homebrew (macOS)
brew install uv
```

### 2. 환경 설정

```bash
# .env 파일 복사 및 수정
cp .env.example .env
```

### 3. 실행 (uv 방식)

```bash
# uv로 의존성 자동 관리하며 실행
uv run generate_data.py
```

## 🔄 기존 방식도 지원

### pip 방식

```bash
# 1. 의존성 설치
pip install -r requirements.txt

# 2. 실행
python generate_data.py
```

### uv 방식 (권장 ⚡️)

```bash
# 프로젝트 동기화 (한 번만)
uv sync

# 실행
uv run generate_data.py

# 또는 직접 실행
uv run python generate_data.py
```

## ⚙️ 환경변수 설정

`.env` 파일에서 데이터 생성 개수를 설정할 수 있습니다:

```bash
# 기본 설정
USERS_COUNT=1000000      # 사용자 100만명
PRODUCTS_COUNT=100000    # 상품 10만개
ORDERS_COUNT=1000000     # 주문 100만건
OUTPUT_DIR=csv_data      # 출력 디렉토리
```

## 🎯 프리셋 설정

### 소규모 테스트 (1분 이내)

```bash
USERS_COUNT=10000
PRODUCTS_COUNT=1000
ORDERS_COUNT=50000
```

### 중간 규모 (5-10분)

```bash
USERS_COUNT=100000
PRODUCTS_COUNT=10000
ORDERS_COUNT=500000
```

### 대규모 (30분-1시간)

```bash
USERS_COUNT=10000000
PRODUCTS_COUNT=1000000
ORDERS_COUNT=20000000
```

### 초대규모 (2-4시간)

```bash
USERS_COUNT=50000000
PRODUCTS_COUNT=5000000
ORDERS_COUNT=100000000
```

## 🏃‍♂️ uv 장점

### 왜 uv를 사용하나요?

- **⚡️ 속도**: pip보다 10-100배 빠른 의존성 설치
- **🔒 안정성**: 자동 잠금 파일 생성
- **🎯 간편함**: Python 버전 관리까지 통합
- **🚀 현대적**: 최신 Python 표준 준수

### 성능 비교

| 작업            | pip  | uv       |
| --------------- | ---- | -------- |
| 의존성 설치     | 30초 | **2초**  |
| 가상환경 생성   | 5초  | **1초**  |
| 프로젝트 초기화 | 수동 | **자동** |

## 📊 예상 파일 크기

| 규모     | Users | Products | Orders | 총 크기     |
| -------- | ----- | -------- | ------ | ----------- |
| 소규모   | ~1MB  | ~0.1MB   | ~5MB   | **~6MB**    |
| 중간     | ~10MB | ~1MB     | ~50MB  | **~61MB**   |
| 대규모   | ~1GB  | ~100MB   | ~2GB   | **~3.1GB**  |
| 초대규모 | ~5GB  | ~500MB   | ~10GB  | **~15.5GB** |

## 📁 출력 파일

생성되는 CSV 파일:

- `csv_data/users.csv`: 사용자 데이터
- `csv_data/products.csv`: 상품 데이터
- `csv_data/orders.csv`: 주문 데이터

## 🔄 데이터베이스 로드

### MariaDB/MySQL

```sql
LOAD DATA LOCAL INFILE 'csv_data/users.csv'
INTO TABLE users
FIELDS TERMINATED BY ','
LINES TERMINATED BY '\n'
(name, email, age);
```

### PostgreSQL

```sql
COPY users(name, email, age)
FROM 'csv_data/users.csv'
DELIMITER ',' CSV;
```

## ⏱️ 성능 팁

### 빠른 생성을 위해:

1. **SSD 사용**: HDD보다 3-5배 빠름
2. **충분한 메모리**: 최소 8GB 이상 권장
3. **uv 사용**: pip보다 훨씬 빠름
4. **적절한 배치 크기**: 환경변수로 조정 가능

### 데이터베이스 로드 최적화:

1. **인덱스 비활성화**: 로드 전 인덱스 삭제
2. **트랜잭션 설정**: autocommit=0, bulk insert
3. **동시성 조정**: 병렬 로드 고려

## 🛠️ 커스터마이징

### 데이터 형식 변경

`generate_data.py`에서 다음 함수들을 수정:

- `generate_users_csv()`: 사용자 데이터 형식
- `generate_products_csv()`: 상품 데이터 형식
- `generate_orders_csv()`: 주문 데이터 형식

### 출력 형식 변경

CSV 구분자나 인코딩 변경 가능:

```python
writer = csv.writer(f, delimiter='\t')  # 탭 구분
```

### 새 의존성 추가

```bash
# uv로 패키지 추가
uv add pandas  # 예시: pandas 추가

# 또는 pyproject.toml 직접 수정 후
uv sync
```

## 🔍 트러블슈팅

### uv 관련

```bash
# uv 캐시 클리어
uv cache clean

# 의존성 다시 설치
uv sync --reinstall
```

### 메모리 부족

```bash
# 배치 크기 줄이기
ORDERS_COUNT=1000000  # 작은 수로 시작
```

### 디스크 공간 부족

```bash
# 출력 디렉토리 변경
OUTPUT_DIR=/path/to/large/disk
```

### 권한 오류

```bash
# 출력 디렉토리 권한 확인
chmod 755 csv_data/
```

## 🔗 유용한 명령어

### uv 명령어

```bash
# 프로젝트 초기화
uv init

# 의존성 추가
uv add package-name

# 의존성 제거
uv remove package-name

# 가상환경에서 셸 실행
uv shell

# Python 버전 관리
uv python install 3.12
uv python pin 3.12
```

## 📝 라이선스

이 프로젝트는 테스트 목적으로 만들어졌습니다.
