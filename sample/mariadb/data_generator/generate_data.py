#!/usr/bin/env python3
import csv
import random
import uuid
from datetime import datetime, timedelta
import os
import shutil
from dotenv import load_dotenv

# .env 파일 로드
load_dotenv()

# 고정된 출력 경로
OUTPUT_DIR = './output/csv_data'

def get_env_int(key, default):
    """환경변수를 정수로 가져오기"""
    try:
        return int(os.getenv(key, default))
    except (ValueError, TypeError):
        print(f"⚠️  Warning: Invalid value for {key}, using default: {default:,}")
        return default

def generate_users_csv(filename, count):
    """사용자 CSV 생성"""
    print(f"Generating {count:,} users...")
    
    with open(filename, 'w', newline='', encoding='utf-8') as f:
        writer = csv.writer(f)
        
        for i in range(1, count + 1):
            name = f"User_{uuid.uuid4().hex[:8]}"
            email = f"user{i}@test.com"
            age = random.randint(18, 80)
            
            writer.writerow([name, email, age])
            
            if i % 100_000 == 0:
                print(f"  Generated {i:,} users...")
    
    print(f"✅ Users CSV generated: {filename}")

def generate_products_csv(filename, count):
    """상품 CSV 생성"""
    print(f"Generating {count:,} products...")
    
    product_names = [
        'iPhone', 'Galaxy', 'MacBook', 'iPad', 'ThinkPad',
        'Surface', 'AirPods', 'Watch', 'Monitor', 'Camera',
        'Laptop', 'Desktop', 'Tablet', 'Smartphone', 'Headphone'
    ]
    
    with open(filename, 'w', newline='', encoding='utf-8') as f:
        writer = csv.writer(f)
        
        for i in range(1, count + 1):
            name = f"{random.choice(product_names)}_{uuid.uuid4().hex[:4]}"
            price = round(random.uniform(10, 2000), 2)
            category = uuid.uuid4().hex[:6]
            stock = random.randint(1, 1000)
            description = f"Description for product {i}"
            
            writer.writerow([name, price, category, stock, description])
            
            if i % 50_000 == 0:
                print(f"  Generated {i:,} products...")
    
    print(f"✅ Products CSV generated: {filename}")

def generate_orders_csv(filename, count, max_user_id, max_product_id):
    """주문 CSV 생성"""
    print(f"Generating {count:,} orders...")
    
    with open(filename, 'w', newline='', encoding='utf-8') as f:
        writer = csv.writer(f)
        
        for i in range(1, count + 1):
            user_id = random.randint(1, max_user_id)
            product_id = random.randint(1, max_product_id)
            quantity = random.randint(1, 5)
            total_amount = round(random.uniform(10, 1000), 2)
            status = uuid.uuid4().hex[:6]
            
            # 랜덤 날짜 (지난 1년간)
            days_ago = random.randint(0, 365)
            order_date = (datetime.now() - timedelta(days=days_ago)).strftime('%Y-%m-%d %H:%M:%S')
            
            writer.writerow([user_id, product_id, quantity, total_amount, status, order_date])
            
            if i % 500_000 == 0:
                print(f"  Generated {i:,} orders...")
    
    print(f"✅ Orders CSV generated: {filename}")

def main():
    """메인 실행 함수"""
    print("🚀 Starting CSV generation for MariaDB...")
    
    # 환경변수에서 설정값 읽기 (기본값 제공)
    database_name = os.getenv('DATABASE_NAME', 'testdb')
    users_count = get_env_int('USERS_COUNT', 1_000_000)
    products_count = get_env_int('PRODUCTS_COUNT', 100_000)
    orders_count = get_env_int('ORDERS_COUNT', 1_000_000)
    orders_table_count = get_env_int('ORDERS_TABLE_COUNT', 9)
    
    print(f"📊 Generation settings:")
    print(f"   Database: {database_name}")
    print(f"   Users: {users_count:,}")
    print(f"   Products: {products_count:,}")
    print(f"   Orders: {orders_count:,}")
    print(f"   Orders files: {orders_table_count}")
    print(f"   Output directory: {OUTPUT_DIR}")
    print()
    
    start_time = datetime.now()
    
    # 출력 디렉토리 생성
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    
    # CSV 파일 생성 (하드코딩된 파일명)
    generate_users_csv(f'{OUTPUT_DIR}/users.csv', users_count)
    generate_products_csv(f'{OUTPUT_DIR}/products.csv', products_count)
    
    # 주문 파일 생성: 첫 번째만 생성하고 나머지는 복사
    if orders_table_count > 0:
        # 첫 번째 주문 파일 생성
        first_orders_file = f'{OUTPUT_DIR}/orders1.csv'
        generate_orders_csv(first_orders_file, orders_count, users_count, products_count)
        
        # 나머지 주문 파일들은 복사로 생성
        if orders_table_count > 1:
            print(f"📄 Copying orders1.csv to {orders_table_count - 1} additional files...")
            for i in range(2, orders_table_count + 1):
                target_file = f'{OUTPUT_DIR}/orders{i}.csv'
                shutil.copy2(first_orders_file, target_file)
                print(f"  ✅ Copied to orders{i}.csv")
    
    end_time = datetime.now()
    duration = end_time - start_time
    
    print(f"\n🎉 All CSV files generated!")
    print(f"⏱️  Total time: {duration}")
    print(f"📁 Files saved in: {OUTPUT_DIR}/")
    
    # 파일 크기 확인
    total_size = 0
    file_list = ['users.csv', 'products.csv']
    for i in range(1, orders_table_count + 1):
        file_list.append(f"orders{i}.csv")
    
    for filename in file_list:
        filepath = f'{OUTPUT_DIR}/{filename}'
        if os.path.exists(filepath):
            size_mb = os.path.getsize(filepath) / (1024 * 1024)
            total_size += size_mb
            print(f"   {filename}: {size_mb:.1f} MB")
    
    print(f"   Total: {total_size:.1f} MB")

if __name__ == "__main__":
    main() 