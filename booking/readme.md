# 배타락을 이용한 이중예약 문제 테스트

이 프로젝트는 PostgreSQL과 Go(Gin 프레임워크)를 사용하여 **배타락(Exclusive Lock)** 기반으로 **좌석 이중예약(Double Booking)** 문제를 방지하는 테스트 서버입니다.

---

## 📌 테이블 모델

```sql
CREATE TABLE seats (
  id INT PRIMARY KEY,
  isbooked INT,
  name TEXT
);
```

### ✅ 초기 데이터 삽입

```sql
INSERT INTO seats (id, isbooked, name)
SELECT generate_series(1, 15), 0, '';
```

---

## 🌐 웹 접속 주소

애플리케이션 실행 후 웹 브라우저에서 아래 주소로 접속하세요:

```
http://0.0.0.0:8080/
```

---

## ⚙️ CORS 설정 변경 (main.go)

`main.go` 파일에서 **CORS 정책에 허용할 origins URL을 반드시 수정**해야 합니다.  
예를 들어 프론트엔드가 다른 포트나 도메인에서 동작할 경우:

```go
r.Use(cors.New(cors.Config{
  AllowOrigins: []string{"http://localhost:3000"}, // 여기를 수정
  ...
}))
```

---

## 🧪 주요 기능

- `/seats` (GET): 전체 좌석 조회
- `/:id/:name` (PUT): 좌석 예약 → 이미 예약된 좌석은 배타락으로 막음

---

## 🔒 배타락 기반 동작

- 예약 시 `SELECT ... FOR UPDATE`로 배타락을 걸어 다른 트랜잭션의 예약 충돌 방지
- 동일한 좌석에 대해 동시에 PUT 요청 시, 한 트랜잭션만 성공

---

## ✅ 실행 요약

1. PostgreSQL 실행 및 테이블 생성
2. `main.go` 실행
3. 웹 페이지 접속: [http://0.0.0.0:8080/](http://0.0.0.0:8080/)
4. 좌석 예약 테스트 (이중예약 방지 확인)
