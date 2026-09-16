# Tugas 03 — Database & Repository Pattern

**Mata Kuliah:** Pemrograman Backend Lanjut (SIP375)  
**Modul:** 3 — Database & Repository Pattern  
**Bahasa:** Go  
**Framework:** Fiber v2  
**Database:** PostgreSQL  
**Driver:** pgx/v5  # Modul 4 — Clean Architecture

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
curl -X PUT http://localhost:3000/api/v1/students/1 -H "Content-Type: application/json" -d '{"nim":"S101","name":"Thaariq Updated","grade":95.0,"is_active":false}'
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

## Status

**Status: Selesai**

Modul 4 berhasil merestrukturisasi API Students menjadi arsitektur yang lebih terorganisir dengan pemisahan business rules, repository, helper, middleware, routing, configuration, database, logging, serta unit testing.

Fokus utama modul ini bukan sekadar membuat API tetap berjalan, tetapi memahami bagaimana struktur internal aplikasi dapat dirancang agar lebih mudah diuji, dipelihara, dan dikembangkan.

---

## Catatan

Implementasi ini merupakan proyek pembelajaran untuk mata kuliah **Pemrograman Backend Lanjut (SIP375)**.

Struktur Clean Architecture pada modul ini merupakan versi yang disederhanakan dari implementasi kanonik agar sesuai dengan skala proyek dan tujuan pembelajaran.

**Konfigurasi:** godotenv  

---

## Deskripsi

Modul 3 merupakan pengembangan dari REST API Students yang telah dibuat pada Modul 2.

Pada Modul 2, data mahasiswa masih disimpan di dalam memory menggunakan `slice`, sehingga seluruh data akan hilang ketika server dihentikan atau dijalankan ulang.

Pada Modul 3, mekanisme penyimpanan tersebut dipindahkan ke **PostgreSQL** sehingga data dapat disimpan secara permanen. Selain itu, implementasi API mulai menggunakan **Repository Pattern** untuk memisahkan logika handler dengan detail akses database.

Perubahan utama pada modul ini meliputi:

* Migrasi penyimpanan data dari memory ke PostgreSQL.
* Implementasi connection pool menggunakan `pgxpool`.
* Penerapan Repository Pattern.
* Pemindahan filter, search, sort, dan pagination ke sisi database.
* Penggunaan parameterized query untuk mencegah SQL Injection.
* Penerapan database constraint untuk menjaga integritas data.
* Penerjemahan error database menjadi HTTP status code yang sesuai.
* Penggunaan context timeout untuk operasi database.
* Penambahan endpoint health check yang memeriksa koneksi database.

---

## Tujuan Pembelajaran

Setelah menyelesaikan modul ini, beberapa konsep yang dipelajari meliputi:

1. Menghubungkan aplikasi Go dengan PostgreSQL.
2. Menggunakan `pgx/v5` dan `pgxpool` untuk koneksi database.
3. Mengelola konfigurasi aplikasi menggunakan environment variable.
4. Membuat dan menjalankan database migration.
5. Memahami dan menerapkan Repository Pattern.
6. Memisahkan handler dari implementasi database.
7. Menggunakan SQL parameter binding.
8. Menerapkan filtering, searching, sorting, dan pagination langsung pada database.
9. Menangani database error secara terstruktur.
10. Memetakan error database ke HTTP status code.
11. Menggunakan database constraint untuk menjaga konsistensi data.
12. Menerapkan context timeout pada operasi database.

---

## Teknologi & Tools

| Teknologi / Tool | Penggunaan                      |
| :--------------- | :------------------------------ |
| Go               | Bahasa pemrograman utama        |
| Fiber v2         | HTTP web framework              |
| PostgreSQL 15+   | Database relasional             |
| pgx/v5           | PostgreSQL driver untuk Go      |
| pgxpool          | Connection pool PostgreSQL      |
| godotenv         | Membaca konfigurasi dari `.env` |
| SQL              | Query dan database migration    |
| cURL             | Pengujian REST API              |
| Git & GitHub     | Version control dan repository  |

---

## Struktur Project

```text
tugas03/
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── main.go
├── handler.go
├── helper.go
│
├── config/
│   └── env.go
│
├── database/
│   └── postgres.go
│
├── app/
│   ├── model/
│   │   └── student.go
│   │
│   └── repository/
│       └── student_repository.go
│
└── migrations/
    └── 001_create_students.sql
```

Struktur tersebut menerapkan prinsip **separation of concerns**, sehingga masing-masing bagian aplikasi memiliki tanggung jawab yang lebih terarah.

---

## Penjelasan Struktur

### `main.go`

Berfungsi sebagai entry point aplikasi.

Tanggung jawab utamanya:

* Memuat konfigurasi environment.
* Membuat koneksi database.
* Membuat repository.
* Membuat handler.
* Menginisialisasi Fiber.
* Memasang middleware.
* Mendaftarkan routing.
* Menjalankan HTTP server.

---

### `handler.go`

Berisi handler/controller untuk endpoint Student.

Handler bertanggung jawab terhadap:

* Menerima HTTP request.
* Parsing request body.
* Validasi input.
* Memanggil repository.
* Menangani hasil operasi.
* Mengembalikan HTTP response.

Handler tidak berisi query SQL secara langsung.

---

### `helper.go`

Berisi fungsi-fungsi pembantu yang digunakan oleh handler, antara lain:

* Response builder.
* Validation response.
* Parsing query parameter.
* Parsing ID dari URL.
* Context timeout untuk operasi database.

Contoh fungsi:

```go
func reqCtx(c *fiber.Ctx) (context.Context, context.CancelFunc) {
    return context.WithTimeout(c.UserContext(), 5*time.Second)
}
```

Dengan demikian, setiap operasi database memiliki batas waktu maksimal 5 detik.

---

### `config/env.go`

Digunakan untuk mengelola environment variable menggunakan `godotenv`.

Beberapa konfigurasi yang digunakan:

```env
APP_PORT=3000

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=db_students
DB_SSLMODE=disable
DB_MAX_CONNS=10
```

File `.env` digunakan untuk konfigurasi lokal dan tidak disimpan ke repository.

---

### `database/postgres.go`

Berfungsi membuat koneksi PostgreSQL menggunakan `pgxpool`.

Connection pool dikonfigurasi dengan:

```text
MaxConns          = 10
MinConns          = 2
MaxConnLifetime   = 1 hour
MaxConnIdleTime   = 30 minutes
```

Aplikasi juga melakukan `Ping()` ketika inisialisasi untuk memastikan database dapat diakses.

---

### `app/model/student.go`

Berisi model dan DTO yang digunakan oleh aplikasi.

Model utama:

```go
type Student struct {
    ID        int       `json:"id"`
    NIM       string    `json:"nim"`
    Name      string    `json:"name"`
    Grade     float64   `json:"grade"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
}
```

Selain model `Student`, terdapat beberapa request DTO:

* `CreateStudentRequest`
* `ReplaceStudentRequest`
* `PatchStudentRequest`
* `ListQuery`

Terdapat juga response wrapper:

* `WebResponse`
* `Meta`

---

### `app/repository/student_repository.go`

Merupakan bagian utama penerapan **Repository Pattern**.

Interface repository mendefinisikan kontrak:

```go
type StudentRepository interface {
    FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error)
    FindByID(ctx context.Context, id int) (model.Student, error)
    Create(ctx context.Context, s model.Student) (model.Student, error)
    Update(ctx context.Context, s model.Student) (model.Student, error)
    Delete(ctx context.Context, id int) error
}
```

Implementasi konkretnya menggunakan PostgreSQL:

```text
StudentRepository
        |
        v
studentPostgresRepository
        |
        v
    PostgreSQL
```

Dengan pendekatan ini, handler tidak perlu mengetahui detail bagaimana data disimpan.

---

### `migrations/001_create_students.sql`

Berisi SQL untuk membuat database schema.

Migration digunakan untuk memastikan struktur database dapat dibuat secara konsisten.

---

## Database Design

### Tabel `students`

Struktur tabel:

```sql
CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    nim VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    grade DECIMAL(5,2) NOT NULL CHECK (grade >= 0 AND grade <= 100),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Kolom

| Kolom        | Tipe           | Keterangan             |
| :----------- | :------------- | :--------------------- |
| `id`         | `SERIAL`       | Primary key            |
| `nim`        | `VARCHAR(20)`  | Nomor induk mahasiswa  |
| `name`       | `VARCHAR(100)` | Nama mahasiswa         |
| `grade`      | `DECIMAL(5,2)` | Nilai mahasiswa, 0–100 |
| `is_active`  | `BOOLEAN`      | Status aktif mahasiswa |
| `created_at` | `TIMESTAMPTZ`  | Waktu data dibuat      |

---

## Database Constraint & Index

### Unique NIM

Database menggunakan unique index:

```sql
CREATE UNIQUE INDEX IF NOT EXISTS students_nim_lower_key
ON students (LOWER(nim));
```

Penggunaan `LOWER(nim)` membuat pengecekan bersifat **case-insensitive**.

Dengan demikian:

```text
S001
s001
```

dianggap sebagai NIM yang sama.

Keunikan dijaga langsung oleh database sehingga lebih aman terhadap kondisi konkurensi dibandingkan hanya melakukan pengecekan dari kode aplikasi.

---

### Index Nama

Untuk mendukung pencarian nama:

```sql
CREATE INDEX IF NOT EXISTS students_name_lower_idx
ON students (LOWER(name));
```

Index ini mendukung pencarian case-insensitive pada kolom `name`.

---

### Check Constraint Grade

Kolom `grade` memiliki constraint:

```sql
CHECK (grade >= 0 AND grade <= 100)
```

Dengan demikian, database tidak akan menerima nilai di luar rentang 0–100.

Constraint ini menjadi lapisan validasi tambahan selain validasi pada aplikasi.

---

## Database Migration

Buat database PostgreSQL:

```bash
psql -U postgres -c "CREATE DATABASE db_students;"
```

Kemudian jalankan migration:

```bash
psql -U postgres -d db_students -f migrations/001_create_students.sql
```

Setelah migration berhasil, tabel `students` beserta constraint dan index akan tersedia.

---

## Repository Pattern

Repository Pattern digunakan untuk memisahkan **logika aplikasi** dari **detail penyimpanan data**.

### Tanpa Repository

Secara sederhana:

```text
Handler
   |
   +---- SQL Query
   |
   +---- PostgreSQL
```

Handler harus mengetahui bagaimana database bekerja.

### Dengan Repository

Pada Modul 3:

```text
HTTP Request
     |
     v
  Handler
     |
     v
Repository Interface
     |
     v
PostgreSQL Repository
     |
     v
PostgreSQL
```

Handler hanya mengetahui kontrak repository.

Contohnya:

```go
student, err := h.repo.FindByID(ctx, id)
```

Handler tidak perlu mengetahui query SQL yang digunakan untuk mengambil data.

### Keuntungan

Repository Pattern memberikan beberapa keuntungan:

* Separation of concerns.
* Handler menjadi lebih sederhana.
* Detail database terisolasi.
* Lebih mudah melakukan perubahan database di masa depan.
* Lebih mudah melakukan testing dengan repository mock/fake.
* Logika penyimpanan tidak bercampur dengan HTTP handling.

---

## Connection Pool

Koneksi database menggunakan `pgxpool`.

Konfigurasi:

```go
cfg.MaxConns = 10
cfg.MinConns = 2
cfg.MaxConnLifetime = time.Hour
cfg.MaxConnIdleTime = 30 * time.Minute
```

Connection pool memungkinkan aplikasi menggunakan kembali koneksi database yang tersedia daripada membuat koneksi baru untuk setiap request.

Selain itu, aplikasi melakukan pengecekan koneksi saat startup:

```go
if err := pool.Ping(pingCtx); err != nil {
    pool.Close()
    return nil, fmt.Errorf("gagal terhubung ke database: %w", err)
}
```

---

## Query Parameterization

Query database menggunakan parameter binding:

```sql
SELECT *
FROM students
WHERE nim = $1
```

Nilai input diberikan sebagai parameter terpisah.

Contoh:

```go
pool.QueryRow(
    ctx,
    "SELECT * FROM students WHERE nim = $1",
    nim,
)
```

Pendekatan ini mencegah input pengguna disisipkan secara langsung ke query SQL dan membantu melindungi aplikasi dari **SQL Injection**.

---

## Filter, Search, Sort, dan Pagination

Pada Modul 2, pengolahan data dilakukan di sisi aplikasi.

Pada Modul 3, proses tersebut dipindahkan ke PostgreSQL.

| Fitur      | Implementasi SQL     |
| :--------- | :------------------- |
| Search     | `ILIKE`              |
| Filter     | `WHERE`              |
| Sort       | `ORDER BY`           |
| Pagination | `LIMIT` dan `OFFSET` |
| Total data | `COUNT(*)`           |

### Search

Contoh:

```text
GET /api/v1/students?search-ri
```

Query akan mencari nama atau NIM yang mengandung teks tersebut secara case-insensitive.

### Filter Grade

```text
GET /api/v1/students?min_grade=80&max_grade=90
```

### Sorting

```text
GET /api/v1/students?sort=name&order=asc
```

Kolom yang diperbolehkan untuk sorting dibatasi menggunakan whitelist:

```text
id
nim
name
grade
created_at
```

Hal ini mencegah input kolom SQL secara sembarangan.

### Pagination

Contoh:

```text
GET /api/v1/students?page=1&limit=2
```

Query menggunakan:

```sql
LIMIT ...
OFFSET ...
```

Sehingga database hanya mengirimkan data yang diperlukan untuk halaman tersebut.

---

# API Documentation

Base URL:

```text
http://localhost:3000/api/v1
```

## Health Check

### `GET /health`

Memeriksa apakah server dan database dapat digunakan.

Response berhasil:

```json
{
  "success": true,
  "message": "server dan database berjalan",
  "data": {
    "timestamp": "..."
  }
}
```

Status:

```text
200 OK
```

Jika database tidak dapat dihubungi:

```text
503 Service Unavailable
```

---

## Get All Students

### `GET /students`

Mengambil daftar mahasiswa.

### Query Parameter

| Parameter   | Keterangan                          |
| :---------- | :---------------------------------- |
| `page`      | Nomor halaman                       |
| `limit`     | Jumlah data per halaman             |
| `search`    | Pencarian berdasarkan nama atau NIM |
| `sort`      | Kolom pengurutan                    |
| `order`     | `asc` atau `desc`                   |
| `is_active` | Filter status aktif                 |
| `min_grade` | Nilai minimum                       |
| `max_grade` | Nilai maksimum                      |

Contoh:

```bash
curl "http://localhost:3000/api/v1/students?page=1&limit=2&sort=name&order=asc"
```

---

## Get Student by ID

### `GET /students/:id`

Mengambil satu mahasiswa berdasarkan ID.

Contoh:

```bash
curl http://localhost:3000/api/v1/students/1
```

Response:

```text
200 OK
```

Jika ID tidak ditemukan:

```text
404 Not Found
```

---

## Create Student

### `POST /students`

Menambahkan mahasiswa baru.

Request:

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

Response:

```text
201 Created
```

Database akan mengisi:

* `id`
* `created_at`

secara otomatis.

---

## Replace Student

### `PUT /students/:id`

PUT digunakan untuk mengganti seluruh data student.

Request:

```json
{
  "nim": "S101",
  "name": "Thoriq",
  "grade": 95.0,
  "is_active": false
}
```

Contoh:

```bash
curl -X PUT http://localhost:3000/api/v1/students/1 -H "Content-Type: application/json" -d '{"nim":"S101","name":"Thoriq","grade":95.0,"is_active":false}'
```

Response:

```text
200 OK
```

Berbeda dengan PATCH, seluruh field yang dibutuhkan harus diberikan.

---

## Partial Update Student

### `PATCH /students/:id`

PATCH digunakan untuk mengubah sebagian data.

Contoh hanya mengubah status:

```json
{
  "is_active": true
}
```

Request:

```bash
curl -X PATCH http://localhost:3000/api/v1/students/1 -H "Content-Type: application/json" -d '{"is_active":true}'
```

Response:

```text
200 OK
```

Field lain tetap mempertahankan nilai sebelumnya.

Jika body kosong:

```json
{}
```

API akan memberikan:

```text
400 Bad Request
```

---

## Delete Student

### `DELETE /students/:id`

Menghapus mahasiswa berdasarkan ID.

Contoh:

```bash
curl -X DELETE http://localhost:3000/api/v1/students/2 -i
```

Response berhasil:

```text
204 No Content
```

Tidak terdapat response body.

---

# HTTP Status Code

API menggunakan status code sesuai kondisi request.

| Status                       | Kondisi                            |
| :--------------------------- | :--------------------------------- |
| `200 OK`                     | Request berhasil                   |
| `201 Created`                | Data berhasil dibuat               |
| `204 No Content`             | Data berhasil dihapus              |
| `400 Bad Request`            | Request tidak valid                |
| `404 Not Found`              | Data atau endpoint tidak ditemukan |
| `409 Conflict`               | NIM sudah digunakan                |
| `415 Unsupported Media Type` | Content-Type bukan JSON            |
| `422 Unprocessable Entity`   | Validasi data gagal                |
| `503 Service Unavailable`    | Database tidak tersedia            |
| `500 Internal Server Error`  | Error server yang tidak terduga    |

---

# Error Translation

Repository menggunakan sentinel error:

```go
var (
    ErrNotFound  = errors.New("data tidak ditemukan")
    ErrDuplicate = errors.New("data sudah ada")
)
```

Kemudian handler menerjemahkannya menjadi HTTP status.

```text
PostgreSQL / pgx
       |
       v
Repository Error
       |
       v
translateError()
       |
       +---- ErrNotFound  -> 404
       |
       +---- ErrDuplicate -> 409
       |
       +---- Error lain   -> 500
```

Contoh:

```text
pgx.ErrNoRows
      ↓
repository.ErrNotFound
      ↓
404 Not Found
```

Untuk duplicate NIM:

```text
PostgreSQL Code 23505
      ↓
repository.ErrDuplicate
      ↓
409 Conflict
```

Error mentah dari PostgreSQL tidak dikirim langsung kepada client sehingga detail internal database tidak terekspos.

---

# Middleware

Aplikasi menggunakan beberapa middleware Fiber.

### Request ID

```go
requestid.New()
```

Digunakan untuk memberikan identifier pada setiap request.

### Logger

Logger mencatat informasi seperti:

```text
time
request ID
HTTP method
path
status
latency
```

### CORS

```go
cors.New()
```

Digunakan untuk mengaktifkan dukungan Cross-Origin Resource Sharing.

### Content-Type Validation

Request dengan body (`POST`, `PUT`, `PATCH`) harus menggunakan:

```text
Content-Type: application/json
```

Jika tidak:

```text
415 Unsupported Media Type
```

---

# Context Timeout

Operasi database menggunakan timeout:

```go
context.WithTimeout(c.UserContext(), 5*time.Second)
```

Tujuannya adalah mencegah query yang bermasalah atau terlalu lama menahan resource database tanpa batas waktu.

Health check menggunakan timeout yang lebih singkat:

```go
context.WithTimeout(c.UserContext(), 2*time.Second)
```

---

# Cara Menjalankan Project

## Prasyarat

Pastikan telah tersedia:

* Go 1.27.0
* PostgreSQL 15+
* `psql`
* Git

---

## 1. Clone Repository

```bash
git clone https://github.com/vxpal3n/go-workspace.git
cd go-workspace/tugas03
```

## 2. Konfigurasi Environment

Salin `.env.example` menjadi `.env`.

```env
APP_PORT=3000

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=db_students
DB_SSLMODE=disable
DB_MAX_CONNS=10
```

Sesuaikan `DB_USER` dan `DB_PASSWORD` dengan konfigurasi PostgreSQL lokal.

> File `.env` tidak boleh di-commit ke repository karena dapat berisi kredensial database.

## 3. Buat Database

```bash
psql -U postgres -c "CREATE DATABASE db_students;"
```

## 4. Jalankan Migration

```bash
psql -U postgres -d db_students -f migrations/001_create_students.sql
```

## 5. Install Dependency

```bash
go mod tidy
```

## 6. Jalankan Server

```bash
go run .
```

Server akan berjalan pada:

```text
http://localhost:3000
```

---

# Testing

Berikut urutan pengujian yang digunakan untuk memastikan seluruh fitur API berjalan.

## 1. Reset Database

Untuk memulai pengujian dari kondisi bersih:

```bash
psql -U postgres -d db_students -c "TRUNCATE TABLE students RESTART IDENTITY;"
```

Perintah tersebut menghapus seluruh data dan mengembalikan sequence ID ke awal.

---

## 2. Insert Data Awal

Tambahkan tiga mahasiswa:

```bash
curl -X POST http://localhost:3000/api/v1/students -H "Content-Type: application/json" -d '{"nim":"S001","name":"Thaariq","grade":85.5}'
```

```bash
curl -X POST http://localhost:3000/api/v1/students -H "Content-Type: application/json" -d '{"nim":"S002","name":"Valen","grade":92.0}'
```

```bash
curl -X POST http://localhost:3000/api/v1/students -H "Content-Type: application/json" -d '{"nim":"S003","name":"Rizki","grade":78.0}'
```

Ketiganya diharapkan menghasilkan:

```text
201 Created
```

---

## 3. Pagination & Sorting

```bash
curl "http://localhost:3000/api/v1/students?page=1&limit=2&sort=name&order=asc"
```

Hasil yang diharapkan:

```text
200 OK
```

Dengan:

```text
total = 3
total_pages = 2
```

---

## 4. Search

```bash
curl "http://localhost:3000/api/v1/students?search-ri"
```

Pencarian dilakukan terhadap nama atau NIM secara case-insensitive.

---

## 5. Filter Grade

```bash
curl "http://localhost:3000/api/v1/students?min_grade=80&max_grade=90"
```

Data yang memenuhi kondisi:

```text
Thaariq — 85.5
```

---

## 6. Get by ID

Request berhasil:

```bash
curl http://localhost:3000/api/v1/students/1
```

Expected:

```text
200 OK
```

Request dengan ID yang tidak tersedia:

```bash
curl http://localhost:3000/api/v1/students/999
```

Expected:

```text
404 Not Found
```

---

## 7. Duplicate NIM

Pada tahap ini, **S001 harus masih dimiliki oleh Thaariq**.

Jangan melakukan PUT terlebih dahulu karena PUT pada tahap berikutnya akan mengubah NIM tersebut.

Test:

```bash
curl -X POST http://localhost:3000/api/v1/students -H "Content-Type: application/json" -d '{"nim":"S001","name":"Dup","grade":80}'
```

Expected:

```text
409 Conflict
```

Response message:

```text
NIM sudah digunakan
```

Alurnya:

```text
NIM S001 sudah ada
       ↓
INSERT baru dengan S001
       ↓
UNIQUE INDEX violation
       ↓
PostgreSQL error 23505
       ↓
ErrDuplicate
       ↓
409 Conflict
```

---

## 8. PUT

Setelah pengujian 409 selesai, lakukan PUT:

```bash
curl -X PUT http://localhost:3000/api/v1/students/1 -H "Content-Type: application/json" -d '{"nim":"S101","name":"Thoriq","grade":95.0,"is_active":false}'
```

Expected:

```text
200 OK
```

Data ID 1 berubah menjadi:

```text
NIM       : S101
Name      : Thoriq
Grade     : 95
Is Active : false
```

---

## 9. PUT Validation

Test PUT tanpa field `name`:

```bash
curl -X PUT http://localhost:3000/api/v1/students/1 -H "Content-Type: application/json" -d '{"nim":"S001","grade":90}'
```

Expected:

```text
422 Unprocessable Entity
```

Karena `name` wajib diberikan pada PUT.

---

## 10. PATCH

Ubah hanya `is_active`:

```bash
curl -X PATCH http://localhost:3000/api/v1/students/1 -H "Content-Type: application/json" -d '{"is_active":true}'
```

Expected:

```text
200 OK
```

Hasil:

```text
NIM       : S101
Name      : Thoriq
Grade     : 95
Is Active : true
```

Hanya `is_active` yang berubah.

---

## 11. PATCH Empty Body

```bash
curl -X PATCH http://localhost:3000/api/v1/students/1 -H "Content-Type: application/json" -d '{}'
```

Expected:

```text
400 Bad Request
```

Message:

```text
tidak ada field yang diubah
```

---

## 12. DELETE

Hapus student dengan ID 2:

```bash
curl -X DELETE http://localhost:3000/api/v1/students/2 -i
```

Expected:

```text
204 No Content
```

Kemudian coba ambil kembali:

```bash
curl http://localhost:3000/api/v1/students/2
```

Expected:

```text
404 Not Found
```

---

## 13. Content-Type Validation

POST tanpa `Content-Type`:

```bash
curl -X POST http://localhost:3000/api/v1/students -d '{"nim":"x"}'
```

Expected:

```text
415 Unsupported Media Type
```

---

## 14. Health Check

Dengan database aktif:

```bash
curl http://localhost:3000/api/v1/health
```

Expected:

```text
200 OK
```

Jika PostgreSQL dihentikan dan endpoint health dipanggil:

```bash
curl http://localhost:3000/api/v1/health
```

Expected:

```text
503 Service Unavailable
```

Message:

```text
database tidak dapat dihubungi
```

---

# Ringkasan Skenario Pengujian

| Skenario                | Method | Expected Status |
| :---------------------- | :----: | :-------------: |
| Membuat student         |  POST  |      `201`      |
| List student            |   GET  |      `200`      |
| Search                  |   GET  |      `200`      |
| Filter grade            |   GET  |      `200`      |
| Get by ID               |   GET  |      `200`      |
| ID tidak ditemukan      |   GET  |      `404`      |
| Duplicate NIM           |  POST  |      `409`      |
| PUT berhasil            |   PUT  |      `200`      |
| PUT invalid             |   PUT  |      `422`      |
| PATCH berhasil          |  PATCH |      `200`      |
| PATCH tanpa field       |  PATCH |      `400`      |
| DELETE berhasil         | DELETE |      `204`      |
| DELETE data tidak ada   | DELETE |      `404`      |
| Content-Type salah      |  POST  |      `415`      |
| Health check            |   GET  |      `200`      |
| Database tidak tersedia |   GET  |      `503`      |

---

# Perbandingan Modul 2 dan Modul 3

| Aspek               | Modul 2        | Modul 3        |
| :------------------ | :------------- | :------------- |
| Storage             | Memory / Slice | PostgreSQL     |
| Persistence         | Tidak permanen | Permanen       |
| Filter              | Go             | SQL            |
| Search              | Go             | PostgreSQL     |
| Sorting             | Go             | SQL            |
| Pagination          | Go             | SQL            |
| Repository Pattern  | Belum          | Implementasi   |
| Connection Pool     | Belum          | `pgxpool`      |
| Database Constraint | Belum          | Implementasi   |
| Parameterized Query | Belum          | Implementasi   |
| Error Translation   | Dasar          | Database-aware |
| Health Check DB     | Belum          | Implementasi   |

Perubahan tersebut membuat aplikasi lebih dekat dengan pola backend yang digunakan pada aplikasi nyata.

---

# Security & Best Practices

Beberapa praktik yang diterapkan pada Modul 3:

### 1. Parameterized Query

Input pengguna tidak digabung langsung ke query SQL.

### 2. Sorting Whitelist

Kolom untuk `ORDER BY` dibatasi:

```text
id
nim
name
grade
created_at
```

### 3. Database Constraint

Keunikan NIM dijaga langsung oleh database.

### 4. Environment Variable

Credential database disimpan melalui `.env`, bukan hard-code di source code.

### 5. `.gitignore`

File `.env`, log, dan temporary file tidak dimasukkan ke repository.

### 6. Context Timeout

Operasi database dibatasi waktu untuk mencegah resource tertahan terlalu lama.

### 7. Error Abstraction

Error internal PostgreSQL tidak dikirim langsung kepada client.

### 8. Pagination Limit

Nilai `limit` dibatasi maksimal 100 untuk mencegah request mengambil data dalam jumlah berlebihan.

---

# Keterbatasan

Implementasi Modul 3 masih memiliki beberapa batasan:

* Database masih menggunakan satu tabel utama `students`.
* Belum menggunakan migration tool khusus.
* Belum menggunakan authentication dan authorization.
* Belum menggunakan automated integration test.
* Belum menggunakan ORM.
* Repository masih memiliki implementasi khusus PostgreSQL.
* Konfigurasi masih ditujukan untuk lingkungan development.

Batasan tersebut dapat dikembangkan pada modul atau project berikutnya.

---

# Learning Outcomes

Melalui Modul 3, REST API Students mengalami peningkatan dari aplikasi yang hanya menyimpan data di memory menjadi aplikasi yang menggunakan database relasional secara permanen.

Konsep penting yang berhasil diterapkan meliputi:

1. Integrasi Go dengan PostgreSQL.
2. Connection pooling menggunakan `pgxpool`.
3. Database migration.
4. Repository Pattern.
5. Separation of concerns.
6. Parameterized SQL query.
7. Database constraint dan index.
8. Filtering dan pagination pada database.
9. Error translation.
10. Context timeout.
11. HTTP status code yang sesuai dengan kondisi aplikasi.

Dengan implementasi tersebut, struktur aplikasi menjadi lebih terorganisir dan lebih siap dikembangkan dibandingkan implementasi Modul 2.

---

# Referensi

* Dokumentasi Fiber v2
* Dokumentasi `pgx/v5`
* Dokumentasi PostgreSQL
* Dokumentasi Go
* Materi Praktikum Pemrograman Backend Lanjut — Modul 3

---

## Repository

[GitHub — go-workspace](https://github.com/vxpal3n/go-workspace/tree/main/tugas03)

## Catatan

Project ini merupakan implementasi pembelajaran untuk mata kuliah **Pemrograman Backend Lanjut (SIP375)**.

Pengembangan kode memanfaatkan bantuan AI untuk debugging dan penyusunan struktur pada beberapa bagian, sementara implementasi dan penyesuaian logika dilakukan sesuai kebutuhan modul.
