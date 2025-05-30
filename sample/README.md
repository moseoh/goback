## 테스트 데이터 CSV 파일 준비

## Database Docker 실행 방법

```shell
# 1. 처음 시작하거나 Dockerfile/SQL 파일이 변경되었을 때
docker compose up --build -d
# 이후 실행
docker compose up -d

# 2. 로그 확인
docker compose logs -f mariadb-test

# 3. 완전히 정리하고 다시 시작하고 싶을 때
docker compose down -v
docker compose up --build -d

# 4. 개발 중 자주 사용하는 명령어 조합
docker compose down -v && docker compose up --build -d
```

## 백업

```shell
time mariadb-dump -h localhost -P 25001 -u root -ppwd testdb > test.sql
```
