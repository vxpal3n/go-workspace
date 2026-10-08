# UTS – SIAKAD Mini RESTful API Backend

Implementasi **Ujian Tengah Semester (UTS)** untuk mata kuliah **Pemrograman Backend Lanjut (SIP375)**.

SIAKAD Mini merupakan RESTful API backend untuk layanan akademik sederhana yang mengelola tiga domain utama:

* Data mahasiswa.
* Mata kuliah.
* Kartu Rencana Studi (KRS) / enrollment.

Project ini menggabungkan konsep-konsep yang telah dibangun pada modul sebelumnya ke dalam satu aplikasi backend yang menerapkan **Clean Architecture, authentication, authorization, business rules, database transaction, concurrency control, soft delete, rate limiting, dan automated testing**.

---

## Tujuan

UTS ini bertujuan membangun backend SIAKAD sederhana dengan kontrak API yang jelas serta business rules yang dapat diuji secara terisolasi.

Implementasi mencakup:

* RESTful API menggunakan Go dan Fiber.
* Authentication menggunakan JWT.
* Authorization berdasarkan role.
* Pengelolaan data mahasiswa.
* Pengelolaan mata kuliah.
* Enrollment KRS.
* Batas SKS berdasarkan IPK.
* Pencegahan duplicate enrollment.
* Pencegahan enrollment ketika kuota penuh.
* Ownership pada data mahasiswa.
* Database transaction.
* Row-level locking.
* Soft delete.
* Login guard untuk student yang telah dihapus.
* Login rate limiting.
* Unit testing.
* End-to-end testing.

---

# Arsitektur

Project menggunakan **Clean Architecture** yang dikembangkan sejak Modul 4.

Pembagian layer:

| Layer                | Folder                                                     | Tanggung Jawab                                  |
| :------------------- | :--------------------------------------------------------- | :---------------------------------------------- |
| Entities             | `app/model/`                                               | Struct, DTO, dan validation tag                 |
| Use Cases            | `app/service/*_rules.go`                                   | Business rules murni                            |
| Interface Adapters   | `app/service/*_service.go`, `app/repository/`, `helper/`   | Orchestration, repository, dan helper           |
| Frameworks & Drivers | `config/`, `database/`, `middleware/`, `route/`, `main.go` | Fiber, PostgreSQL, logger, dependency injection |

Business rules dipisahkan dari HTTP dan database sehingga dapat diuji tanpa menjalankan server.

Arsitektur ini merupakan kelanjutan langsung dari struktur yang dibangun pada Modul 4, 5, 6, dan 7.

---

# Domain

SIAKAD Mini memiliki tiga domain utama.

## Students

Data mahasiswa dikelola oleh admin.

Fitur:

* Membuat mahasiswa.
* Melihat mahasiswa.
* Mengubah mahasiswa.
* Menghapus mahasiswa menggunakan soft delete.
* Ownership untuk akses mahasiswa terhadap datanya sendiri.

---

## Courses

Mata kuliah dikelola oleh admin dan dapat dibaca oleh seluruh role yang telah terautentikasi.

Data mata kuliah digunakan sebagai sumber enrollment.

---

## Enrollments / KRS

Mahasiswa dapat:

* Mengambil mata kuliah.
* Membatalkan enrollment miliknya sendiri.

Enrollment memiliki aturan:

* Batas total SKS.
* Tidak boleh mengambil mata kuliah yang sama pada tahun akademik yang sama.
* Tidak dapat mengambil mata kuliah yang kuotanya penuh.
* Hanya mahasiswa yang dapat melakukan enrollment.

---

# Role

SIAKAD Mini menggunakan dua role:

| Role        | Tanggung Jawab                                            |
| :---------- | :-------------------------------------------------------- |
| `admin`     | Mengelola data mahasiswa dan mata kuliah                  |
| `mahasiswa` | Melihat data yang diperbolehkan dan mengelola KRS sendiri |

Access control diterapkan pada endpoint berdasarkan role dan ownership.

---

# Business Rules

Empat business rule utama diterapkan dan diuji.

## 1. Batas SKS Berdasarkan IPK

Batas pengambilan SKS ditentukan berdasarkan IPK terakhir:

| IPK           | Batas SKS |
| :------------ | :-------- |
| `>= 3.00`     | 24 SKS    |
| `2.50 – 2.99` | 21 SKS    |
| `< 2.50`      | 18 SKS    |

Contoh data seeder:

| Mahasiswa    |  IPK | Batas |
| :----------- | ---: | ----: |
| Rina Putri   | 3.55 |    24 |
| Citra Dewi   | 2.87 |    21 |
| Doni Pratama | 2.34 |    18 |

Business rule dipisahkan menjadi pure function sehingga dapat diuji tanpa database.

---

## 2. Duplicate Enrollment

Mahasiswa tidak dapat mengambil mata kuliah yang sama pada tahun akademik yang sama.

Database memiliki constraint:

```sql
UNIQUE (student_id, course_id, tahun_akademik)
```

Aturan dijaga pada dua lapis:

```text
Service check
     +
Database constraint
```

Jika request kedua dilakukan, API mengembalikan:

```text
409 Conflict
```

Pendekatan dua lapis digunakan sebagai defense in depth.

---

## 3. Kuota Mata Kuliah

Mahasiswa tidak dapat mengambil mata kuliah ketika kuotanya telah penuh.

Response:

```text
422 Unprocessable Entity
```

Pemeriksaan kuota dilakukan di dalam transaction menggunakan row lock:

```sql
SELECT ...
FOR UPDATE
```

Dengan demikian dua request concurrent tidak dapat secara bersamaan menganggap slot terakhir masih tersedia.

---

## 4. Ownership

Mahasiswa hanya dapat mengakses atau memodifikasi data yang menjadi miliknya.

Contoh:

```text
Mahasiswa A
    |
    +---- GET student A
    |        -> 200 OK
    |
    +---- GET student B
             -> 403 Forbidden
```

Ownership check ditempatkan pada service karena keputusan membutuhkan data target.

---

# Database Design

Database menggunakan PostgreSQL.

Beberapa keputusan desain penting:

### Student–User

`students.user_id` memiliki constraint `UNIQUE` untuk merepresentasikan hubungan satu-ke-satu dengan user.

### Enrollment

```sql
UNIQUE (student_id, course_id, tahun_akademik)
```

digunakan untuk mencegah duplicate enrollment pada level database.

### Course Deletion

`enrollments.course_id` menggunakan:

```text
ON DELETE RESTRICT
```

sehingga course yang sudah digunakan oleh enrollment tidak dapat dihapus secara sembarangan.

### Soft Delete

Student menggunakan:

```text
deleted_at
```

sebagai penanda penghapusan.

Query repository selalu memfilter:

```sql
WHERE deleted_at IS NULL
```

Index parsial digunakan untuk mempercepat akses terhadap student yang masih aktif.

### GPA

IPK disimpan menggunakan:

```sql
NUMERIC(3,2)
```

agar representasi nilai desimal tetap tepat.

---

# Seeder

Seeder menyediakan data awal untuk kebutuhan pengujian.

## Admin

```text
Email    : admin@siakad.test
Password : admin12345
Role     : admin
```

## Mahasiswa

Terdapat 20 mahasiswa.

Format email:

```text
{NIM}@student.siakad.test
```

Password menggunakan NIM masing-masing.

## Mata Kuliah

Seeder menyediakan 10 mata kuliah yang tersebar pada semester 1 sampai 6.

Password di-hash menggunakan PostgreSQL `pgcrypto` dengan:

```sql
crypt(..., gen_salt('bf', 12))
```

yang kompatibel dengan verifikasi bcrypt pada aplikasi Go.

---

# Transaction pada `POST /students`

Pembuatan student melibatkan lebih dari satu tabel:

```text
users
  |
  v
students
```

Karena itu proses dilakukan menggunakan database transaction.

Alur:

```text
BEGIN
  |
  +---- INSERT users
  |
  +---- INSERT students
  |
  +---- COMMIT
```

Jika salah satu operasi gagal:

```text
ROLLBACK
```

Dengan demikian tidak terjadi data parsial seperti:

```text
users berhasil dibuat
students gagal dibuat
```

Transaction memastikan kedua operasi berhasil bersama-sama atau seluruh perubahan dibatalkan.

---

# Row Locking pada `POST /enrollments`

Enrollment memiliki risiko race condition pada kuota.

Contoh:

```text
Quota = 1

Request A ----+
              |
Request B ----+---- melihat slot terakhir
```

Tanpa locking, keduanya dapat menganggap slot tersedia.

Implementasi menggunakan:

```sql
SELECT ...
FOR UPDATE
```

sehingga row course dikunci selama transaction.

Alurnya:

```text
BEGIN
  |
  v
Lock course row
  |
  v
Check quota
  |
  v
Insert enrollment
  |
  v
COMMIT
```

Dengan demikian pemeriksaan quota dan perubahan data terjadi secara aman.

---

# Soft Delete

Student tidak langsung dihapus dari database.

Sebagai gantinya:

```text
deleted_at = current timestamp
```

Student yang telah di-soft delete tidak muncul pada query aktif.

Keuntungan:

* Riwayat data tetap tersedia.
* Referential integrity lebih terjaga.
* Penghapusan dapat dilacak.
* Data masih dapat dipulihkan jika dibutuhkan.

Konsekuensinya, setiap query student harus disiplin menggunakan:

```sql
deleted_at IS NULL
```

---

# Login Guard

Soft delete juga harus diperhitungkan pada authentication.

Sebelumnya terdapat kemungkinan:

```text
Student di-soft delete
        |
        v
Password masih valid
        |
        v
Login berhasil
```

Hal tersebut tidak sesuai dengan status bisnis student.

Implementasi kemudian memastikan `FindByUserID` hanya mengambil student aktif.

Hasil akhirnya:

```text
Student aktif
    -> login berhasil

Student soft-deleted
    -> 401 Unauthorized
```

Bug ini ditemukan selama pengembangan dan diperbaiki sebelum pengujian akhir.

---

# UPDATE RETURNING Bug

Salah satu bug yang ditemukan terjadi pada PostgreSQL ketika query `UPDATE ... RETURNING` menggunakan alias yang tidak sesuai.

Query update membutuhkan column list tanpa alias tertentu untuk bagian `RETURNING`.

Perbaikannya adalah memisahkan:

```text
studentColumns
```

untuk query yang membutuhkan alias dengan:

```text
studentColumnsBare
```

untuk query `UPDATE ... RETURNING`.

Setelah perbaikan:

```text
PUT /students/1
    |
    v
200 OK
```

dan response student berhasil dikembalikan.

---

# Authentication

Authentication menggunakan JWT.

Login menghasilkan access token dengan masa berlaku:

```text
900 detik
```

atau:

```text
15 menit
```

Request terproteksi menggunakan:

```text
Authorization: Bearer <token>
```

Endpoint login merupakan endpoint publik.

Endpoint lainnya memerlukan authentication sesuai access matrix.

---

# Rate Limiting

Endpoint login dilindungi rate limiter:

```text
5 attempts / minute / IP
```

Pengujian:

```text
Attempt 1 → 401
Attempt 2 → 401
Attempt 3 → 401
Attempt 4 → 401
Attempt 5 → 401
Attempt 6 → 429
```

Response `429 Too Many Requests` dilengkapi:

```text
Retry-After: 60
```

Rate limiting digunakan untuk mengurangi risiko brute-force login.

---

# API Contract

Base URL:

```text
http://localhost:3000/api/v1
```

SIAKAD Mini memiliki 10 endpoint utama.

| No. | Method   | Endpoint            | Fungsi                           |
| :-: | :------- | :------------------ | :------------------------------- |
|  1  | `POST`   | `/auth/login`       | Login                            |
|  2  | `GET`    | `/auth/me`          | Informasi user yang sedang login |
|  3  | `GET`    | `/students`         | Daftar mahasiswa                 |
|  4  | `POST`   | `/students`         | Membuat mahasiswa                |
|  5  | `GET`    | `/students/{id}`    | Detail mahasiswa                 |
|  6  | `PUT`    | `/students/{id}`    | Mengubah mahasiswa               |
|  7  | `DELETE` | `/students/{id}`    | Soft delete mahasiswa            |
|  8  | `GET`    | `/courses`          | Daftar mata kuliah               |
|  9  | `POST`   | `/enrollments`      | Mengambil mata kuliah            |
|  10 | `DELETE` | `/enrollments/{id}` | Membatalkan enrollment sendiri   |

---

# Access Matrix

Authorization diuji untuk role `admin`, `mahasiswa`, dan request tanpa authentication.

| Endpoint                       | Admin | Mahasiswa | Public |
| :----------------------------- | :---: | :-------: | :----: |
| `POST /auth/login`             | `200` |   `200`   |  `200` |
| `GET /auth/me`                 | `200` |   `200`   |  `401` |
| `GET /students`                | `200` |   `403`   |  `401` |
| `POST /students`               | `201` |   `403`   |  `401` |
| `GET /students/{id}` self      | `200` |   `200`   |  `401` |
| `GET /students/{id}` other     | `200` |   `403`   |  `401` |
| `PUT /students/{id}`           | `200` |   `403`   |  `401` |
| `DELETE /students/{id}`        | `204` |   `403`   |  `401` |
| `GET /courses`                 | `200` |   `200`   |  `401` |
| `POST /enrollments`            | `403` |   `201`   |  `401` |
| `DELETE /enrollments/{id}` own | `403` |   `204`   |  `401` |

---

# HTTP Status Code

API membedakan authentication, authorization, validation, dan conflict.

| Status | Makna                                       |
| :----: | :------------------------------------------ |
|  `200` | Request berhasil                            |
|  `201` | Resource berhasil dibuat                    |
|  `204` | Request berhasil tanpa response body        |
|  `401` | Authentication tidak valid / tidak tersedia |
|  `403` | Tidak memiliki hak akses                    |
|  `404` | Resource tidak ditemukan                    |
|  `409` | Conflict dengan data yang sudah ada         |
|  `422` | Request melanggar business rule             |
|  `429` | Rate limit terlampaui                       |
|  `500` | Internal server error                       |

---

# Pengujian

Pengujian dilakukan pada tiga level utama:

```text
Unit Test
    +
End-to-End Test
    +
Manual API Verification
```

---

## Unit Test

Total:

```text
44 cases
44 PASS
```

Test mencakup:

| Package       | Test                                         |
| :------------ | :------------------------------------------- |
| `app/service` | `TestBatasSKS`                               |
| `app/service` | `TestTotalSKS`                               |
| `app/service` | `TestParseTahunAkademik_Valid`               |
| `app/service` | `TestParseTahunAkademik_Invalid`             |
| `app/service` | `TestCanEnrollSKS`                           |
| `helper`      | `TestValidateStruct_LoginRequest`            |
| `helper`      | `TestValidateStruct_CreateStudentRequest`    |
| `helper`      | `TestValidateStruct_CreateEnrollmentRequest` |

Business rules yang bersifat pure function dapat diuji tanpa database atau Fiber.

Keuntungannya:

* Cepat.
* Deterministik.
* Tidak bergantung pada environment server.
* Mudah mengisolasi kesalahan.

44 test case selesai dalam waktu kurang dari 0.5 detik.

---

# End-to-End Testing

Automated E2E test menggunakan:

```text
scripts/test-uts.ps1
```

Total:

```text
31 assertions
31 PASS
0 FAIL
```

Test mencakup:

1. Login.
2. Authorization matrix.
3. GET courses.
4. Business rules enrollment.
5. DELETE enrollment.
6. Soft delete student.
7. Login rate limiting.

Script dibuat idempotent sehingga dapat dijalankan berulang kali.

---

# Business Rule Testing

## Batas SKS

Tiga tier diuji:

```text
IPK >= 3.00
    -> 24 SKS

2.50 <= IPK < 3.00
    -> 21 SKS

IPK < 2.50
    -> 18 SKS
```

Ketika batas terlampaui:

```text
422 Unprocessable Entity
```

Response memberikan informasi mengenai batas dan sisa SKS.

---

## Duplicate Enrollment

Request kedua untuk kombinasi:

```text
student_id
course_id
tahun_akademik
```

yang sama menghasilkan:

```text
409 Conflict
```

Database juga memiliki unique constraint sebagai lapisan pertahanan kedua.

---

## Full Quota

Course yang sudah penuh:

```text
422 Unprocessable Entity
```

Pemeriksaan dilakukan dengan row locking untuk mencegah race condition.

---

## Ownership

Mahasiswa yang mencoba mengakses student lain:

```text
403 Forbidden
```

Mahasiswa yang mengakses datanya sendiri:

```text
200 OK
```

---

# Pengujian Soft Delete

Skenario:

```text
DELETE student
        |
        v
deleted_at terisi
        |
        +---- tidak muncul di GET /students
        |
        +---- GET /students/{id} → 404
        |
        +---- login → 401
```

Hasil pengujian menunjukkan bahwa student yang telah di-soft delete tidak dapat login kembali menggunakan credential yang valid.

---

# Pengujian Rate Limiting

Skenario brute-force sederhana:

```text
5 login gagal
    |
    v
401 Unauthorized

login gagal ke-6
    |
    v
429 Too Many Requests
```

Header:

```text
Retry-After: 60
```

diverifikasi pada response.

---

# Environment

| Komponen       | Versi / Teknologi             |
| :------------- | :---------------------------- |
| Language       | Go 1.27                       |
| Framework      | Fiber v2                      |
| Database       | PostgreSQL 16                 |
| Driver         | `pgx/v5`                      |
| JWT            | `golang-jwt/jwt/v5`           |
| Password       | `golang.org/x/crypto/bcrypt`  |
| Validation     | `go-playground/validator/v10` |
| Logging        | `lumberjack.v2`               |
| Shell          | PowerShell                    |
| API Testing    | `curl.exe`                    |
| Manual Testing | Postman                       |

---

# Menjalankan Project

## 1. Masuk ke folder UTS

```bash
cd uts
```

## 2. Install dependency

```bash
go mod tidy
```

## 3. Konfigurasi environment

Gunakan file `.env.example` sebagai referensi.

Pastikan PostgreSQL tersedia dan konfigurasi database telah sesuai.

## 4. Jalankan migration

Jalankan seluruh migration sesuai urutan pada project.

## 5. Jalankan seeder

Seeder digunakan untuk menyediakan:

```text
1 admin
20 mahasiswa
10 mata kuliah
```

## 6. Jalankan server

```bash
go run .
```

Server tersedia pada:

```text
http://localhost:3000
```

API:

```text
http://localhost:3000/api/v1
```

---

# Menjalankan Unit Test

```bash
go test ./... -v
```

Expected:

```text
44 PASS
```

---

# Menjalankan E2E Test

PowerShell:

```powershell
.\scripts\test-uts.ps1
```

Expected:

```text
31 assertions
31 PASS
```

---

# Struktur Project

Struktur utama:

```text
uts/
├── app/
│   ├── model/
│   ├── repository/
│   └── service/
├── config/
├── database/
├── helper/
├── middleware/
├── route/
├── migrations/
├── scripts/
│   └── test-uts.ps1
├── main.go
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

---

# Design Decisions

## Transaction

Digunakan ketika satu operasi melibatkan beberapa perubahan database yang harus berhasil secara atomik.

Contoh:

```text
POST /students
```

---

## Row Locking

Digunakan ketika beberapa request dapat berkompetisi terhadap resource yang sama.

Contoh:

```text
POST /enrollments
```

khususnya pada pemeriksaan quota.

---

## Soft Delete

Dipilih karena student merupakan data akademik yang memiliki nilai historis.

Konsekuensinya adalah setiap query harus memperhatikan:

```text
deleted_at IS NULL
```

---

## Ownership di Service

Ownership membutuhkan informasi tentang record target sehingga pemeriksaannya ditempatkan pada service, bukan middleware.

---

## Business Rules sebagai Pure Function

Fungsi seperti:

```text
BatasSKS
TotalSKS
CanEnrollSKS
ParseTahunAkademik
```

dipisahkan dari database dan HTTP.

Hal ini memungkinkan unit test yang cepat dan deterministik.

---

# Keterbatasan

Satu keterbatasan masih tersisa pada authentication.

JWT bersifat stateless.

Jika sebuah token diterbitkan sebelum student di-soft delete:

```text
Token
expires_in = 900 seconds
```

token tersebut tetap valid hingga expiration.

Artinya terdapat window maksimal sekitar:

```text
15 menit
```

antara soft delete dan expiration token.

Beberapa solusi yang memungkinkan:

1. Token blocklist menggunakan Redis.
2. `RequireAuth` memeriksa status user ke database pada setiap request.
3. Memperpendek lifetime access token.

Untuk UTS ini, keterbatasan tersebut **didokumentasikan**, bukan diperbaiki, karena solusi tersebut berada di luar lingkup spesifikasi dan menambah kompleksitas/dependency.

---

# Hasil Akhir

Implementasi UTS menghasilkan:

```text
10 endpoint
        |
        +---- 4 business rules
        |
        +---- transaction
        |
        +---- row locking
        |
        +---- soft delete
        |
        +---- login guard
        |
        +---- rate limiting
        |
        +---- 44 unit tests
        |
        +---- 31 E2E assertions
```

Hasil pengujian:

```text
44 / 44 unit tests       PASS
31 / 31 E2E assertions   PASS
```

Dua bug penting juga ditemukan dan diperbaiki:

```text
1. UPDATE ... RETURNING dengan alias
2. Soft-deleted student masih dapat login
```

Satu limitation terkait JWT stateless didokumentasikan secara eksplisit.

---

# Repositori

Source code UTS:

`https://github.com/vxpal3n/go-workspace/tree/main/uts`

---

# Sumber Bantuan

* Dokumentasi Go.
* Dokumentasi Fiber v2.
* Dokumentasi PostgreSQL.
* Dokumentasi `pgx/v5`.
* Dokumentasi JWT.
* Dokumentasi bcrypt.
* Dokumentasi `go-playground/validator`.
* Dokumentasi HTTP status code.

Beberapa bagian implementasi dan dokumentasi dibantu oleh alat bantu AI untuk debugging dan penyusunan struktur, namun logika utama disesuaikan dengan kebutuhan dan spesifikasi UTS.
