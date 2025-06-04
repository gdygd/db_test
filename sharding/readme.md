# PostgreSQL Shard 컨테이너 설정

- PostgreSQL 컨테이너를 3개 생성하고 초기화 시 `init.sql`을 통해 테이블을 자동 생성 합니다. 
- shard를 테스트하기 위한 프로제긑 입니다.

---

## 📄 1. init.sql 생성

```bash
vim init.sql
```

```sql
-- script for making table
CREATE TABLE URL_TABLE (
    id SERIAL NOT NULL PRIMARY KEY,
    URL TEXT,
    URL_ID CHARACTER(5)
);
```

---

## 📄 2. Dockerfile 생성

```bash
vim Dockerfile
```

```dockerfile
FROM postgres
COPY init.sql /docker-entrypoint-initdb.d
```

---

## 🛠️ 3. Docker 이미지 빌드

```bash
docker build -t pgshard .
```

---

## 🚀 4. 컨테이너 실행 (포트: 5433 ~ 5435)

```bash
docker run --name pgshard1 -p 5433:5432 -e POSTGRES_PASSWORD=secret -d pgshard
docker run --name pgshard2 -p 5434:5432 -e POSTGRES_PASSWORD=secret -d pgshard
docker run --name pgshard3 -p 5435:5432 -e POSTGRES_PASSWORD=secret -d pgshard
```

---

## ✅ 5. 컨테이너 내부 테이블 확인

```bash
docker exec -it pgshard1 psql -U postgres -c '\dt'
```

또는 테이블 조회:

```bash
docker exec -it pgshard1 psql -U postgres -c 'SELECT * FROM URL_TABLE;'
```

---

## 📌 참고사항

- `POSTGRES_PASSWORD`는 필수이며, 지정하지 않으면 컨테이너가 종료됩니다.
- `init.sql`의 위치는 `Dockerfile`과 같은 디렉토리에 있어야 합니다.
- 컨테이너가 실행되지 않으면 `docker logs <컨테이너명>` 명령어로 원인 확인이 가능합니다.
