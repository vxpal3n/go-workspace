# Go Workspace – Pemrograman Backend Lanjut (SIP375)

Repositori ini merupakan kumpulan tugas praktikum mata kuliah **Pemrograman Backend Lanjut (SIP375)**.

Setiap modul dikerjakan secara bertahap dengan mengembangkan API Students menggunakan Go, mulai dari pengenalan sintaks, REST API, database, Clean Architecture, authentication, authorization, advanced API design, hingga implementasi UTS dalam bentuk SIAKAD Mini.

**Author:** Valentino Chandra
**Semester:** Gasal 2026/2027
**Dosen:** Rachman Sinatriya Marjianto, B.Eng, M.Sc.

---

## Struktur Repositori

| Folder                  | Modul   | Deskripsi                                 |
| :---------------------- | :------ | :---------------------------------------- |
| [`tugas01/`](./tugas01) | Modul 1 | Persiapan lingkungan dan sintaks dasar Go |
| [`tugas02/`](./tugas02) | Modul 2 | REST API dan HTTP Deep Dive               |
| [`tugas03/`](./tugas03) | Modul 3 | Database dan Repository Pattern           |
| [`tugas04/`](./tugas04) | Modul 4 | Clean Architecture                        |
| [`tugas05/`](./tugas05) | Modul 5 | Authentication & Security                 |
| [`tugas06/`](./tugas06) | Modul 6 | Authorization & Role-Based Access Control |
| [`tugas07/`](./tugas07) | Modul 7 | Advanced API Design                       |
| [`uts/`](./uts)         | UTS     | SIAKAD Mini – RESTful API Backend         |

---

# Modul 1 – Persiapan & Sintaks Go

Folder:

```text
tugas01/
```

Modul pertama berfokus pada persiapan lingkungan pengembangan Go dan pemahaman sintaks dasar.

Materi yang diterapkan meliputi:

* Variabel.
* Slice.
* Map.
* Pointer.
* Struct.
* Method.
* Hello World menggunakan Fiber.

Dokumentasi modul:

[`tugas01/README.md`](./tugas01/README.md)

---

# Modul 2 – REST API & HTTP Deep Dive

Folder:

```text
tugas02/
```

Modul kedua membangun REST API untuk entitas `Student` menggunakan Go dan Fiber v2.

Fitur utama:

* RESTful endpoint.
* HTTP method GET, POST, PUT, PATCH, dan DELETE.
* HTTP status code yang sesuai.
* Paginasi.
* Search.
* Sorting.
* Filtering.
* Validasi request.
* Error handling.
* Perbedaan PUT dan PATCH.

Dokumentasi lengkap:

[`tugas02/README.md`](./tugas02/README.md)

---

# Modul 3 – Database & Repository Pattern

Folder:

```text
tugas03/
```

Modul ketiga mengembangkan API Students dengan penyimpanan permanen menggunakan **PostgreSQL**.

Konsep yang diterapkan:

* PostgreSQL.
* `pgx/v5` dan `pgxpool`.
* Database migration.
* Repository Pattern.
* Parameterized query.
* Filtering dan search melalui SQL.
* Sorting dengan whitelist.
* Pagination pada database.
* Translasi database error menjadi HTTP response.

Struktur repository digunakan untuk memisahkan logika penyimpanan data dari komponen API.

Dokumentasi modul:

[`tugas03/README.md`](./tugas03/README.md)

---

# Modul 4 – Clean Architecture

Folder:

```text
tugas04/
```

Modul keempat melakukan restrukturisasi API Students ke dalam pendekatan **Clean Architecture**.

Fokus utama:

* Dependency Rule.
* Pemisahan entities, use cases, interface adapters, dan frameworks.
* Pemisahan business rules dari framework.
* Unit testing business rules.
* Dependency injection.
* Structured logging.
* Pemeriksaan layer leakage.

Business rules dipindahkan ke file khusus agar tidak bergantung pada Fiber maupun database.

Dokumentasi modul:

[`tugas04/README.md`](./tugas04/README.md)

---

# Modul 5 – Authentication & Security

Folder:

```text
tugas05/
```

Modul kelima menambahkan lapisan **authentication dan security** pada API.

Fitur utama:

* JWT authentication.
* Password hashing menggunakan bcrypt.
* Login.
* Authenticated user context.
* Authentication middleware.
* Rate limiting.
* CORS.
* Request body limit.
* Security helper.
* Token repository.
* Auth service.
* Unit test authentication rules.

Modul ini menjadi fondasi authorization yang dikembangkan pada Modul 6.

Dokumentasi modul:

[`tugas05/README.md`](./tugas05/README.md)

---

# Modul 6 – Authorization & Role-Based Access Control

Folder:

```text
tugas06/
```

Setelah authentication dari Modul 5 selesai, Modul 6 berfokus pada pertanyaan berikut:

```text
"Pengguna ini boleh melakukan apa?"
```

Modul ini menerapkan **Role-Based Access Control (RBAC)** dan ownership.

Fitur utama:

* Role `admin`, `staff`, dan `user`.
* Permission.
* `roles`.
* `permissions`.
* `role_permissions`.
* Permission middleware.
* Ownership check.
* `owner_id`.
* Fail closed.
* Pencegahan `owner_id` spoofing.
* Business rule untuk mencegah penghapusan akun sendiri.
* Business rule untuk mencegah perubahan role diri sendiri.
* Unit testing authorization rules.

Pembagian tanggung jawab:

```text
┌─────────────────────────────────────┐
│  Middleware                         │
│  ─────────────────────────────────  │
│  Role / Permission decision         │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Service                            │
│  ─────────────────────────────────  │
│  Ownership / Data-dependent         │
│  decision                           │
└─────────────────────────────────────┘
```

Dokumentasi modul:

[`tugas06/README.md`](./tugas06/README.md)

---

# Modul 7 – Advanced API Design

Folder:

```text
tugas07/
```

Modul ketujuh berfokus pada penyempurnaan API melalui debugging dan peningkatan desain.

Modul ini menemukan dan memperbaiki:

```text
┌─────────────────────────────────────┐
│  3 compiler errors                  │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  6 behavioral bugs                  │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  9 bugs total                       │
└─────────────────────────────────────┘
```

Fitur dan konsep utama:

* Declarative validation.
* Custom validator.
* Centralized error handling.
* `AppError`.
* Stable error code.
* `request_id`.
* Correct HTTP status mapping.
* Cursor / keyset pagination.
* PostgreSQL index untuk pagination.
* `EXPLAIN ANALYZE`.
* Content negotiation.
* JSON response.
* CSV response.
* Logging berdasarkan status HTTP.
* Error handling untuk unknown error.

Salah satu fokus penting adalah memastikan API tidak hanya “berjalan”, tetapi juga memiliki kontrak error, pagination, validation, dan response format yang konsisten.

Dokumentasi modul:

[`tugas07/README.md`](./tugas07/README.md)

---

# UTS – SIAKAD Mini RESTful API Backend

Folder:

```text
uts/
```

UTS menggabungkan berbagai konsep yang telah dibangun sepanjang modul sebelumnya ke dalam satu aplikasi **SIAKAD Mini**.

Domain utama:

```text
Students
Courses
Enrollments / KRS
```

API menyediakan **10 endpoint** dengan dua role utama:

```text
admin
mahasiswa
```

---

## Konsep yang Digabungkan

UTS mengintegrasikan:

* Clean Architecture.
* PostgreSQL.
* Repository Pattern.
* JWT authentication.
* Role-based authorization.
* Ownership.
* Declarative validation.
* Business rules.
* Database transaction.
* Row-level locking.
* Soft delete.
* Login guard.
* Rate limiting.
* Unit testing.
* End-to-end testing.

---

## Business Rules

Empat business rule utama:

### Batas SKS

```text
┌─────────────────────────────────────┐
│  IPK >= 3.00                        │
│  ─────────────────────────────────  │
│  → 24 SKS                           │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  2.50 <= IPK < 3.00                 │
│  ─────────────────────────────────  │
│  → 21 SKS                           │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  IPK < 2.50                         │
│  ─────────────────────────────────  │
│  → 18 SKS                           │
└─────────────────────────────────────┘
```

### Duplicate Enrollment

Kombinasi:

```text
┌─────────────────────────────────────┐
│  student_id                         │
│  course_id                          │
│  tahun_akademik                     │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  harus unik                         │
└─────────────────────────────────────┘
```

### Full Quota

Enrollment tidak dapat dilakukan jika quota mata kuliah sudah penuh.

Pemeriksaan dilakukan menggunakan:

```sql
SELECT ... FOR UPDATE
```

untuk mencegah race condition.

### Ownership

Mahasiswa hanya dapat mengakses data yang menjadi miliknya.

---

## Transaction

`POST /students` menggunakan transaction karena proses pembuatan student melibatkan:

```text
┌─────────────────────────────────────┐
│  users                              │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  students                           │
└─────────────────────────────────────┘
```

Jika salah satu operasi gagal:

```text
┌─────────────────────────────────────┐
│  ROLLBACK                           │
└─────────────────────────────────────┘
```

---

## Soft Delete

Student menggunakan:

```text
deleted_at
```

sehingga data historis tetap tersedia.

Student yang telah di-soft delete:

```text
┌─────────────────────────────────────┐
│  tidak muncul pada query aktif      │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  tidak dapat login                  │
└─────────────────────────────────────┘
```

---

## Testing

UTS memiliki dua automated test layer utama:

```text
┌─────────────────────────────────────┐
│  44 Unit Test Cases                 │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  31 E2E Assertions                  │
└─────────────────────────────────────┘
```

Hasil akhir:

```text
┌─────────────────────────────────────┐
│  44 / 44  PASS                      │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  31 / 31  PASS                      │
└─────────────────────────────────────┘
```

Rate limiting juga diverifikasi:

```text
┌─────────────────────────────────────┐
│  5 failed login attempts            │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  401 Unauthorized                   │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  6th attempt                        │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  429 Too Many Requests              │
└─────────────────────────────────────┘
```

---

# Progression

Keseluruhan repository dirancang sebagai progression pembelajaran:

```text
┌─────────────────────────────────────┐
│  Modul 1 — Go Fundamentals          │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Modul 2 — REST API & HTTP          │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Modul 3 — PostgreSQL & Repository  │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Modul 4 — Clean Architecture       │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Modul 5 — Authentication &         │
│            Security                 │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Modul 6 — Authorization & RBAC     │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Modul 7 — Advanced API Design      │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  UTS — SIAKAD Mini                  │
└─────────────────────────────────────┘
```

Setiap modul memperluas fondasi modul sebelumnya.

---

# Teknologi Utama

| Teknologi     | Penggunaan               |
| :------------ | :----------------------- |
| Go            | Bahasa pemrograman utama |
| Fiber v2      | HTTP framework           |
| PostgreSQL    | Relational database      |
| pgx/v5        | PostgreSQL driver        |
| JWT           | Authentication           |
| bcrypt        | Password hashing         |
| Validator v10 | Declarative validation   |
| slog          | Structured logging       |
| lumberjack    | Log rotation             |
| curl          | API testing              |
| PowerShell    | Automated E2E testing    |
| Postman       | Manual API testing       |

---

# Cara Menjalankan

Setiap modul merupakan project Go yang dapat dijalankan secara independen.

Contoh:

```bash
cd tugas07
go mod tidy
go run .
```

Untuk UTS:

```bash
cd uts
go mod tidy
go run .
```

Perintah testing:

```bash
go test ./... -v
```

Detail konfigurasi, migration, seeder, endpoint, dan testing tersedia pada README masing-masing modul.

---

# Dokumentasi

Dokumentasi lengkap tersedia pada setiap folder:

| Modul   | Dokumentasi                                |
| :------ | :----------------------------------------- |
| Modul 1 | [`tugas01/README.md`](./tugas01/README.md) |
| Modul 2 | [`tugas02/README.md`](./tugas02/README.md) |
| Modul 3 | [`tugas03/README.md`](./tugas03/README.md) |
| Modul 4 | [`tugas04/README.md`](./tugas04/README.md) |
| Modul 5 | [`tugas05/README.md`](./tugas05/README.md) |
| Modul 6 | [`tugas06/README.md`](./tugas06/README.md) |
| Modul 7 | [`tugas07/README.md`](./tugas07/README.md) |
| UTS     | [`uts/README.md`](./uts/README.md)         |

---

# Riwayat Pengembangan

Repository dikembangkan secara bertahap menggunakan Git.

Setiap modul memiliki riwayat commit tersendiri yang merepresentasikan proses pengembangan fitur, perbaikan bug, refactoring, dan penyempurnaan keamanan.

Struktur ini membuat perkembangan project dapat dilacak dari:

```text
┌─────────────────────────────────────┐
│  Go fundamentals                    │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  REST API                           │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Database                           │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Architecture                       │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Authentication                     │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Authorization                      │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Advanced API Design                │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  SIAKAD Mini                        │
└─────────────────────────────────────┘
```

---

# Sumber Bantuan

* Dokumentasi resmi Go.
* Dokumentasi Fiber v2.
* Dokumentasi PostgreSQL.
* Dokumentasi `pgx/v5`.
* Dokumentasi JWT.
* Dokumentasi bcrypt.
* Dokumentasi `go-playground/validator`.
* RFC dan dokumentasi HTTP terkait.

Beberapa bagian kode dan dokumentasi dibantu oleh alat bantu AI untuk debugging, analisis, dan penyusunan struktur, namun implementasi disesuaikan dengan kebutuhan masing-masing modul dan tugas akademik.

---

# Repositori

Source code:

`https://github.com/vxpal3n/go-workspace`

---

**Go Workspace – Pemrograman Backend Lanjut (SIP375)**
**Semester Gasal 2026/2027**