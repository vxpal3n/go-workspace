# Tugas 05 – Authentication & Security

Implementasi Modul 5 untuk mata kuliah **Pemrograman Backend Lanjut (SIP375)**.

Modul ini merupakan pengembangan dari API Students pada Modul 4. Fokus utama bukan menambah fitur CRUD baru, tetapi **menambahkan sistem authentication dan security** ke dalam aplikasi yang sudah menerapkan Clean Architecture.

Perubahan utama meliputi registrasi akun student, login berbasis NIM dan password, penerapan JWT sebagai access token, mekanisme refresh token dengan rotation, penyimpanan refresh token dalam bentuk hash, middleware `RequireAuth` untuk melindungi endpoint, serta beberapa lapisan pengamanan tambahan seperti bcrypt, rate limiting, CORS terbatas, body limit, dan mitigasi timing attack.

## Tujuan

Modul 5 merupakan pengembangan dari API Students pada Modul 4 dengan menambahkan sistem **authentication dan security** menggunakan `Student` sebagai entitas autentikasi.

Implementasi pada modul ini mencakup:

* Registrasi akun student.
* Login menggunakan NIM dan password.
* Password hashing menggunakan **bcrypt**.
* Access token menggunakan **JWT**.
* Refresh token dengan mekanisme **rotation**.
* Penyimpanan refresh token dalam bentuk hash.
* Middleware `RequireAuth` untuk melindungi endpoint.
* Rate limiting pada endpoint login.
* Pembatasan origin menggunakan CORS.
* Pembatasan ukuran request body menggunakan `BodyLimit`.
* Validasi `JWT_SECRET` sebelum server dijalankan.
* Pencegahan **mass assignment** terhadap field `role`.
* Mitigasi timing attack pada proses login.

Struktur aplikasi tetap mempertahankan pendekatan **Clean Architecture** dari Modul 4.

---

## Teknologi

| Komponen         | Teknologi                             |
| :--------------- | :------------------------------------ |
| Bahasa           | Go                                    |
| Framework        | Fiber v2                              |
| Database         | PostgreSQL                            |
| Database Driver  | pgx/v5 + pgxpool                      |
| Authentication   | JWT                                   |
| Password Hashing | bcrypt                                |
| Refresh Token    | Cryptographically Secure Random Token |
| Hashing Token    | SHA-256                               |
| Rate Limiting    | Fiber Limiter                         |
| CORS             | Fiber CORS                            |
| Security Headers | Fiber Helmet                          |
| Logging          | `log/slog` + lumberjack               |
| Testing          | Go testing                            |

---

## Struktur Folder

```text
tugas05/
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── main.go
├── README.md
├── app/
│   ├── model/
│   │   ├── student.go
│   │   └── auth.go
│   ├── repository/
│   │   ├── student_repository.go
│   │   └── token_repository.go
│   └── service/
│       ├── student_service.go
│       ├── student_rules.go
│       ├── auth_service.go
│       ├── auth_rules.go
│       └── auth_rules_test.go
├── helper/
│   ├── response.go
│   ├── request.go
│   ├── security.go
│   ├── jwt.go
│   └── context.go
├── middleware/
│   ├── middleware.go
│   └── auth.go
├── route/
│   └── route.go
├── config/
│   ├── env.go
│   ├── logger.go
│   └── app.go
├── database/
│   └── postgres.go
├── migrations/
│   ├── 001_create_students.sql
│   └── 002_auth.sql
└── logs/
    └── app.log
```

Folder `logs/` tidak disimpan dalam repository karena telah dimasukkan ke `.gitignore`.

---

## Konsep Authentication

Pada Modul 5, entitas `Student` tidak hanya digunakan sebagai business entity, tetapi juga menjadi **entitas autentikasi**.

Field authentication yang ditambahkan:

| Field       | Fungsi                                         |
| :---------- | :--------------------------------------------- |
| `email`     | Identitas email student                        |
| `password`  | Password yang telah di-hash menggunakan bcrypt |
| `role`      | Role pengguna                                  |
| `is_active` | Status akun                                    |

Password menggunakan tag JSON `json:"-"` sehingga tidak pernah dikirimkan dalam response API.

Role tidak dapat dikontrol melalui request. Pada proses registrasi, server secara otomatis menetapkan role:

```text
user
```

Hal ini digunakan untuk mencegah **mass assignment**, sehingga client tidak dapat mendaftarkan dirinya sebagai administrator.

---

## Database

Migration kedua digunakan untuk menambahkan kebutuhan authentication ke database.

### Migration

```text
migrations/
├── 001_create_students.sql
└── 002_auth.sql
```

Migration `002_auth.sql` melakukan beberapa perubahan:

1. Menghapus data student lama karena data tersebut belum memiliki password.
2. Menambahkan kolom `email`.
3. Menambahkan kolom `password`.
4. Menambahkan kolom `role`.
5. Membuat unique index case-insensitive pada email.
6. Membuat tabel `refresh_tokens`.
7. Membuat index untuk `student_id` dan `token_hash`.

Struktur tabel `refresh_tokens`:

| Kolom        | Tipe          | Keterangan                |
| :----------- | :------------ | :------------------------ |
| `id`         | `BIGSERIAL`   | Primary key               |
| `student_id` | `INTEGER`     | Foreign key ke `students` |
| `token_hash` | `TEXT`        | Hash refresh token        |
| `expires_at` | `TIMESTAMPTZ` | Waktu kedaluwarsa         |
| `revoked_at` | `TIMESTAMPTZ` | Waktu token dicabut       |
| `created_at` | `TIMESTAMPTZ` | Waktu token dibuat        |

Refresh token **tidak disimpan dalam bentuk plaintext**, melainkan hash SHA-256.

### Menjalankan Migration

Pastikan database `db_students` sudah tersedia, kemudian jalankan:

```bash
psql -U postgres -d db_students -f tugas05/migrations/001_create_students.sql
psql -U postgres -d db_students -f tugas05/migrations/002_auth.sql
```

Verifikasi struktur tabel:

```bash
psql -U postgres -d db_students -c "\d students"
psql -U postgres -d db_students -c "\d refresh_tokens"
```

---

## Environment Variables

Buat file `.env` di dalam folder `tugas05/`.

```env
APP_PORT=3000
APP_NAME=Praktikum Backend Lanjut - Tugas05

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=db_students
DB_SSLMODE=disable
DB_MAX_CONNS=10

JWT_SECRET=
JWT_ISSUER=praktikum-backend
JWT_ACCESS_TTL_MINUTES=15
JWT_REFRESH_TTL_DAYS=7

ALLOWED_ORIGINS=http://localhost:5173

LOG_LEVEL=info
```

Contoh konfigurasi tersedia pada:

```text
tugas05/.env.example
```

`JWT_SECRET` wajib memiliki panjang minimal **32 karakter**.

Generate secret menggunakan:

```bash
openssl rand -hex 32
```

File `.env` tidak boleh dimasukkan ke repository.

Verifikasi:

```bash
cd tugas05
git check-ignore .env
```

Jika menghasilkan:

```text
.env
```

berarti file telah berhasil di-ignore.

---

## Security Implementation

### 1. Password Hashing

Password tidak disimpan secara langsung ke database.

Proses registrasi:

```text
Plain Password
      |
      v
   bcrypt
      |
      v
Password Hash
      |
      v
 PostgreSQL
```

Implementasi menggunakan bcrypt dengan cost `12`.

---

### 2. JWT Access Token

Access token digunakan untuk mengakses endpoint yang membutuhkan autentikasi.

Payload token membawa informasi minimal:

```json
{
  "student_id": 1,
  "nim": "S001",
  "role": "user"
}
```

Token menggunakan:

```text
Algorithm : HS256
Issuer    : praktikum-backend
Expiration: 15 menit
```

Middleware akan menolak token yang:

* Tidak valid.
* Telah dimodifikasi.
* Kedaluwarsa.
* Menggunakan algoritma yang tidak diharapkan.
* Tidak memiliki expiration claim.

Implementasi juga melakukan pemeriksaan terhadap algoritma JWT untuk mencegah **algorithm confusion**.

---

### 3. Refresh Token

Refresh token dibuat menggunakan random bytes yang aman secara kriptografis.

Alurnya:

```text
Login
  |
  +---- Access Token
  |
  +---- Refresh Token
             |
             v
        SHA-256 Hash
             |
             v
       PostgreSQL
```

Refresh token memiliki masa berlaku default:

```text
7 hari
```

Refresh token lama akan dicabut ketika digunakan untuk mendapatkan token baru.

Dengan demikian, mekanisme yang digunakan adalah **refresh token rotation**.

---

### 4. Timing Attack Mitigation

Jika NIM tidak ditemukan, aplikasi tetap menjalankan operasi bcrypt menggunakan dummy hash.

Tujuannya agar perbedaan waktu response antara:

```text
NIM tidak ditemukan
```

dan:

```text
NIM ditemukan tetapi password salah
```

tidak terlalu mudah digunakan untuk mengetahui apakah sebuah NIM terdaftar.

Pesan error untuk kedua kondisi juga dibuat sama:

```text
NIM atau password salah
```

---

### 5. Rate Limiting

Endpoint login memiliki batas:

```text
5 request / 1 menit / IP
```

Jika batas terlampaui, API memberikan:

```text
429 Too Many Requests
```

serta header:

```text
Retry-After: 60
```

Rate limiter digunakan untuk mengurangi risiko brute-force login.

---

### 6. CORS

CORS tidak lagi menggunakan konfigurasi terbuka.

Origin dikontrol melalui:

```env
ALLOWED_ORIGINS=http://localhost:5173
```

Method yang diizinkan:

```text
GET
POST
PUT
PATCH
DELETE
OPTIONS
```

Header `Authorization` juga secara eksplisit diizinkan untuk kebutuhan Bearer Token.

---

### 7. Body Limit

Ukuran request body dibatasi hingga:

```text
1 MB
```

Konfigurasi ini digunakan untuk mengurangi risiko penggunaan resource secara berlebihan melalui payload berukuran besar.

---

## API Contract

Base URL:

```text
http://localhost:3000/api/v1
```

### Health Check

| Method | Endpoint  | Auth | Fungsi                               |
| :----- | :-------- | :--: | :----------------------------------- |
| `GET`  | `/health` |  No  | Mengecek server dan koneksi database |

Response sukses:

```json
{
  "success": true,
  "message": "server dan database berjalan"
}
```

---

## Authentication Endpoints

| Method | Endpoint         | Auth | Fungsi                                      |
| :----- | :--------------- | :--: | :------------------------------------------ |
| `POST` | `/auth/register` |  No  | Membuat akun student                        |
| `POST` | `/auth/login`    |  No  | Login dan mendapatkan token                 |
| `POST` | `/auth/refresh`  |  No  | Memperbarui access token                    |
| `POST` | `/auth/logout`   |  No  | Mencabut refresh token                      |
| `GET`  | `/auth/me`       |  Yes | Mengambil profil pengguna yang sedang login |

Endpoint authentication didaftarkan pada route `/api/v1/auth`.

---

## Student Endpoints

Mulai Modul 5, seluruh endpoint student membutuhkan access token.

| Method   | Endpoint        | Auth | Fungsi                         |
| :------- | :-------------- | :--: | :----------------------------- |
| `GET`    | `/students`     |  Yes | Daftar student                 |
| `GET`    | `/students/:id` |  Yes | Detail student                 |
| `PUT`    | `/students/:id` |  Yes | Mengganti data student         |
| `PATCH`  | `/students/:id` |  Yes | Mengubah sebagian data student |
| `DELETE` | `/students/:id` |  Yes | Menghapus student              |

Endpoint:

```text
POST /students
```

dihapus.

Pembuatan student sekarang dilakukan melalui:

```text
POST /auth/register
```

Hal ini mencegah adanya dua jalur pembuatan akun dengan aturan authentication yang berbeda.

---

## Contoh Request

### Register

```bash
curl -i -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"nim":"S001","name":"Thaariq","email":"thaariq@example.com","grade":85.5,"password":"rahasia123"}'
```

Response yang diharapkan:

```text
201 Created
```

---

### Login

```bash
curl -X POST http://localhost:3000/api/v1/auth/login -H "Content-Type: application/json" -d '{"nim":"S001","password":"rahasia123"}'
```

Response akan berisi:

```json
{
  "success": true,
  "message": "login berhasil",
  "data": {
    "access_token": "...",
    "refresh_token": "...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

---

### Mengakses Profil

Gunakan access token dari hasil login:

```bash
curl -i http://localhost:3000/api/v1/auth/me -H "Authorization: Bearer ACCESS_TOKEN"
```

---

### Mengakses Student

```bash
curl http://localhost:3000/api/v1/students -H "Authorization: Bearer ACCESS_TOKEN"
```

---

### Refresh Token

```bash
curl -X POST http://localhost:3000/api/v1/auth/refresh -H "Content-Type: application/json" -d '{"refresh_token":"REFRESH_TOKEN"}'
```

Refresh token yang sama tidak dapat digunakan kembali setelah berhasil di-rotate.

---

### Logout

```bash
curl -X POST http://localhost:3000/api/v1/auth/logout -H "Content-Type: application/json" -d '{"refresh_token":"REFRESH_TOKEN"}'
```

---

## Validasi Authentication

Password pada proses registrasi harus:

* Minimal 8 karakter.
* Mengandung huruf.
* Mengandung angka.
* Tidak menggunakan password umum tertentu.

Contoh password yang ditolak:

```text
abc1
rahasiaku
12345678
password1
```

Business rules authentication dibuat sebagai fungsi terpisah sehingga dapat diuji tanpa menjalankan Fiber maupun database.

---

## Testing

Unit test untuk authentication rules dapat dijalankan dengan:

```bash
cd tugas05
go test ./app/service -v -run TestValidate
```

Validasi yang diuji meliputi:

* Register berhasil.
* Password terlalu pendek.
* Password tanpa angka.
* Password tanpa huruf.
* Password terlalu umum.
* Format email tidak valid.
* Grade berada di luar rentang 0–100.

Untuk memastikan seluruh project dapat dikompilasi dan lolos static analysis:

```bash
go build ./...
go vet ./...
```

---

## Skenario Verifikasi

Implementasi Modul 5 diuji menggunakan beberapa skenario utama:

| No. | Skenario                        | Expected Result                           |
| :-: | :------------------------------ | :---------------------------------------- |
|  1  | Register dengan data valid      | `201 Created`                             |
|  2  | Register dengan password lemah  | `422 Unprocessable Entity`                |
|  3  | Memeriksa password di database  | Password tersimpan sebagai bcrypt hash    |
|  4  | Akses `/students` tanpa token   | `401 Unauthorized`                        |
|  5  | Login dengan password salah     | `401 Unauthorized`                        |
|  6  | Login dengan NIM tidak ada      | `401 Unauthorized` dengan pesan yang sama |
|  7  | Login dengan credential benar   | `200 OK` + access/refresh token           |
|  8  | Access token valid              | Request berhasil                          |
|  9  | Access token dimodifikasi       | `401 Unauthorized`                        |
|  10 | Refresh token digunakan kembali | `401 Unauthorized`                        |
|  11 | Percobaan login berlebihan      | `429 Too Many Requests`                   |
|  12 | Register dengan `role: admin`   | Role tetap `user`                         |

Skenario tersebut mencakup authentication, authorization middleware, token security, refresh token rotation, rate limiting, dan mass assignment protection.

---

## Menjalankan Aplikasi

### 1. Masuk ke folder project

```bash
cd tugas05
```

### 2. Install dependency

```bash
go mod tidy
```

### 3. Konfigurasi environment

Buat `.env` berdasarkan `.env.example` dan isi:

```env
DB_PASSWORD=PASSWORD_POSTGRES
JWT_SECRET=JWT_SECRET_MINIMAL_32_KARAKTER
```

### 4. Jalankan migration

```bash
psql -U postgres -d db_students -f migrations/001_create_students.sql
psql -U postgres -d db_students -f migrations/002_auth.sql
```

### 5. Jalankan server

```bash
go run .
```

Server berjalan pada:

```text
http://localhost:3000
```

---

## Konvensi Commit

Pengerjaan Modul 5 menggunakan pendekatan **Conventional Commits** untuk mendokumentasikan progres implementasi secara bertahap.

Kategori yang digunakan:

| Prefix     | Penggunaan                              |
| :--------- | :-------------------------------------- |
| `feat`     | Fitur baru                              |
| `fix`      | Perbaikan bug                           |
| `refactor` | Restrukturisasi tanpa mengubah perilaku |
| `test`     | Pengujian                               |
| `docs`     | Dokumentasi                             |
| `chore`    | Dependency, konfigurasi, dan tooling    |

---

## Catatan

* `JWT_SECRET` wajib memiliki minimal 32 karakter.
* File `.env` tidak boleh di-commit.
* Password tidak pernah dikembalikan dalam response JSON.
* Refresh token disimpan dalam database sebagai hash.
* Role ditentukan oleh server pada saat registrasi.
* Endpoint `/students` membutuhkan access token.
* Endpoint `POST /students` tidak tersedia pada Modul 5.
* Refresh token menggunakan mekanisme rotation.
* Login memiliki rate limiter berdasarkan IP.
* Request body dibatasi hingga 1 MB.
* CORS menggunakan daftar origin yang dikonfigurasi melalui environment variable.

---

## Repositori

Kode sumber Modul 5:

`https://github.com/vxpal3n/go-workspace/tree/main/tugas05`

---

## Sumber Bantuan

* Dokumentasi Go
* Dokumentasi Fiber v2
* Dokumentasi PostgreSQL
* Dokumentasi JWT
* Dokumentasi bcrypt

Beberapa bagian implementasi dan dokumentasi dibantu oleh alat bantu AI untuk debugging dan penyusunan struktur, sedangkan logika utama disesuaikan dengan kebutuhan praktikum.
