# Tugas 04 — Clean Architecture

Implementasi Modul 4 untuk mata kuliah **Pemrograman Backend Lanjut (SIP375)**.

Modul ini merupakan pengembangan dari REST API Students pada Modul 3. Fokus utama bukan menambah fitur HTTP baru, tetapi **merestrukturisasi aplikasi ke dalam pendekatan Clean Architecture** dengan tetap mempertahankan perilaku API yang sudah ada.

Perubahan utama meliputi pemisahan business rules, repository, controller/service, helper, middleware, routing, konfigurasi aplikasi, database, serta penambahan structured logging dan unit testing.

---

## Tujuan

Modul ini bertujuan untuk:

* Menerapkan prinsip **Clean Architecture** pada REST API yang sudah ada.
* Menerapkan **Dependency Rule**, yaitu dependensi diarahkan menuju bagian inti aplikasi.
* Memisahkan business rules dari framework dan infrastruktur.
* Membuat business rules dapat diuji tanpa menjalankan server atau database.
* Menambahkan structured logging untuk monitoring request.
* Mempertahankan perilaku HTTP dari Modul 3 setelah restrukturisasi.
* Meningkatkan maintainability dan testability aplikasi.

---

## Fitur Utama

* REST API Students berbasis Fiber v2.
* PostgreSQL sebagai persistent storage.
* Repository Pattern.
* Clean Architecture dengan beberapa penyederhanaan untuk kebutuhan pembelajaran.
* Business rules terpisah dari Fiber.
* Unit testing menggunakan package `testing`.
* Structured logging menggunakan `log/slog`.
* Log rotation menggunakan `lumberjack`.
* Request ID untuk setiap HTTP request.
* Request logging dalam format JSON.
* CORS dan Helmet middleware.
* Global error handler.
* Context timeout untuk operasi database.
* Graceful shutdown.
* Validasi `Content-Type: application/json`.
* Pagination, search, sorting, dan filtering.
* Error translation dari repository ke HTTP response.
* Health check dengan pemeriksaan koneksi database.

---

## Tech Stack

| Komponen        | Teknologi                               |
| --------------- | --------------------------------------- |
| Language        | Go 1.27.0                               |
| Framework       | Fiber v2                                |
| Database        | PostgreSQL 15+                          |
| Database Driver | pgx/v5                                  |
| Connection Pool | pgxpool                                 |
| Environment     | godotenv                                |
| Logger          | log/slog                                |
| Log Rotation    | lumberjack                              |
| Testing         | Go `testing`                            |
| API Testing     | cURL                                    |
| Architecture    | Clean Architecture + Repository Pattern |

---

## Struktur Proyek

```text
tugas04/
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── main.go
│
├── app/
│   ├── model/
│   │   └── student.go
│   │
│   ├── repository/
│   │   └── student_repository.go
│   │
│   └── service/
│       ├── student_rules.go
│       ├── student_rules_test.go
│       └── student_service.go
│
├── config/
│   ├── env.go
│   ├── logger.go
│   └── app.go
│
├── database/
│   └── postgres.go
│
├── helper/
│   ├── response.go
│   └── request.go
│
├── middleware/
│   └── middleware.go
│
├── route/
│   └── route.go
│
├── migrations/
│   └── 001_create_students.sql
│
└── logs/
    └── app.log
```

Folder `logs/` dibuat secara otomatis oleh logger dan tidak di-commit ke repository.

---

## Arsitektur

Implementasi Modul 4 memetakan struktur proyek ke empat kelompok utama Clean Architecture.

### 1. Entities

Lokasi:

```text
app/model/
```

Berisi struktur data inti aplikasi, antara lain:

* `Student`
* `CreateStudentRequest`
* `ReplaceStudentRequest`
* `PatchStudentRequest`
* `WebResponse`
* `Meta`
* `ListQuery`

Bagian ini tidak bergantung pada framework HTTP maupun database.

---

### 2. Use Cases / Business Rules

Lokasi:

```text
app/service/student_rules.go
```

Berisi fungsi-fungsi bisnis murni:

* `ValidateCreate`
* `ValidateReplace`
* `ApplyPatch`
* `IsEmptyPatch`
* `CountTotalPages`

Fungsi-fungsi tersebut tidak menerima `fiber.Ctx` dan tidak menjalankan query database sehingga dapat diuji secara langsung menggunakan unit test.

---

### 3. Interface Adapters

Bagian ini terdiri dari:

```text
app/service/
app/repository/
helper/
```

#### Service

`student_service.go` bertindak sebagai controller sekaligus use-case coordinator.

Service:

* menerima `fiber.Ctx`;
* membaca request;
* menjalankan business rules;
* memanggil repository;
* menerjemahkan error;
* menentukan HTTP response.

#### Repository

`student_repository.go` bertanggung jawab terhadap akses PostgreSQL.

Service tidak perlu mengetahui detail query SQL karena komunikasi dilakukan melalui interface repository.

#### Helper

`helper/` menangani kebutuhan umum HTTP seperti:

* response presenter;
* parsing query parameter;
* parsing parameter ID;
* context timeout.

---

### 4. Frameworks & Drivers

Bagian terluar aplikasi terdiri dari:

```text
config/
database/
middleware/
route/
main.go
```

Komponen ini menangani:

* Fiber;
* PostgreSQL connection pool;
* environment configuration;
* middleware;
* routing;
* logger;
* application assembly;
* graceful shutdown.

---

## Dependency Rule

Arah dependensi dirancang agar bagian inti aplikasi tidak bergantung pada framework atau infrastruktur.

Gambaran sederhananya:

```text
main.go
├── config
├── database
├── app/repository
└── app/service
        ├── app/model
        ├── app/repository
        └── helper
              └── app/model

route
├── app/service
├── helper
└── middleware
```

Prinsip yang dipertahankan:

* `app/model` tidak bergantung pada package proyek lain.
* `app/service/student_rules.go` tidak bergantung pada Fiber.
* Repository menangani akses database.
* Routing tidak menangani business logic.
* `main.go` bertanggung jawab terhadap perakitan aplikasi.
* Framework dan infrastruktur berada di bagian luar.

---

## Penyederhanaan Clean Architecture

Implementasi ini tidak menggunakan Clean Architecture secara sepenuhnya kanonik karena disesuaikan dengan skala proyek dan kebutuhan pembelajaran.

### Controller dan Use Case

Pada implementasi ini, controller dan use-case coordinator berada dalam:

```text
app/service/student_service.go
```

Konsekuensinya, service masih mengenal `fiber.Ctx`.

Namun business rules yang benar-benar membutuhkan pengujian dipisahkan ke:

```text
app/service/student_rules.go
```

---

### Interface Repository

Interface dan implementasi repository ditempatkan pada package yang sama:

```text
app/repository/
```

Dalam implementasi Clean Architecture yang lebih ketat, interface repository biasanya didefinisikan pada layer use case dan implementasinya berada pada adapter.

Pendekatan pada modul ini dipilih untuk menjaga struktur tetap sederhana.

---

### Presenter

Presenter juga disederhanakan menjadi:

```text
helper/response.go
```

Fungsi seperti `Success`, `SuccessList`, `Created`, `Fail`, dan `FailValidation` digunakan untuk menjaga format response tetap konsisten.

---

## Business Rules

### ValidateCreate

Memvalidasi data sebelum student dibuat:

* NIM tidak boleh kosong.
* Nama tidak boleh kosong.
* Grade harus berada pada rentang `0–100`.

### ValidateReplace

Memvalidasi request PUT:

* NIM wajib diisi.
* Nama wajib diisi.
* Grade harus berada pada rentang `0–100`.

### ApplyPatch

Menerapkan perubahan parsial.

Field hanya diubah apabila field tersebut dikirim dalam request PATCH.

### IsEmptyPatch

Memastikan request PATCH memiliki setidaknya satu field yang ingin diubah.

### CountTotalPages

Menghitung jumlah halaman berdasarkan total data dan nilai limit.

---

## Unit Testing

Business rules diuji menggunakan package `testing` bawaan Go.

Test yang tersedia:

```text
TestCountTotalPages
TestApplyPatch
TestValidateCreate
```

Menjalankan test:

```bash
go test ./app/service -v
```

Contoh hasil:

```text
=== RUN   TestCountTotalPages
--- PASS: TestCountTotalPages (0.00s)

=== RUN   TestApplyPatch
--- PASS: TestApplyPatch (0.00s)

=== RUN   TestValidateCreate
--- PASS: TestValidateCreate (0.00s)

PASS
ok      tugas04/app/service       0.002s
```

---

## Structured Logging

Logger menggunakan:

```text
log/slog
```

dengan `lumberjack` untuk rotasi file log.

Output diarahkan ke dua tujuan:

```text
stdout
logs/app.log
```

Konfigurasi rotasi:

| Konfigurasi |   Nilai |
| ----------- | ------: |
| Max Size    |   10 MB |
| Max Backups |  5 file |
| Max Age     | 14 hari |
| Compression |   Aktif |

Log request menggunakan format JSON dan mencatat informasi seperti:

* request ID;
* HTTP method;
* path;
* status code;
* duration;
* IP address.

Contoh:

```json
{
  "time": "2026-09-09T10:15:23.123Z",
  "level": "INFO",
  "msg": "http_request",
  "request_id": "abc123",
  "method": "GET",
  "path": "/api/v1/students",
  "status": 200,
  "duration": 45.678,
  "ip": "::1"
}
```

---

## Middleware

Middleware yang digunakan:

```text
Request ID
Recover
Helmet
CORS
Request Logger
Require JSON
```

### Request ID

Memberikan ID unik pada request sehingga aktivitas request dapat dilacak melalui log.

### Recover

Mencegah panic pada handler menyebabkan server berhenti secara tidak terkontrol.

### Helmet

Menambahkan security-related HTTP headers.

### CORS

Mengatur Cross-Origin Resource Sharing.

### Request Logger

Mencatat informasi setiap HTTP request.

### RequireJSON

Memastikan request yang memiliki body menggunakan:

```text
Content-Type: application/json
```

Jika tidak sesuai, API mengembalikan:

```text
415 Unsupported Media Type
```

---

## API

Base URL:

```text
http://localhost:3000/api/v1
```

### Health Check

```http
GET /health
```

Response ketika server dan database tersedia:

```json
{
  "success": true,
  "message": "server dan database berjalan"
}
```

Jika database tidak dapat dihubungi:

```text
503 Service Unavailable
```

---

### List Students

```http
GET /students
```

Query parameter:

| Parameter   | Default | Keterangan              |
| ----------- | ------: | ----------------------- |
| `page`      |       1 | Nomor halaman           |
| `limit`     |      10 | Jumlah data per halaman |
| `search`    |       - | Pencarian nama          |
| `sort`      |      id | Kolom sorting           |
| `order`     |     asc | `asc` atau `desc`       |
| `is_active` |       - | Filter status aktif     |
| `min_grade` |       - | Grade minimum           |
| `max_grade` |       - | Grade maksimum          |

Kolom sorting yang diperbolehkan:

```text
id
nim
name
grade
created_at
```

Nilai `limit` dibatasi maksimal `100`.

Contoh:

```bash
curl "http://localhost:3000/api/v1/students?page=1&limit=10&sort=name&order=asc"
```

---

### Get Student

```http
GET /students/:id
```

Contoh:

```bash
curl http://localhost:3000/api/v1/students/1
```

---

### Create Student

```http
POST /students
```

Header:

```text
Content-Type: application/json
```

Body:

```json
{
  "nim": "S001",
  "name": "Thaariq",
  "grade": 85.5
}
```

Contoh:

```bash
curl -X POST http://localhost:3000/api/v1/students -H "Content-Type: application/json" -d '{"nim":"S001","name":"Thaariq","grade":85.5}'
```

Response sukses:

```text
201 Created
```

---

### Replace Student

```http
PUT /students/:id
```

PUT mengganti seluruh data yang dapat diubah.

Contoh:

```bash
curl -X PUT http://localhost:3000/api/v1/students/1 -H "Content-Type: application/json" -d '{"nim":"S101","name":"Thoriq","grade":95.0,"is_active":false}'
```

Response:

```text
200 OK
```

---

### Patch Student

```http
PATCH /students/:id
```

PATCH hanya mengubah field yang dikirim.

Contoh:

```bash
curl -X PATCH http://localhost:3000/api/v1/students/1 -H "Content-Type: application/json" -d '{"is_active":true}'
```

---

### Delete Student

```http
DELETE /students/:id
```

Contoh:

```bash
curl -X DELETE http://localhost:3000/api/v1/students/2 -i
```

Response sukses:

```text
204 No Content
```

---

## HTTP Status Code

| Status | Penggunaan                    |
| -----: | ----------------------------- |
|  `200` | Request berhasil              |
|  `201` | Data berhasil dibuat          |
|  `204` | Data berhasil dihapus         |
|  `400` | Request tidak valid           |
|  `404` | Data/endpoint tidak ditemukan |
|  `409` | NIM sudah digunakan           |
|  `415` | Content-Type tidak sesuai     |
|  `422` | Validasi business rules gagal |
|  `503` | Database tidak tersedia       |
|  `500` | Internal server error         |

---

## Error Translation

Repository menggunakan sentinel error:

```go
ErrNotFound
ErrDuplicate
```

Error PostgreSQL kemudian diterjemahkan menjadi error aplikasi dan akhirnya menjadi HTTP response.

Contoh:

```text
pgx.ErrNoRows
      ↓
repository.ErrNotFound
      ↓
HTTP 404
```

Untuk NIM duplikat:

```text
PostgreSQL 23505
      ↓
repository.ErrDuplicate
      ↓
HTTP 409
```

Dengan pendekatan ini, detail database tidak perlu diketahui oleh layer HTTP.

---

## Graceful Shutdown

Aplikasi menangani:

```text
SIGINT
SIGTERM
```

Ketika menerima signal shutdown:

1. Server berhenti menerima request baru.
2. Aplikasi menjalankan proses shutdown.
3. Context shutdown memiliki timeout 10 detik.
4. Server ditutup secara terkontrol.
5. Connection pool database ditutup.

Pendekatan ini mencegah aplikasi berhenti secara tiba-tiba ketika masih terdapat request yang sedang diproses.

---

## Environment Configuration

Buat file:

```text
.env
```

berdasarkan:

```text
.env.example
```

Contoh konfigurasi:

```env
APP_PORT=3000
APP_NAME=API Students - Modul 4 (Clean Architecture)

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=db_students
DB_SSLMODE=disable
DB_MAX_CONNS=10

LOG_LEVEL=info
```

File `.env` tidak boleh di-commit.

---

## `.gitignore`

Konfigurasi utama:

```gitignore
.env
logs/
*.log
```

Dengan demikian credential database dan file log lokal tidak masuk ke repository.

---

## Menjalankan Project

### 1. Masuk ke folder

```bash
cd tugas04
```

### 2. Siapkan environment

```bash
cp .env.example .env
```

Sesuaikan konfigurasi PostgreSQL pada `.env`.

### 3. Install dependency

```bash
go mod tidy
```

### 4. Jalankan unit test

```bash
go test ./app/service/... -v
```

### 5. Jalankan aplikasi

```bash
go run .
```

Server berjalan pada:

```text
http://localhost:3000
```

---

## Build dan Static Check

Menjalankan build:

```bash
go build ./...
```

Menjalankan static analysis:

```bash
go vet ./...
```

Keduanya digunakan untuk memastikan seluruh package dapat dikompilasi dan tidak ditemukan masalah yang terdeteksi oleh `go vet`.

---

## Pemeriksaan Layer Leakage

Beberapa pemeriksaan dapat dilakukan untuk memastikan dependency rule tetap terjaga.

### Repository tidak mengimpor Fiber

```bash
go list -deps ./app/repository | grep fiber
```

Output yang diharapkan:

```text
(kosong)
```

### Service tidak mengandung SQL langsung

```bash
grep -r "SELECT\|INSERT\|UPDATE\|DELETE" app/service/*.go
```

Output yang diharapkan:

```text
(tidak ditemukan)
```

### Route tidak mengandung business rules

```bash
grep -E "if.*==\"\"|Validate|ApplyPatch" route/route.go
```

Output yang diharapkan:

```text
(tidak ditemukan)
```

### Main tidak mengandung handler HTTP

```bash
grep -E "func.*fiber.Ctx|c.JSON" main.go
```

Output yang diharapkan:

```text
(tidak ditemukan)
```

---

## Pengujian Perilaku API

Restrukturisasi Modul 4 ditujukan untuk mengubah struktur internal tanpa mengubah kontrak API.

Data pengujian:

| NIM  | Nama    | Grade |
| ---- | ------- | ----: |
| S001 | Thaariq |  85.5 |
| S002 | Valen   |  92.0 |
| S003 | Rizki   |  78.0 |

Pengujian meliputi:

* POST data student.
* GET seluruh student.
* Pagination.
* Sorting.
* Search.
* GET berdasarkan ID.
* Duplicate NIM.
* PUT.
* PATCH.
* DELETE.
* Invalid Content-Type.
* ID tidak ditemukan.
* Health check.
* Database offline.

Dengan pengujian tersebut, perilaku HTTP dari Modul 3 dapat dibandingkan dengan implementasi Modul 4.

---

## Perbandingan Modul 3 dan Modul 4

| Aspek              | Modul 3           | Modul 4                    |
| ------------------ | ----------------- | -------------------------- |
| Storage            | PostgreSQL        | PostgreSQL                 |
| Repository Pattern | Ya                | Ya                         |
| Fiber              | Ya                | Ya                         |
| Business Rules     | Service/handler   | Dipisahkan                 |
| Architecture       | Layered sederhana | Clean Architecture         |
| Unit Test          | Belum fokus       | Business rules             |
| Logger             | Basic             | Structured logging         |
| Request ID         | Belum             | Ya                         |
| Graceful Shutdown  | Belum fokus       | Ya                         |
| Error Handler      | Service-level     | Global + service           |
| Helper             | Terbatas          | Presenter + request parser |
| Dependency Check   | Belum fokus       | Ya                         |

Perubahan utama Modul 4 berada pada **struktur internal dan pemisahan tanggung jawab**, bukan pada penambahan endpoint baru.

---

## Kelebihan

* Struktur kode lebih terorganisir.
* Business rules dapat diuji secara terisolasi.
* Dependency antar-layer lebih jelas.
* Repository menyembunyikan detail database.
* Logging lebih mudah dianalisis.
* Error handling lebih konsisten.
* Lebih siap dikembangkan untuk proyek yang lebih besar.

---

## Trade-off

Clean Architecture menambah jumlah package dan file.

Untuk proyek kecil, konsekuensinya adalah:

* navigasi kode menjadi lebih panjang;
* satu request dapat melewati beberapa layer;
* boilerplate meningkat;
* struktur terasa lebih kompleks dibanding Modul 3.

Untuk proyek dengan banyak endpoint, developer, atau sumber data, pemisahan tersebut menjadi lebih bermanfaat.

---

## Learning Outcomes

Setelah menyelesaikan Modul 4, konsep yang dipelajari meliputi:

1. Clean Architecture.
2. Dependency Rule.
3. Separation of Concerns.
4. Business Rules.
5. Repository Pattern.
6. Dependency Inversion melalui interface.
7. Unit Testing.
8. Structured Logging.
9. Log Rotation.
10. Middleware composition.
11. Graceful Shutdown.
12. Error Translation.
13. Layer Leakage Analysis.
14. Maintainable backend structure.

---

## Repository

Source code:

```text
https://github.com/vxpal3n/go-workspace/tree/main/tugas04
```

---

## Catatan

Implementasi ini merupakan proyek pembelajaran untuk mata kuliah **Pemrograman Backend Lanjut (SIP375)**.

Struktur Clean Architecture pada modul ini merupakan versi yang disederhanakan dari implementasi kanonik agar sesuai dengan skala proyek dan tujuan pembelajaran.
