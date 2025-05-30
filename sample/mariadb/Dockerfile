FROM mariadb:11.4

# 환경 변수 설정 (초기화 스크립트 실행을 위해 필요)
ENV MYSQL_ROOT_PASSWORD=qwe123
ENV MYSQL_DATABASE=testdb

# 초기화 스크립트 복사
# /docker-entrypoint-initdb.d/ 디렉토리의 파일들은 컨테이너 초기화 시 자동 실행됩니다
COPY init.sql /docker-entrypoint-initdb.d/01-init.sql
COPY load_csv.sql /docker-entrypoint-initdb.d/02-load-csv.sql

# 포트 노출
EXPOSE 3306

VOLUME ["/var/lib/mysql"]

# 헬스체크: 실제 데이터베이스 연결 테스트 (연결 성공 = 로드 완료)
HEALTHCHECK --interval=30s --timeout=30s --start-period=1800s --retries=3 \
  CMD mysql -h localhost -u root -p${MYSQL_ROOT_PASSWORD} -e "SELECT 1;" > /dev/null 2>&1 || exit 1
