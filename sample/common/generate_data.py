#!/usr/bin/env python3
import csv
import random
import uuid
from datetime import datetime, timedelta
import os
from dotenv import load_dotenv

# .env 파일 로드
load_dotenv()

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
            category_id = random.randint(1, 10)
            stock = random.randint(1, 1000)
            description = f"Description for product {i}"
            
            writer.writerow([name, price, category_id, stock, description])
            
            if i % 50_000 == 0:
                print(f"  Generated {i:,} products...")
    
    print(f"✅ Products CSV generated: {filename}")

def generate_orders_csv(filename, count, max_user_id, max_product_id):
    """주문 CSV 생성"""
    print(f"Generating {count:,} orders...")
    
    statuses = ['pending', 'confirmed', 'shipped', 'delivered', 'cancelled']
    
    with open(filename, 'w', newline='', encoding='utf-8') as f:
        writer = csv.writer(f)
        
        for i in range(1, count + 1):
            user_id = random.randint(1, max_user_id)
            product_id = random.randint(1, max_product_id)
            quantity = random.randint(1, 5)
            total_amount = round(random.uniform(10, 1000), 2)
            status = random.choice(statuses)
            
            # 랜덤 날짜 (지난 1년간)
            days_ago = random.randint(0, 365)
            order_date = (datetime.now() - timedelta(days=days_ago)).strftime('%Y-%m-%d %H:%M:%S')
            
            writer.writerow([user_id, product_id, quantity, total_amount, status, order_date])
            
            if i % 500_000 == 0:
                print(f"  Generated {i:,} orders...")
    
    print(f"✅ Orders CSV generated: {filename}")

def main():
    """메인 실행 함수"""
    print("🚀 Starting CSV generation...")
    
    # 환경변수에서 데이터 개수 읽기
    users_count = get_env_int('USERS_COUNT', 1_000_000)
    products_count = get_env_int('PRODUCTS_COUNT', 100_000)
    orders_count = get_env_int('ORDERS_COUNT', 1_000_000)
    output_dir = os.getenv('OUTPUT_DIR', 'csv_data')
    
    print(f"📊 Generation settings:")
    print(f"   Users: {users_count:,}")
    print(f"   Products: {products_count:,}")
    print(f"   Orders: {orders_count:,}")
    print(f"   Output directory: {output_dir}")
    print()
    
    start_time = datetime.now()
    
    # 출력 디렉토리 생성
    os.makedirs(output_dir, exist_ok=True)
    
    # CSV 파일 생성
    generate_users_csv(f'{output_dir}/users.csv', users_count)
    generate_products_csv(f'{output_dir}/products.csv', products_count)
    generate_orders_csv(f'{output_dir}/orders.csv', orders_count, users_count, products_count)
    
    end_time = datetime.now()
    duration = end_time - start_time
    
    print(f"\n🎉 All CSV files generated!")
    print(f"⏱️  Total time: {duration}")
    print(f"📁 Files saved in: {output_dir}/")
    
    # 파일 크기 확인
    total_size = 0
    for filename in ['users.csv', 'products.csv', 'orders.csv']:
        filepath = f'{output_dir}/{filename}'
        if os.path.exists(filepath):
            size_mb = os.path.getsize(filepath) / (1024 * 1024)
            total_size += size_mb
            print(f"   {filename}: {size_mb:.1f} MB")
    
    print(f"   Total: {total_size:.1f} MB")

if __name__ == "__main__":
    main() 