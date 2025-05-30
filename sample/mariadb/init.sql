-- 테스트 데이터베이스 생성
CREATE DATABASE IF NOT EXISTS testdb;
USE testdb;

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

-- 카테고리 테이블
CREATE TABLE categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 상품 테이블
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    category_id INT,
    stock INT DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id),
    INDEX idx_category_id (category_id),
    INDEX idx_price (price)
);

-- 주문 테이블
CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT DEFAULT 1,
    total_amount DECIMAL(10,2) NOT NULL,
    order_status ENUM('pending', 'confirmed', 'shipped', 'delivered', 'cancelled') DEFAULT 'pending',
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    INDEX idx_user_id (user_id),
    INDEX idx_product_id (product_id),
    INDEX idx_order_date (order_date),
    INDEX idx_status (order_status)
);

-- 기본 카테고리 데이터 삽입
INSERT INTO categories (name, description) VALUES
('전자제품', '전자제품 카테고리'),
('의류', '의류 관련 상품'),
('도서', '각종 도서 및 교재'),
('스포츠', '스포츠 용품'),
('가구', '가정용 가구'),
('식품', '식품 및 음료'),
('뷰티', '화장품 및 뷰티 용품'),
('완구', '장난감 및 완구'),
('자동차', '자동차 관련 용품'),
('건강', '건강 관련 제품');
