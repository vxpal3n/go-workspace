# Go Workspace — Pemrograman Backend Lanjut

Repository ini berisi kumpulan tugas dan latihan praktikum untuk mata kuliah **Pemrograman Backend Lanjut (SIP375)**.

Project dikembangkan secara bertahap mengikuti modul pembelajaran, mulai dari dasar bahasa Go, pengembangan REST API, integrasi database PostgreSQL, penerapan Repository Pattern, hingga penerapan Clean Architecture dan pengujian unit.

---

## Daftar Modul

|          Modul         | Topik                         |  Status |
| :--------------------: | :---------------------------- | :-----: |
| [Tugas 01](./tugas01/) | Persiapan & Sintaks Go        | Selesai |
| [Tugas 02](./tugas02/) | REST API & HTTP Deep Dive     | Selesai |
| [Tugas 03](./tugas03/) | Database & Repository Pattern | Selesai |
| [Tugas 04](./tugas04/) | Clean Architecture            | Selesai |

---

## Progress Pembelajaran

### Modul 1 — Persiapan & Sintaks Go

Modul pertama berfokus pada pengenalan lingkungan pengembangan Go dan dasar-dasar bahasa Go.

Materi yang dipelajari:

* Persiapan environment Go
* Struktur program Go
* Variabel dan tipe data
* Slice
* Map
* Pointer
* Struct
* Method

**Dokumentasi:** [Tugas 01](./tugas01/)

---

### Modul 2 — REST API & HTTP Deep Dive

Modul kedua mengembangkan dasar Go menjadi REST API menggunakan Fiber v2.

Implementasi utama:

* REST API
* HTTP method
* CRUD Student
* Request dan response JSON
* Query parameter
* Pagination
* Search
* Sorting
* Filtering
* Validation
* PUT dan PATCH
* Middleware
* Error handling
* DTO
* In-memory data storage

Pada tahap ini data masih disimpan di memory menggunakan `slice`, sehingga data akan hilang ketika aplikasi dihentikan atau dijalankan kembali.

**Dokumentasi:** [Tugas 02](./tugas02/)

---

### Modul 3 — Database & Repository Pattern

Modul ketiga mengembangkan REST API dari Modul 2 dengan mengganti penyimpanan data dari memory menjadi database PostgreSQL.

Implementasi utama:

* PostgreSQL
* `pgx/v5`
* `pgxpool` connection pooling
* Database migration
* Repository Pattern
* Separation of concerns
* Parameterized query
* Database constraint
* Database index
* Filtering dan searching menggunakan SQL
* Sorting dan pagination pada database
* Error translation
* Context timeout
* Database health check

Dengan perubahan ini, data mahasiswa menjadi **persistent** dan aplikasi memiliki struktur yang lebih mendekati arsitektur backend production.

**Dokumentasi:** [Tugas 03](./tugas03/)

---

### Modul 4 — Clean Architecture

Modul keempat melakukan restrukturisasi terhadap API Students dari Modul 3 menggunakan pendekatan **Clean Architecture**.

Fokus utama modul:

* Clean Architecture
* Dependency Rule
* Separation of business rules
* Repository layer
* Service layer
* Helper dan presenter
* Middleware
* Structured logging
* Unit testing
* Global error handling
* Graceful shutdown
* Layer leakage analysis
* Mempertahankan behavior API dari modul sebelumnya

Struktur aplikasi dipisahkan berdasarkan tanggung jawab sehingga business rules tidak bergantung langsung pada framework maupun database.

Business rules utama ditempatkan pada service layer, meliputi:

* `ValidateCreate`
* `ValidateReplace`
* `ApplyPatch`
* `IsEmptyPatch`
* `CountTotalPages`

Modul ini juga menambahkan unit test untuk business rules menggunakan package `testing`, serta structured logging menggunakan `log/slog` dan log rotation menggunakan `lumberjack`.

**Dokumentasi:** [Tugas 04](./tugas04/)

---

## Struktur Repository

```text
go-workspace/

├── README.md
│
├── tugas01/
│   ├── README.md
│   ├── fiber/
│   │   └── main.go
│   └── syntax/
│       ├── variabel/
│       ├── pointer/
│       └── struct/
│
├── tugas02/
│   ├── README.md
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── model.go
│   ├── helper.go
│   └── handler.go
│
├── tugas03/
│   ├── README.md
│   ├── .env.example
│   ├── .gitignore
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── handler.go
│   ├── helper.go
│   │
│   ├── config/
│   │   └── env.go
│   │
│   ├── database/
│   │   └── postgres.go
│   │
│   ├── app/
│   │   ├── model/
│   │   │   └── student.go
│   │   └── repository/
│   │       └── student_repository.go
│   │
│   └── migrations/
│       └── 001_create_students.sql
│
└── tugas04/
    ├── README.md
    ├── .env.example
    ├── .gitignore
    ├── go.mod
    ├── go.sum
    ├── main.go
    │
    ├── config/
    │   └── app.go
    │
    ├── database/
    │   └── postgres.go
    │
    ├── app/
    │   ├── model/
    │   │   └── student.go
    │   ├── repository/
    │   │   └── student_repository.go
    │   └── service/
    │       ├── student_service.go
    │       └── student_rules_test.go
    │
    ├── helper/
    │   └── ...
    │
    ├── middleware/
    │   └── ...
    │
    ├── route/
    │   └── route.go
    │
    ├── logs/
    │   └── app.log
    │
    └── migrations/
        └── ...
```

> Struktur file dapat berkembang mengikuti kebutuhan masing-masing modul.

---

## Teknologi

Teknologi yang digunakan berkembang seiring dengan bertambahnya materi pada setiap modul.

| Teknologi    | Penggunaan                   |
| :----------- | :--------------------------- |
| Go           | Bahasa pemrograman utama     |
| Fiber v2     | Web framework                |
| PostgreSQL   | Database relasional          |
| pgx/v5       | PostgreSQL driver            |
| pgxpool      | Database connection pooling  |
| godotenv     | Environment configuration    |
| SQL          | Query dan database migration |
| `log/slog`   | Structured logging           |
| `lumberjack` | Log rotation                 |
| `testing`    | Unit testing                 |
| cURL         | API testing                  |
| Git & GitHub | Version control              |

---

## Arsitektur Project

Perkembangan arsitektur backend pada repository ini dilakukan secara bertahap.

```text
Go Fundamentals
       |
       v
REST API
       |
       v
CRUD & HTTP
       |
       v
PostgreSQL
       |
       v
Repository Pattern
       |
       v
Database-backed REST API
       |
       v
Clean Architecture
       |
       v
Business Rules & Unit Testing
       |
       v
Structured Logging
```

Setiap modul membangun kemampuan dari modul sebelumnya. Dengan pendekatan ini, project berkembang dari implementasi sederhana menuju aplikasi backend yang memiliki pemisahan tanggung jawab dan struktur yang lebih terorganisir.

---

## API Students

Mulai Modul 2 hingga Modul 4, project menggunakan API Students sebagai project utama yang dikembangkan secara bertahap.

Base URL:

```text
http://localhost:3000/api/v1
```

Endpoint utama:

| Method | Endpoint        | Keterangan                         |
| :----: | :-------------- | :--------------------------------- |
|   GET  | `/health`       | Health check                       |
|   GET  | `/students`     | Mendapatkan daftar student         |
|   GET  | `/students/:id` | Mendapatkan student berdasarkan ID |
|  POST  | `/students`     | Menambahkan student                |
|   PUT  | `/students/:id` | Mengganti data student             |
|  PATCH | `/students/:id` | Memperbarui sebagian data student  |
| DELETE | `/students/:id` | Menghapus student                  |

Behavior API dipertahankan ketika project berpindah dari Modul 3 ke Modul 4. Perubahan utama pada Modul 4 berada pada struktur internal dan arsitektur aplikasi, bukan pada kontrak HTTP API.

---

## Cara Menggunakan Repository

Setiap modul memiliki dokumentasi masing-masing melalui file `README.md`.

Untuk mempelajari modul tertentu:

```bash
cd tugas04
```

Kemudian ikuti instruksi yang terdapat pada README modul tersebut.

Contoh menjalankan Modul 4:

```bash
cd tugas04

go mod tidy

go run .
```

Untuk menjalankan unit test:

```bash
go test ./app/service/... -v
```

Untuk melakukan build dan pemeriksaan kode:

```bash
go build ./...
go vet ./...
```

Konfigurasi dan kebutuhan masing-masing modul dijelaskan secara lengkap pada README modul terkait.

---

## Status

| Modul |  Status | Fokus Utama                           |
| :---: | :-----: | :------------------------------------ |
|   01  | Selesai | Fundamental Go                        |
|   02  | Selesai | REST API & HTTP                       |
|   03  | Selesai | PostgreSQL & Repository Pattern       |
|   04  | Selesai | Clean Architecture, Testing & Logging |

---

## Catatan

Repository ini merupakan bagian dari proses pembelajaran **Pemrograman Backend Lanjut (SIP375)**.

Setiap modul didokumentasikan secara terpisah agar implementasi, konsep, dan perkembangan project dapat ditelusuri dengan lebih mudah.

Pengembangan kode memanfaatkan bantuan AI untuk debugging, eksplorasi konsep, dan penyusunan struktur pada beberapa bagian. Implementasi dan penyesuaian logika dilakukan sesuai kebutuhan masing-masing modul.

---

## Repository

Source code lengkap tersedia pada repository:

**GitHub:**
https://github.com/vxpal3n/go-workspace
