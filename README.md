# Go Workspace — Pemrograman Backend Lanjut

Repository ini berisi kumpulan tugas dan latihan praktikum untuk mata kuliah **Pemrograman Backend Lanjut (SIP375)**.

Project dikembangkan secara bertahap mengikuti modul pembelajaran, mulai dari dasar bahasa Go, pengembangan REST API, hingga integrasi database PostgreSQL dan penerapan Repository Pattern.

---

## Daftar Modul

| Modul                  | Topik                         |  Status |
| :--------------------- | :---------------------------- | :-----: |
| [Tugas 01](./tugas01/) | Persiapan & Sintaks Go        | Selesai |
| [Tugas 02](./tugas02/) | REST API & HTTP Deep Dive     | Selesai |
| [Tugas 03](./tugas03/) | Database & Repository Pattern | Selesai |

---

## Progress Pembelajaran

### Modul 1 — Persiapan & Sintaks Go

Modul pertama berfokus pada pengenalan lingkungan pengembangan Go dan dasar-dasar bahasa Go.

Materi yang dipelajari:

* Persiapan environment Go.
* Struktur program Go.
* Variabel dan tipe data.
* Slice.
* Map.
* Pointer.
* Struct.
* Method.

**Dokumentasi:** [Tugas 01](./tugas01/)

---

### Modul 2 — REST API & HTTP Deep Dive

Modul kedua mengembangkan dasar Go menjadi REST API menggunakan Fiber v2.

Implementasi utama:

* REST API.
* HTTP method.
* CRUD Student.
* Request dan response JSON.
* Query parameter.
* Pagination.
* Search.
* Sorting.
* Filtering.
* Validation.
* PUT dan PATCH.
* Middleware.
* Error handling.
* DTO.
* In-memory data storage.

Pada tahap ini data masih disimpan di memory menggunakan `slice`, sehingga data akan hilang ketika aplikasi dihentikan atau dijalankan kembali.

**Dokumentasi:** [Tugas 02](./tugas02/)

---

### Modul 3 — Database & Repository Pattern

Modul ketiga mengembangkan REST API dari Modul 2 dengan mengganti penyimpanan data dari memory menjadi database PostgreSQL.

Implementasi utama:

* PostgreSQL.
* `pgx/v5`.
* `pgxpool` connection pooling.
* Database migration.
* Repository Pattern.
* Separation of concerns.
* Parameterized query.
* Database constraint.
* Database index.
* Filtering dan searching menggunakan SQL.
* Sorting dan pagination pada database.
* Error translation.
* Context timeout.
* Database health check.

Dengan perubahan ini, data mahasiswa menjadi **persistent** dan aplikasi memiliki struktur yang lebih mendekati arsitektur backend production.

**Dokumentasi:** [Tugas 03](./tugas03/)

---

## Struktur Repository

```text id="s8d2jp"
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
└── tugas03/
    ├── README.md
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
    │   └── repository/
    │       └── student_repository.go
    │
    └── migrations/
        └── 001_create_students.sql
```

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
| cURL         | API testing                  |
| Git & GitHub | Version control              |

---

## Perkembangan Project

Repository ini dikembangkan secara bertahap:

```text id="q5p2qz"
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
```

Setiap modul membangun kemampuan dari modul sebelumnya sehingga implementasi dapat berkembang secara bertahap dari konsep dasar menuju backend application yang lebih terstruktur.

---

## Cara Menggunakan Repository

Setiap modul memiliki dokumentasi masing-masing melalui file `README.md`.

Untuk mempelajari modul tertentu:

```bash id="v5y8fh"
cd tugas03
```

Kemudian ikuti instruksi yang terdapat pada README modul tersebut.

Contoh untuk Modul 3:

```bash id="0xy7ja"
cd tugas03
go mod tidy
go run .
```

Konfigurasi dan kebutuhan database untuk Modul 3 dijelaskan secara lengkap pada [README Tugas 03](./tugas03/).

---

## Status

| Modul |  Status | Fokus Utama                     |
| :---: | :-----: | :------------------------------ |
|   01  | Selesai | Fundamental Go                  |
|   02  | Selesai | REST API & HTTP                 |
|   03  | Selesai | PostgreSQL & Repository Pattern |

---

## Catatan

Repository ini merupakan bagian dari proses pembelajaran **Pemrograman Backend Lanjut (SIP375)**.

Setiap modul didokumentasikan secara terpisah agar implementasi, konsep, dan perkembangan project dapat ditelusuri dengan lebih mudah.

Pengembangan kode memanfaatkan bantuan AI untuk debugging dan penyusunan struktur pada beberapa bagian, sementara implementasi dan penyesuaian logika dilakukan sesuai kebutuhan masing-masing modul.

---

## Repository

Source code lengkap tersedia pada repository:

```text id="a6eqk2"
https://github.com/vxpal3n/go-workspace
```
