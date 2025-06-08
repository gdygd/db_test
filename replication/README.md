# PostgreSQL 17 Master/Standby Replication (Docker 환경)

이 문서는 Docker를 이용하여 PostgreSQL 17에서 **Master/Standby 복제 환경**을 구성하고 테스트하는 과정을 안내합니다. 기본 계정은 `postgres`, 비밀번호는 `secret`을 사용합니다.

---

## 1. Docker로 DB 인스턴스 생성

```bash
# 마스터 DB 컨테이너 생성
docker run --name pmaster -p 5433:5432 \
  -v /home/gildong/work/db/replication/pmaster_data:/var/lib/postgresql/data \
  -e POSTGRES_PASSWORD=secret \
  -d postgres:17-alpine

# 스탠바이 DB 컨테이너 생성
docker run --name pstandby -p 5434:5432 \
  -v /home/gildong/work/db/replication/pstandby_data:/var/lib/postgresql/data \
  -e POSTGRES_PASSWORD=secret \
  -d postgres:17-alpine
```

---

## 2. 마스터 데이터 디렉터리 복사

```bash
# 기존 standby 디렉터리 백업 및 초기화
mv pstandby_data pstandby_data_bk
sudo cp -R pmaster_data pstandby_data
```

---

## 3. Master 설정

### `pg_hba.conf` 수정

```conf
host replication postgres all md5
```

- `postgres` 유저가 replication 접속할 수 있도록 허용
- 모든 IP에서 접근 허용 (`all` → 필요시 IP 범위 제한 가능)

---

## 4. Standby 설정

### `postgresql.conf` 수정

```conf
primary_conninfo = 'application_name=standby1 host=xxx.xxx.xxx.xxx port=5433 user=postgres password=secret'
```

- `application_name=standby1` 은 Primary에서 인식할 standby 이름

---

## 5. standby.signal 파일 생성

```bash
cd pstandby_data
touch standby.signal
```

- 이 파일이 존재하면 PostgreSQL은 Standby 모드로 실행됨

---

## 6. Master `postgresql.conf` 추가 설정

```conf
synchronous_standby_names = 'first 1 (standby1)'
```

- `standby1`이 반드시 WAL 수신해야 Primary가 커밋 진행
- 동기 복제 활성화

---

## 7. 복제 상태 확인

Primary 컨테이너에 접속 후 확인:

```sql
SELECT * FROM pg_stat_replication;
```

- `application_name = standby1` 이 표시되어야 함
- 복제 상태는 `streaming` 이 되어야 정상

---

## 복제 테스트

- 마스터에서 테이블 생성:

```sql
CREATE TABLE test_tbl (id serial PRIMARY KEY, name text);
```

- 스탠바이에서는 해당 테이블이 **자동 복제되어 존재**하지만, **쓰기 작업은 불가능** (`read-only` 상태)

---

## 참고 사항

- 스탠바이 서버는 `read-only` 상태로 실행됨 (쓰기 불가)
- 마스터에서 스탠바이로의 WAL 복제를 통해 데이터 동기화
- `standby.signal` 파일은 PostgreSQL 12 이상 버전에서 사용되는 스탠바이 작동 신호

---




> 작성자: gildong  
> 📅 작성일: 2025-06  
