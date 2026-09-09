# Tugas 3 — Database & Repository Pattern

**Mata Kuliah:** Pemrograman Backend Lanjut (SIP375)
**Modul:** 3 — Database & Repository Pattern
**Bahasa:** Go
**Framework:** Fiber v2
**Database:** PostgreSQL
**Driver:** pgx/v5
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
