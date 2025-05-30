#!/usr/bin/env python3
"""
공통 설정을 기반으로 SQL 스크립트 생성
"""

import os
from dotenv import load_dotenv

# .env 파일 로드
load_dotenv()

# 고정된 출력 경로
SQL_OUTPUT_DIR = './output/sql'
CSV_MOUNT_PATH = '/csv_data'

def get_env_int(key, default):
    """환경변수를 정수로 가져오기"""
    try:
        return int(os.getenv(key, default))
    except (ValueError, TypeError):
        print(f"⚠️  Warning: Invalid value for {key}, using default: {default}")
        return default

def generate_init_sql():
    """init.sql 생성"""
    database_name = os.getenv('DATABASE_NAME', 'testdb')
    orders_table_count = get_env_int('ORDERS_TABLE_COUNT', 9)
    
    sql_content = f"""-- 테스트 데이터베이스 생성
CREATE DATABASE IF NOT EXISTS {database_name};
USE {database_name};

-- 사용자 테이블
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    age INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_created_at (created_at)
);

-- 상품 테이블
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    category VARCHAR(100),
    stock INT DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_category (category),
    INDEX idx_price (price)
);
"""

    # 각 orders 테이블을 개별적으로 생성
    for i in range(1, orders_table_count + 1):
        sql_content += f"""
-- 주문 테이블 {i}
CREATE TABLE orders{i} (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT DEFAULT 1,
    total_amount DECIMAL(10,2) NOT NULL,
    order_status VARCHAR(100),
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    INDEX idx_user_id (user_id),
    INDEX idx_product_id (product_id),
    INDEX idx_order_date (order_date),
    INDEX idx_status (order_status)
);"""
    
    return sql_content

def generate_load_csv_sql():
    """load_csv.sql 생성"""
    # 환경변수에서 설정값 읽기
    database_name = os.getenv('DATABASE_NAME', 'testdb')
    orders_table_count = get_env_int('ORDERS_TABLE_COUNT', 9)
    
    sql_content = f"""USE {database_name};

-- 성능 최적화 설정
SET SESSION sql_log_bin = 0;
SET SESSION autocommit = 0;
SET SESSION unique_checks = 0;
SET SESSION foreign_key_checks = 0;

-- users 로드
SELECT 'Loading users...' as status;
LOAD DATA LOCAL INFILE '{CSV_MOUNT_PATH}/users.csv'
INTO TABLE users
FIELDS TERMINATED BY ','
LINES TERMINATED BY '\\n'
(name, email, age);
COMMIT;

-- products 로드  
SELECT 'Loading products...' as status;
LOAD DATA LOCAL INFILE '{CSV_MOUNT_PATH}/products.csv'
INTO TABLE products
FIELDS TERMINATED BY ','
LINES TERMINATED BY '\\n'
(name, price, category, stock, description);
COMMIT;
"""

    # 각 orders 파일을 각각의 테이블에 로드
    for i in range(1, orders_table_count + 1):
        sql_content += f"""
-- orders{i} 로드
SELECT 'Loading orders{i} from orders{i}.csv...' as status;
LOAD DATA LOCAL INFILE '{CSV_MOUNT_PATH}/orders{i}.csv'
INTO TABLE orders{i}
FIELDS TERMINATED BY ','
LINES TERMINATED BY '\\n'
(user_id, product_id, quantity, total_amount, order_status, order_date);
COMMIT;"""

    sql_content += f"""

-- 설정 복원
SET SESSION sql_log_bin = 1;
SET SESSION autocommit = 1;
SET SESSION unique_checks = 1;
SET SESSION foreign_key_checks = 1;

-- 통계 확인
SELECT 
    (SELECT COUNT(*) FROM users) as total_users,
    (SELECT COUNT(*) FROM products) as total_products"""

    # 각 orders 테이블의 통계 추가
    for i in range(1, orders_table_count + 1):
        sql_content += f""",
    (SELECT COUNT(*) FROM orders{i}) as total_orders{i}"""

    sql_content += """;

SELECT 'Data loading completed!' as status;
"""
    
    return sql_content

def main():
    """SQL 파일 생성"""
    print("🔧 Generating SQL files...")
    
    # SQL 출력 디렉토리 생성
    os.makedirs(SQL_OUTPUT_DIR, exist_ok=True)
    
    # 01-init.sql 생성 (테이블 생성이 먼저)
    init_sql = generate_init_sql()
    with open(f'{SQL_OUTPUT_DIR}/01-init.sql', 'w', encoding='utf-8') as f:
        f.write(init_sql)
    print(f"✅ Generated {SQL_OUTPUT_DIR}/01-init.sql")
    
    # 02-load_csv.sql 생성 (데이터 로드가 나중)
    load_sql = generate_load_csv_sql()
    with open(f'{SQL_OUTPUT_DIR}/02-load_csv.sql', 'w', encoding='utf-8') as f:
        f.write(load_sql)
    print(f"✅ Generated {SQL_OUTPUT_DIR}/02-load_csv.sql")
    
    print("🎉 All SQL files generated successfully!")

if __name__ == "__main__":
    main() 