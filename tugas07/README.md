# Tugas 07 – Advanced API Design

Implementasi Modul 7 untuk mata kuliah **Pemrograman Backend Lanjut (SIP375)**.

Modul ini merupakan tahap penyempurnaan API Students setelah authentication, authorization, dan Clean Architecture diterapkan pada modul-modul sebelumnya.

Fokus Modul 7 bukan sekadar menambahkan endpoint baru, tetapi **memperbaiki kualitas desain API melalui debugging, validasi deklaratif, error handling terpusat, cursor pagination, dan content negotiation**.

Karakteristik utama modul ini adalah proses debugging terhadap **sembilan bug** yang sengaja terdapat pada implementasi awal:

```text
3 compiler errors
+
6 behavioral bugs
=
9 bugs
```

Setiap bug diperbaiki secara bertahap dan diuji ulang sebelum melanjutkan ke bug berikutnya.

---

## Tujuan

Modul 7 memiliki lima fokus utama:

* Menemukan dan memperbaiki sembilan bug pada implementasi API.
* Mengganti validasi manual dengan **declarative validation**.
* Memusatkan error handling melalui `ErrorHandler`.
* Mengganti pagination berbasis offset dengan **cursor/keyset pagination**.
* Menambahkan **content negotiation** untuk JSON dan CSV.

Selain itu, modul ini memperkuat kontrak API melalui:

* Stable error code.
* `request_id`.
* Error response yang seragam.
* Logging dengan status HTTP yang benar.
* HTTP status code yang sesuai secara semantik.

---

## Metodologi Debugging

Debugging dilakukan secara bertahap.

Urutan yang digunakan:

```text
1. Dengarkan compiler
        |
        v
2. Bandingkan komentar dengan kode
        |
        v
3. Bandingkan implementasi dengan teori/specification
        |
        v
4. Uji terhadap acceptance criteria
        |
        v
5. Perbaiki satu bug → test ulang
```

Pendekatan ini penting karena memperbaiki banyak bagian sekaligus membuat sulit menentukan perubahan mana yang benar-benar menyelesaikan masalah.

---

## Teknologi

| Komponen         | Teknologi                   |
| :--------------- | :-------------------------- |
| Bahasa           | Go 1.27                     |
| Framework        | Fiber v2                    |
| Database         | PostgreSQL 16               |
| Database Driver  | pgx/v5                      |
| Authentication   | golang-jwt/jwt/v5           |
| Password Hashing | golang.org/x/crypto/bcrypt  |
| Validation       | go-playground/validator/v10 |
| Logging          | `log/slog` + lumberjack     |
| Shell Testing    | PowerShell + `curl.exe`     |
| Database Testing | `psql` + `EXPLAIN ANALYZE`  |
| Unit Testing     | Go testing                  |

---

# Bug Fixing

## Ringkasan Sembilan Bug

| Bug | Jenis      | Masalah Utama                                            | Perbaikan                                          |
| :-: | :--------- | :------------------------------------------------------- | :------------------------------------------------- |
|  #1 | Compiler   | `AppError.cause` tidak dapat diakses dari package lain   | Menambahkan method `Cause()`                       |
|  #2 | Compiler   | Variabel `status` dideklarasikan tetapi tidak digunakan  | Menggunakan `status` hasil koreksi pada access log |
|  #3 | Compiler   | Field PATCH bukan pointer                                | Mengubah field opsional menjadi pointer            |
|  #4 | Behavioral | Validator `strongpassword` memiliki kondisi terbalik     | Mengubah `!= ""` menjadi `== ""`                   |
|  #5 | Behavioral | Error database tidak dikenal diterjemahkan menjadi `nil` | Default error menjadi `Internal`                   |
|  #6 | Behavioral | Log 4xx dicatat sebagai `ERROR`                          | 4xx → `WARN`, 5xx → `ERROR`                        |
|  #7 | Behavioral | Pagination menggunakan arah `ORDER BY` yang salah        | Menggunakan `created_at DESC, id DESC`             |
|  #8 | Behavioral | Response CSV kosong                                      | Memanggil `writer.Flush()`                         |
|  #9 | Behavioral | Validation error menghasilkan `400`                      | Mengubah menjadi `422`                             |

---

## Bug #1 – Private Error Cause

`AppError` memiliki field:

```text
cause
```

yang sengaja dibuat private agar detail error internal tidak langsung digunakan sebagai response client.

Namun `config` mencoba mengakses field tersebut secara langsung.

Perbaikannya adalah menyediakan method:

```go
func (e *AppError) Cause() error {
    return e.cause
}
```

Dengan demikian:

```text
cause tetap private
        |
        v
Cause() hanya digunakan ketika diperlukan
        |
        v
logging
```

Detail teknis database tetap tidak bocor ke client.

---

## Bug #2 – Status Log Tidak Digunakan

Request logger menghitung status yang benar ketika handler mengembalikan error.

Namun log masih menggunakan:

```text
c.Response().StatusCode()
```

sehingga status yang dicatat dapat berbeda dengan status yang diterima client.

Perbaikannya adalah menggunakan variable:

```text
status
```

yang telah dikoreksi berdasarkan error.

Hasil akhirnya:

```text
Client response status
        =
Access log status
```

---

## Bug #3 – PATCH Request Bukan Pointer

Partial update membutuhkan kemampuan membedakan:

```text
field tidak dikirim
```

dengan:

```text
field dikirim dengan nilai kosong
```

Karena itu field optional pada `PatchStudentRequest` menggunakan pointer:

```text
*string
*float64
*bool
```

Contohnya:

```json
{}
```

berbeda dengan:

```json
{
  "name": ""
}
```

Pointer memungkinkan kedua kondisi tersebut dibedakan.

---

## Bug #4 – Custom Password Validator Terbalik

Validator `strongpassword` memiliki kondisi:

```text
CheckPasswordStrength(...) != ""
```

Padahal fungsi tersebut menghasilkan:

```text
""     → password kuat
pesan  → password lemah
```

Akibatnya logika menjadi terbalik.

Perbaikan:

```text
CheckPasswordStrength(...) == ""
```

Setelah diperbaiki:

```text
password kuat
    -> diterima

password lemah
    -> ditolak
```

---

## Bug #5 – Unknown Error Menjadi 200

Sebelum perbaikan, `translateError` hanya mengenali beberapa error tertentu.

Untuk error lain, fungsi mengembalikan:

```text
nil
```

Akibatnya error database dapat terlihat sebagai:

```text
200 OK
```

dengan response kosong.

Ini sangat berbahaya karena kegagalan sistem terlihat seperti keberhasilan.

Perbaikannya:

```text
unknown error
    |
    v
helper.Internal(err)
    |
    v
500 Internal Server Error
```

Response client menggunakan informasi generik, sedangkan detail error asli hanya dicatat pada log.

---

## Bug #6 – Log Level Tidak Sesuai

Error `4xx` merupakan kesalahan request/client, sedangkan `5xx` menunjukkan masalah server.

Setelah perbaikan:

```text
4xx
  -> WARN
  -> request_rejected

5xx
  -> ERROR
  -> request_failed
```

Pembagian ini membuat log lebih mudah dianalisis.

---

## Bug #7 – Cursor Pagination

Pagination sebelumnya menggunakan konsep offset.

Modul ini menggantinya dengan **cursor/keyset pagination**.

Pasangan kolom yang digunakan:

```text
created_at DESC
id DESC
```

`created_at` digunakan sebagai urutan kronologis.

`id` digunakan sebagai tie-breaker karena timestamp tidak dijamin unik.

Dengan kombinasi tersebut, urutan menjadi:

```text
unik
+
deterministik
```

---

## Index Cursor Pagination

Database memiliki index:

```sql
CREATE INDEX IF NOT EXISTS students_created_at_id_desc_idx
ON students (created_at DESC, id DESC);
```

Arah dan urutan index dibuat sama dengan query:

```text
ORDER BY created_at DESC, id DESC
```

Tujuannya agar PostgreSQL dapat menggunakan index ketika mencari posisi berdasarkan cursor.

---

## Cursor Pagination vs Offset

### Offset

```text
page=100
limit=10
```

Database perlu melewati sejumlah besar record sebelum mendapatkan halaman yang diminta.

Masalah lain muncul ketika data berubah di antara dua request.

Record dapat:

```text
terduplikasi
```

atau:

```text
terlewat
```

### Cursor

Cursor menyimpan posisi terakhir.

Contoh konsep:

```text
Page 1
  |
  +---- User3
  +---- User2
          |
          v
       cursor
          |
          v
Page 2
  |
  +---- User1
  +---- Staff
```

Dengan keyset pagination, perubahan data di antara request tidak menyebabkan masalah offset yang sama.

---

## Hasil Pagination

Hasil pengujian:

| Halaman | Data         | `has_more` |
| :------ | :----------- | :--------: |
| 1       | User3, User2 |   `true`   |
| 2       | User1, Staff |   `true`   |
| 3       | Admin        |   `false`  |

Tidak terdapat:

```text
duplicate
```

dan tidak terdapat:

```text
missing record
```

---

## EXPLAIN ANALYZE

Query cursor pagination diverifikasi menggunakan:

```sql
EXPLAIN ANALYZE
```

Hasil pengujian menunjukkan PostgreSQL menggunakan:

```text
Index Scan
```

bukan:

```text
Seq Scan
```

Query halaman berikutnya juga menggunakan row value comparison:

```text
(created_at, id) < (...)
```

yang sesuai dengan pasangan kolom cursor.

---

# Declarative Validation

## Sebelum

Validasi sebelumnya dilakukan secara manual:

```text
if NIM kosong
if name kosong
if email tidak valid
if grade di luar range
if password lemah
```

Logika tersebut berada di business rule.

---

## Sesudah

Validasi dipindahkan ke struct menggunakan tag:

```go
type CreateStudentRequest struct {
    NIM      string  `json:"nim" validate:"required,nim"`
    Name     string  `json:"name" validate:"required,min=2,max=100"`
    Email    string  `json:"email" validate:"required,email,max=120"`
    Grade    float64 `json:"grade" validate:"min=0,max=100"`
    Password string  `json:"password" validate:"required,min=8,max=72,nospace"`
}
```

Service kemudian cukup memanggil:

```text
helper.ValidateStruct(req)
```

Pendekatan ini membuat aturan validasi terlihat langsung pada definisi request.

---

## Custom Validator `nim`

Format NIM memiliki aturan domain sendiri sehingga validator bawaan tidak cukup.

Custom validator:

```text
nim
```

memastikan:

* Panjang 3–20 karakter.
* Karakter pertama huruf kapital.
* Karakter berikutnya berupa angka.

Contoh:

```text
S001
A12345
```

valid.

Sedangkan format yang tidak mengikuti aturan tersebut ditolak.

Pesan validasi dipusatkan melalui:

```text
messageFor
```

sehingga client mendapatkan pesan yang konsisten.

---

## PATCH dan `omitnil`

Pointer pada PATCH dikombinasikan dengan:

```text
omitnil
```

untuk membedakan:

```text
field tidak dikirim
```

dan:

```text
field dikirim tetapi kosong
```

Hasil pengujian:

| Request            | Expected | Actual |
| :----------------- | :------: | :----: |
| `{"name":""}`      |   `422`  |  `422` |
| `{"nim":"abc123"}` |   `422`  |  `422` |
| `{"grade":72}`     |   `200`  |  `200` |

---

# Centralized Error Handling

## AppError

Modul ini menggunakan satu bentuk error aplikasi:

```text
AppError
```

Informasi error meliputi:

```text
Status
Code
Message
Fields
cause
```

`cause` tetap private.

Client hanya menerima informasi yang memang merupakan bagian dari API contract.

---

## ErrorHandler

Semua error diproses melalui centralized:

```text
ErrorHandler
```

Alurnya:

```text
Error
  |
  v
ErrorHandler
  |
  +---- AppError
  |
  +---- Fiber Error
  |
  +---- Unknown Error
  |
  v
Standard Error Response
```

Dengan pendekatan ini, seluruh endpoint menggunakan bentuk response error yang konsisten.

---

## Error Response

Format response kegagalan memiliki informasi seperti:

```json
{
  "success": false,
  "code": "VALIDATION_ERROR",
  "message": "validasi gagal",
  "fields": {},
  "request_id": "..."
}
```

`request_id` digunakan untuk menghubungkan response client dengan log server.

---

## Error Code

Error code yang digunakan pada API:

| Code                     | Status | Situasi                                              |
| :----------------------- | :----: | :--------------------------------------------------- |
| `VALIDATION_ERROR`       |  `422` | Body melanggar aturan validasi                       |
| `BAD_REQUEST`            |  `400` | Request tidak dapat diproses / parameter tidak valid |
| `UNAUTHORIZED`           |  `401` | Authentication tidak valid                           |
| `FORBIDDEN`              |  `403` | Tidak memiliki permission / bukan pemilik            |
| `NOT_FOUND`              |  `404` | Resource tidak ditemukan                             |
| `CONFLICT`               |  `409` | NIM atau email sudah digunakan                       |
| `UNSUPPORTED_MEDIA_TYPE` |  `415` | Content-Type tidak didukung                          |
| `NOT_ACCEPTABLE`         |  `406` | Format response tidak dapat dipenuhi                 |
| `TOO_MANY_REQUESTS`      |  `429` | Rate limit terlampaui                                |
| `INTERNAL_ERROR`         |  `500` | Error internal / database                            |
| `SERVICE_UNAVAILABLE`    |  `503` | Service atau database tidak tersedia                 |

---

# Content Negotiation

Satu endpoint dapat memberikan resource dalam beberapa format berdasarkan header:

```text
Accept
```

Tidak diperlukan endpoint terpisah untuk CSV.

---

## JSON

Request:

```text
Accept: application/json
```

Response:

```text
200 OK
Content-Type: application/json
```

JSON menjadi format default.

---

## CSV

Request:

```text
Accept: text/csv
```

Response:

```text
200 OK
Content-Type: text/csv; charset=utf-8
Content-Disposition: attachment; filename="students.csv"
```

Response berisi header CSV dan data student.

---

## Unsupported Format

Request:

```text
Accept: application/xml
```

Response:

```text
406 Not Acceptable
```

Response error tetap menggunakan JSON sehingga client mendapatkan informasi kesalahan dalam format yang konsisten.

---

## Bug CSV dan `Flush()`

Salah satu behavioral bug terjadi ketika CSV response selalu kosong.

Penyebabnya adalah:

```text
csv.Writer
    |
    v
bufio.Writer
    |
    v
strings.Builder
```

Data masih berada di buffer internal.

Memanggil:

```text
Error()
```

tidak otomatis memindahkan data dari buffer.

Perbaikannya:

```go
writer.Flush()

if err := writer.Error(); err != nil {
    return Internal(err)
}
```

Setelah `Flush()`:

```text
CSV body
    |
    v
tersedia di response
```

---

# API Contract

Base URL:

```text
http://localhost:3000/api/v1
```

Endpoint utama yang terpengaruh oleh Modul 7 tetap menggunakan resource dari modul sebelumnya.

### Students

| Method   | Endpoint        | Fungsi                                  |
| :------- | :-------------- | :-------------------------------------- |
| `GET`    | `/students`     | Daftar student dengan cursor pagination |
| `GET`    | `/students/:id` | Detail student                          |
| `POST`   | `/students`     | Membuat student                         |
| `PUT`    | `/students/:id` | Mengganti student                       |
| `PATCH`  | `/students/:id` | Partial update                          |
| `DELETE` | `/students/:id` | Menghapus student                       |

`GET /students` menjadi endpoint utama untuk fitur:

```text
cursor pagination
+
content negotiation
+
centralized error handling
```

---

# Contoh Request

## Cursor Pagination

Konsep request:

```bash
curl -i "http://localhost:3000/api/v1/students?limit=2"
```

Request berikutnya menggunakan cursor yang diberikan response sebelumnya.

---

## JSON

```bash
curl -i http://localhost:3000/api/v1/students -H "Accept: application/json"
```

---

## CSV

```bash
curl -i http://localhost:3000/api/v1/students -H "Accept: text/csv"
```

---

## Unsupported Media

```bash
curl -i http://localhost:3000/api/v1/students -H "Accept: application/xml"
```

Expected:

```text
406 Not Acceptable
```

---

# Aturan yang Tidak Bisa Dijadikan Validation Tag

Tidak semua business rule dapat dipindahkan ke declarative validation.

Contoh:

```text
IsEmptyPatch
```

Aturannya:

```text
setidaknya satu field harus dikirim
```

Aturan tersebut melibatkan hubungan beberapa field sekaligus.

Validator tag lebih cocok untuk:

```text
satu field
```

sedangkan:

```text
relasi antar-field
```

tetap membutuhkan business rule biasa.

Hal serupa berlaku untuk:

```text
ValidateAssignRole
```

karena aturan tersebut membutuhkan:

```text
current.StudentID
```

dan:

```text
targetID
```

Tag validation tidak memiliki akses terhadap konteks pengguna tersebut.

---

# Security Consideration untuk Cursor

Cursor hanya boleh digunakan sebagai:

```text
penanda posisi
```

Cursor tidak boleh menjadi sumber informasi authentication atau authorization.

Jangan menyimpan:

```text
user_id
role
access_token
permission
```

di dalam cursor sebagai dasar keputusan akses.

Cursor yang dikirim client tidak boleh dipercaya sebagai sumber identitas.

Authorization tetap harus berasal dari:

```text
JWT / authenticated context
```

bukan dari isi cursor.

---

# Trade-off Cursor Pagination

Cursor pagination memberikan:

* Tidak mudah mengalami duplicate record ketika data berubah.
* Tidak mudah mengalami missing record akibat insert/delete di antara request.
* Lebih efisien untuk pagination pada dataset besar.
* Dapat memanfaatkan index yang sesuai.

Namun terdapat trade-off:

```text
Tidak dapat dengan mudah melompat ke:
"halaman 3 dari 40"
```

Untuk kebutuhan seperti itu, offset pagination masih dapat berguna.

Pendekatan yang dipilih untuk API ini adalah menjadikan cursor sebagai mekanisme utama untuk endpoint list.

---

# Testing

Menjalankan seluruh test:

```bash
go test ./... -v
```

Static build:

```bash
go build ./...
```

Testing database:

```bash
EXPLAIN ANALYZE
```

digunakan untuk memverifikasi penggunaan index pada cursor pagination.

Pengujian runtime menggunakan:

```text
PowerShell
curl.exe
PostgreSQL / psql
```

---

## Skenario Verifikasi

Beberapa hasil penting yang diverifikasi:

| Skenario                      | Expected Result               |
| :---------------------------- | :---------------------------- |
| Password kuat                 | Diterima                      |
| Password lemah                | `422`                         |
| Validation error              | `422`                         |
| Database error tidak dikenal  | `500`                         |
| Error 4xx                     | Log `WARN`                    |
| Error 5xx                     | Log `ERROR`                   |
| Pagination halaman berikutnya | Tidak duplicate               |
| Pagination                    | Tidak missing record          |
| CSV                           | Body tidak kosong             |
| `Accept: application/json`    | JSON                          |
| `Accept: text/csv`            | CSV                           |
| `Accept: application/xml`     | `406`                         |
| Error response                | Format konsisten              |
| Access log                    | Status sesuai response client |

---

# Struktur Folder

Modul 7 mempertahankan struktur Clean Architecture dari modul sebelumnya.

```text
tugas07/
├── app/
│   ├── model/
│   ├── repository/
│   └── service/
├── helper/
│   ├── errors.go
│   ├── validator.go
│   └── ...
├── middleware/
│   └── middleware.go
├── route/
│   └── route.go
├── config/
│   ├── app.go
│   └── logger.go
├── database/
├── migrations/
├── main.go
├── go.mod
└── README.md
```

Perubahan Modul 7 terutama memperkuat:

```text
helper
   |
   +---- error contract
   +---- validation

config
   |
   +---- centralized ErrorHandler
   +---- logging

service
   |
   +---- business rules
   +---- pagination logic

repository
   |
   +---- cursor-based query
```

---

# Menjalankan Aplikasi

### 1. Masuk ke folder

```bash
cd tugas07
```

### 2. Install dependency

```bash
go mod tidy
```

### 3. Konfigurasi environment

Gunakan `.env.example` sebagai referensi.

Pastikan PostgreSQL tersedia dan konfigurasi dari modul sebelumnya telah disiapkan.

### 4. Jalankan migration

Jalankan migration sesuai urutan yang tersedia pada project.

Pastikan index cursor pagination telah dibuat:

```text
students_created_at_id_desc_idx
```

### 5. Jalankan server

```bash
go run .
```

Server:

```text
http://localhost:3000
```

---

# Kesimpulan

Modul 7 memperkuat API Students bukan dengan menambah banyak endpoint baru, tetapi dengan memperbaiki kualitas internal dan kontrak API.

Hasil utama:

```text
9 bugs fixed
        |
        +---- 3 compiler errors
        |
        +---- 6 behavioral bugs
```

Kemudian API diperkuat dengan:

```text
Declarative Validation
        +
Centralized Error Handling
        +
Cursor Pagination
        +
Content Negotiation
        +
Stable Error Codes
        +
Request ID
```

Dengan demikian, API menjadi lebih konsisten, lebih mudah diuji, dan lebih siap digunakan sebagai fondasi pengembangan modul berikutnya.

---

## Repositori

Kode sumber Modul 7:

`https://github.com/vxpal3n/go-workspace/tree/main/tugas07`

---

## Sumber Bantuan

* Dokumentasi Go.
* Dokumentasi Fiber v2.
* Dokumentasi `go-playground/validator`.
* Dokumentasi PostgreSQL.
* RFC 9110 – HTTP Semantics.
* RFC 4180 – CSV.
* Dokumentasi JWT.
* Dokumentasi `pgx`.

Beberapa bagian implementasi dan dokumentasi dibantu oleh alat bantu AI untuk debugging dan penyusunan struktur, sedangkan seluruh bug diverifikasi melalui compiler output, test suite, dan pengujian runtime sesuai kebutuhan praktikum.
