# Tugas 02 — REST API & HTTP Deep Dive

## Deskripsi

Tugas ini merupakan implementasi **REST API menggunakan Go dan Fiber v2** dengan fokus pada pemahaman HTTP dan prinsip RESTful API.

Implementasi mencakup operasi CRUD, penggunaan HTTP method, HTTP status code, validasi request, pagination, filtering, searching, sorting, serta perbedaan antara `PUT` dan `PATCH`. 

## Tujuan Pembelajaran

Implementasi ini mencakup:

* Penggunaan HTTP method `GET`, `POST`, `PUT`, `PATCH`, dan `DELETE`
* Penggunaan HTTP status code yang sesuai
* Implementasi CRUD untuk entitas `Student`
* Pagination
* Searching
* Sorting
* Filtering
* Validasi data
* Perbedaan `PUT` dan `PATCH`
* Dokumentasi kontrak API

## Teknologi

* **Go**
* **Fiber v2**
* **REST API**
* **HTTP**
* **cURL / Postman**

## Struktur Project

```text
tugas02/
├── README.md
├── go.mod
├── go.sum
├── main.go
├── model.go
├── helper.go
└── handler.go
```

| File         | Tanggung Jawab                                 |
| ------------ | ---------------------------------------------- |
| `main.go`    | Inisialisasi aplikasi, middleware, dan routing |
| `model.go`   | Model `Student`, request DTO, dan response     |
| `helper.go`  | Helper response, parsing query, dan ID         |
| `handler.go` | CRUD, filtering, sorting, dan pagination       |
| `go.mod`     | Dependency management                          |
| `go.sum`     | Dependency checksums                           |

## Menjalankan Project

Masuk ke direktori project:

```bash
cd tugas02
```

Install dependency:

```bash
go mod tidy
```

Jalankan server:

```bash
go run .
```

Server akan tersedia pada:

```text
http://localhost:3000
```

### Base URL

```text
http://localhost:3000/api/v1
```

---

# API Documentation

## 1. Get All Students

```http
GET /students
```

Mengambil daftar student dengan dukungan pagination, searching, sorting, dan filtering.

### Query Parameters

| Parameter   | Tipe     | Deskripsi                           | Default |
| ----------- | -------- | ----------------------------------- | ------- |
| `page`      | `int`    | Nomor halaman                       | `1`     |
| `limit`     | `int`    | Jumlah data per halaman             | `10`    |
| `search`    | `string` | Pencarian berdasarkan nama atau NIM | —       |
| `sort`      | `string` | Field pengurutan                    | `id`    |
| `order`     | `string` | `asc` atau `desc`                   | `asc`   |
| `is_active` | `bool`   | Filter status aktif                 | —       |
| `min_grade` | `float`  | Nilai minimum                       | —       |
| `max_grade` | `float`  | Nilai maksimum                      | —       |

Field yang dapat digunakan untuk sorting:

```text
id
nim
name
grade
created_at
```

### Contoh

```bash
curl "http://localhost:3000/api/v1/students?page=1&limit=10&sort=name&order=asc"
```

### Response

```json
{
  "success": true,
  "message": "daftar student berhasil diambil",
  "data": [
    {
      "id": 1,
      "nim": "S001",
      "name": "StudentA",
      "grade": 85.5,
      "is_active": true,
      "created_at": "2026-09-08T10:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 3,
    "total_pages": 1
  }
}
```

---

## 2. Get Student by ID

```http
GET /students/:id
```

Mengambil detail student berdasarkan ID.

### Contoh

```bash
curl "http://localhost:3000/api/v1/students/1"
```

### Response

```json
{
  "success": true,
  "message": "student ditemukan",
  "data": {
    "id": 1,
    "nim": "S001",
    "name": "StudentA",
    "grade": 85.5,
    "is_active": true,
    "created_at": "2026-09-08T10:00:00Z"
  }
}
```

Jika student tidak ditemukan:

```json
{
  "success": false,
  "message": "student tidak ditemukan"
}
```

**Status:** `404 Not Found`

---

## 3. Create Student

```http
POST /students
```

Membuat student baru.

### Request Body

```json
{
  "nim": "S004",
  "name": "StudentD",
  "grade": 88.0
}
```

### Validasi

* `NIM` harus unik
* `name` wajib diisi
* `grade` harus berada pada rentang `0–100`

### Response

**Status:** `201 Created`

```json
{
  "success": true,
  "message": "student berhasil dibuat",
  "data": {
    "id": 4,
    "nim": "S004",
    "name": "StudentD",
    "grade": 88,
    "is_active": true,
    "created_at": "..."
  }
}
```

### Validation Error

**Status:** `422 Unprocessable Entity`

```json
{
  "success": false,
  "message": "validasi gagal",
  "errors": {
    "nim": "sudah digunakan",
    "grade": "harus antara 0 dan 100"
  }
}
```

---

## 4. Replace Student — PUT

```http
PUT /students/:id
```

`PUT` digunakan untuk **mengganti seluruh data student**.

### Request Body

Seluruh field harus diberikan:

```json
{
  "nim": "S004",
  "name": "StudentD Baru",
  "grade": 95.0,
  "is_active": false
}
```

### Response

**Status:** `200 OK`

Data student terbaru dikembalikan.

> **PUT vs PATCH:**
> PUT mengganti seluruh representasi resource. Field yang tidak dikirim dianggap tidak lengkap dan akan gagal validasi.

---

## 5. Partial Update — PATCH

```http
PATCH /students/:id
```

`PATCH` digunakan untuk **mengubah sebagian data student**.

### Request Body

Hanya field yang ingin diubah yang perlu dikirim:

```json
{
  "grade": 100.0,
  "is_active": true
}
```

Field lainnya akan mempertahankan nilai sebelumnya.

### Implementasi Pointer

Request `PATCH` menggunakan pointer:

```text
*string
*float64
*bool
```

Pendekatan ini memungkinkan aplikasi membedakan antara:

* Field tidak dikirim → `nil`
* Field dikirim dengan nilai tertentu

---

## 6. Delete Student

```http
DELETE /students/:id
```

Menghapus student berdasarkan ID.

### Contoh

```bash
curl -X DELETE http://localhost:3000/api/v1/students/2 -i
```

### Response

**Status:** `204 No Content`

Tidak terdapat response body.

---

# Searching, Filtering & Sorting

### Searching

Pencarian berdasarkan nama atau NIM:

```bash
curl "http://localhost:3000/api/v1/students?search=li"
```

Pencarian bersifat **case-insensitive**.

### Filtering Grade

Menampilkan student dengan grade antara 80 dan 90:

```bash
curl "http://localhost:3000/api/v1/students?min_grade=80&max_grade=90"
```

### Filtering Status

```bash
curl "http://localhost:3000/api/v1/students?is_active=true"
```

### Sorting

```bash
curl "http://localhost:3000/api/v1/students?sort=name&order=asc"
```

### Pagination

```bash
curl "http://localhost:3000/api/v1/students?page=1&limit=2"
```

Nilai maksimum `limit` adalah `100`.

---

# Testing

## Membuat Data Awal

```bash
curl -X POST http://localhost:3000/api/v1/students -H "Content-Type: application/json" -d "{\"nim\":\"S001\",\"name\":\"StudentA\",\"grade\":85.5}"
```

```bash
curl -X POST http://localhost:3000/api/v1/students -H "Content-Type: application/json" -d "{\"nim\":\"S002\",\"name\":\"Bob\",\"grade\":92.0}"
```

```bash
curl -X POST http://localhost:3000/api/v1/students -H "Content-Type: application/json" -d "{\"nim\":\"S003\",\"name\":\"Charlie\",\"grade\":78.0}"
```

## Duplicate NIM

```bash
curl -X POST http://localhost:3000/api/v1/students -H "Content-Type: application/json" -d "{\"nim\":\"S001\",\"name\":\"Dup\",\"grade\":80}"
```

Digunakan untuk menguji konflik NIM.

**Status:** `409 Conflict`

## Invalid Content-Type

```bash
curl -X POST http://localhost:3000/api/v1/students -d "{\"nim\":\"x\"}"
```

Digunakan untuk menguji request dengan `Content-Type` yang tidak sesuai.

**Status:** `415 Unsupported Media Type`

---

# Security & Best Practices

Implementasi menerapkan beberapa praktik dasar:

### Pagination Limit

`limit` dibatasi maksimal `100` untuk mencegah request mengambil data dalam jumlah yang terlalu besar.

### Sorting Whitelist

Field sorting dibatasi pada:

```text
id
nim
name
grade
created_at
```

Hal ini mencegah penggunaan field yang tidak diizinkan.

### JSON Content-Type Validation

Middleware `requireJSON` memastikan request dengan body menggunakan:

```http
Content-Type: application/json
```

Request yang tidak memenuhi persyaratan akan mendapatkan:

```text
415 Unsupported Media Type
```

### Global Error Handler

Error dikonversi ke format response yang konsisten tanpa mengekspos stack trace kepada client.

### Unique NIM

NIM harus unik pada operasi:

* `POST`
* `PUT`
* `PATCH`

Pada proses update, ID student yang sedang diperbarui tidak dianggap sebagai duplikasi dirinya sendiri.

---

# Limitasi

Data student saat ini disimpan **di memory menggunakan slice**.

Artinya, data akan hilang ketika server dihentikan atau di-restart.

```text
Application
     │
     ▼
 In-Memory Slice
     │
     └── Data hilang saat restart
```

Persistensi database direncanakan untuk modul berikutnya menggunakan **PostgreSQL**.

---

# Learning Outcomes

Melalui implementasi ini, beberapa konsep utama yang dipraktikkan meliputi:

* RESTful API
* HTTP Methods
* HTTP Status Codes
* CRUD
* Request Validation
* Pagination
* Searching
* Filtering
* Sorting
* PUT vs PATCH
* Middleware
* Error Handling
* Data Transfer Object
* In-memory data management

---

## Repository

[GitHub — go-workspace](https://github.com/vxpal3n/go-workspace/tree/main/tugas02)
