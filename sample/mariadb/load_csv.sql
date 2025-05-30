USE testdb;

-- 성능 최적화 설정
SET SESSION sql_log_bin = 0;
SET SESSION autocommit = 0;
SET SESSION unique_checks = 0;
SET SESSION foreign_key_checks = 0;

-- users 로드
SELECT 'Loading users...' as status;
LOAD DATA LOCAL INFILE '/csv_data/users.csv'
INTO TABLE users
FIELDS TERMINATED BY ','
LINES TERMINATED BY '\n'
(name, email, age);
COMMIT;

-- products 로드  
SELECT 'Loading products...' as status;
LOAD DATA LOCAL INFILE '/csv_data/products.csv'
INTO TABLE products
FIELDS TERMINATED BY ','
LINES TERMINATED BY '\n'
(name, price, category_id, stock, description);
COMMIT;

-- orders 로드
SELECT 'Loading orders...' as status;
LOAD DATA LOCAL INFILE '/csv_data/orders.csv'
INTO TABLE orders
FIELDS TERMINATED BY ','
LINES TERMINATED BY '\n'
(user_id, product_id, quantity, total_amount, order_status, order_date);
COMMIT;

-- 설정 복원
SET SESSION sql_log_bin = 1;
SET SESSION autocommit = 1;
SET SESSION unique_checks = 1;
SET SESSION foreign_key_checks = 1;

-- 통계 확인
SELECT 
    (SELECT COUNT(*) FROM users) as total_users,
    (SELECT COUNT(*) FROM products) as total_products,
    (SELECT COUNT(*) FROM orders) as total_orders;

SELECT 'Data loading completed!' as status; 