# Tugas 06 – Authorization & Role-Based Access Control

Implementasi Modul 6 untuk mata kuliah **Pemrograman Backend Lanjut (SIP375)**.

Modul ini merupakan kelanjutan langsung dari Modul 5. Jika Modul 5 menyelesaikan masalah **authentication** dengan JWT, Modul 6 menyelesaikan masalah berikutnya: **siapa yang boleh melakukan apa terhadap data tertentu**.

Fokus utama modul ini adalah penerapan **Role-Based Access Control (RBAC)**, permission, ownership check, dan business rule yang bergantung pada identitas pengguna.

Implementasi mempertahankan pendekatan **Clean Architecture** dari modul sebelumnya dengan pembagian tanggung jawab:

* **Middleware** menangani keputusan akses yang cukup ditentukan dari role dan permission.
* **Service** menangani keputusan yang membutuhkan data, seperti ownership.
* **Business rules** dipisahkan agar dapat diuji tanpa menjalankan server.
* **Fail closed** diterapkan sehingga role atau permission yang tidak dikenal tidak pernah otomatis mendapatkan akses.

---

## Tujuan

Modul 6 bertujuan menambahkan sistem **authorization** pada API Students yang sebelumnya telah memiliki authentication.

Implementasi mencakup:

* RBAC menggunakan tabel `roles`, `permissions`, dan `role_permissions`.
* Tiga role: `admin`, `staff`, dan `user`.
* Permission untuk operasi pada resource `students`.
* Middleware `RequirePermission`.
* Ownership check menggunakan `owner_id`.
* Pemisahan keputusan authorization antara middleware dan service.
* Pencegahan `owner_id` spoofing.
* Business rule untuk mencegah penghapusan akun sendiri.
* Business rule untuk mencegah perubahan role diri sendiri.
* Prinsip **fail closed**.
* Penggunaan status `401`, `403`, dan `422` sesuai konteks.
* Unit test untuk authorization rules.

Modul ini tidak menggantikan authentication dari Modul 5. Authorization dibangun di atas identitas pengguna yang sudah diperoleh melalui JWT.

---

## Konsep Authentication vs Authorization

Modul 5 menjawab:

```text
"Siapa pengguna ini?"
```

Modul 6 menjawab:

```text
"Pengguna ini boleh melakukan apa?"
```

Alur sederhananya:

```text
Request
   |
   v
RequireAuth
   |
   |  JWT valid?
   v
Current User
   |
   v
Permission / Ownership Check
   |
   +---- allowed ----> Service
   |
   +---- denied -----> 403
```

Authentication tetap menggunakan mekanisme JWT dari Modul 5, sedangkan authorization menggunakan role, permission, dan ownership.

---

## Teknologi

| Komponen         | Teknologi               |
| :--------------- | :---------------------- |
| Bahasa           | Go 1.27.0               |
| Framework        | Fiber v2                |
| Database         | PostgreSQL 15+          |
| Database Driver  | pgx/v5 + pgxpool        |
| Authentication   | JWT                     |
| Password Hashing | bcrypt                  |
| Authorization    | RBAC + Ownership        |
| Testing          | Go testing + `curl`     |
| Logging          | `log/slog` + lumberjack |

---

## Role

Modul ini menggunakan tiga role:

| Role    | Deskripsi                                                                                  |
| :------ | :----------------------------------------------------------------------------------------- |
| `admin` | Memiliki akses penuh terhadap data student dan pengaturan role                             |
| `staff` | Dapat melihat dan membuat data student, tetapi tidak dapat mengubah atau menghapus data    |
| `user`  | Tidak memiliki permission global; akses terhadap data sendiri ditentukan melalui ownership |

Role `user` sengaja tidak diberikan permission seperti `student:list` atau `student:read:any`.

Akses terhadap data milik sendiri berasal dari **ownership**, bukan dari permission global.

---

## Permission

Permission yang digunakan:

| Permission           | Fungsi                         |
| :------------------- | :----------------------------- |
| `student:list`       | Melihat daftar seluruh student |
| `student:read:any`   | Melihat data student mana pun  |
| `student:create`     | Membuat student baru           |
| `student:update:any` | Mengubah data student mana pun |
| `student:delete`     | Menghapus student              |
| `role:assign`        | Mengubah role student lain     |

Matriks dasar permission:

| Permission           | `admin` | `staff` | `user` |
| :------------------- | :-----: | :-----: | :----: |
| `student:list`       |   Yes   |   Yes   |   No   |
| `student:read:any`   |   Yes   |   Yes   |   No   |
| `student:create`     |   Yes   |   Yes   |   No   |
| `student:update:any` |   Yes   |    No   |   No   |
| `student:delete`     |   Yes   |    No   |   No   |
| `role:assign`        |   Yes   |    No   |   No   |

`user` tetap dapat mengakses data miliknya sendiri melalui ownership check.

---

## Database

RBAC menggunakan tiga tabel utama:

```text
roles
  |
  +---- role_permissions ---- permissions
```

### `roles`

Menyimpan daftar role yang tersedia.

```sql
CREATE TABLE IF NOT EXISTS roles (
    name VARCHAR(20) PRIMARY KEY,
    description VARCHAR(150) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Role yang digunakan:

```text
admin
staff
user
```

### `permissions`

Menyimpan daftar permission yang tersedia.

```sql
CREATE TABLE IF NOT EXISTS permissions (
    name VARCHAR(50) PRIMARY KEY,
    description VARCHAR(150) NOT NULL
);
```

### `role_permissions`

Menghubungkan role dengan permission.

```sql
CREATE TABLE IF NOT EXISTS role_permissions (
    role_name VARCHAR(20) NOT NULL REFERENCES roles(name) ON DELETE CASCADE,
    permission_name VARCHAR(50) NOT NULL REFERENCES permissions(name) ON DELETE CASCADE,
    PRIMARY KEY (role_name, permission_name)
);
```

---

## Ownership

Migration authorization menambahkan:

```text
students.owner_id
```

Kolom ini merupakan foreign key yang mengarah kembali ke:

```text
students(id)
```

sehingga ownership dapat direpresentasikan secara langsung pada record student.

Struktur sederhananya:

```text
students
├── id
├── nim
├── name
├── email
├── role
├── owner_id
└── ...
```

`owner_id` bersifat nullable.

Alasannya:

* Student yang melakukan self-register tidak memiliki owner eksternal.
* Student yang dibuat melalui `POST /students` memiliki `owner_id` yang berasal dari pengguna yang sedang login.

---

## Ownership vs Permission

Authorization terhadap student tidak hanya bergantung pada permission.

Untuk endpoint yang membutuhkan akses terhadap satu student, terdapat tiga jalur yang dapat memberikan akses:

```text
1. Self-access
   current.StudentID == targetID

2. Ownership
   current.StudentID == ownerID

3. Permission-based
   role memiliki permission :any
```

Contohnya:

```text
User A
  |
  +---- mengakses student miliknya sendiri
  |          -> ALLOWED
  |
  +---- mengakses student yang bukan miliknya
  |          -> DENIED
  |
  +---- admin mengakses student lain
             -> ALLOWED
```

Ketiga jalur gagal berarti akses ditolak.

---

## Fail Closed

Authorization menggunakan prinsip **fail closed**.

Artinya, jika sistem tidak dapat membuktikan bahwa suatu akses diperbolehkan, hasilnya adalah:

```text
DENY
```

Contohnya:

* `PermissionSet` bernilai `nil` → ditolak.
* Role tidak dikenal → ditolak.
* Permission tidak dikenal → ditolak.
* User bukan pemilik → ditolak.
* User tidak memiliki permission `:any` → ditolak.

Tidak ada logika seperti:

```text
"permission tidak ditemukan, jadi mungkin boleh."
```

Semua kondisi tidak dikenal dianggap tidak memiliki akses.

---

## Implementasi PermissionSet

Permission dimuat ke dalam memory dalam bentuk mapping role → permission.

Konsepnya:

```text
admin
 ├── student:list
 ├── student:read:any
 ├── student:create
 ├── student:update:any
 ├── student:delete
 └── role:assign

staff
 ├── student:list
 ├── student:read:any
 └── student:create

user
 └── tidak memiliki permission global
```

Lookup permission dilakukan melalui `PermissionSet.Can()`.

---

## Middleware Authorization

Middleware:

```text
RequirePermission
```

digunakan untuk keputusan yang tidak membutuhkan isi record database.

Contoh:

```text
GET    /students
POST   /students
DELETE /students/:id
PATCH  /students/:id/role
```

Alurnya:

```text
JWT
 |
 v
CurrentUser
 |
 v
PermissionSet.Can(role, permission)
 |
 +---- true  ---> next handler
 |
 +---- false ---> 403 Forbidden
```

Middleware tidak bertanggung jawab terhadap ownership.

---

## Ownership Check di Service

Ownership ditempatkan di service karena service perlu mengetahui data target.

Fungsi utama:

```text
CanAccessStudent
```

memeriksa:

```text
self-access
    OR
ownership
    OR
permission :any
```

Pendekatan ini menghindari middleware melakukan query database hanya untuk mengetahui siapa pemilik suatu record.

---

## Business Rules

Beberapa aturan tidak cukup diselesaikan dengan permission.

### Tidak boleh menghapus akun sendiri

Walaupun seorang `admin` memiliki:

```text
student:delete
```

admin tetap tidak boleh menghapus dirinya sendiri.

Hasil:

```text
403 Forbidden
```

Aturan ini membutuhkan perbandingan:

```text
current.StudentID
```

dengan:

```text
targetID
```

sehingga keputusan berada di service.

---

### Tidak boleh mengubah role diri sendiri

Admin memiliki:

```text
role:assign
```

tetapi tidak boleh mengubah role dirinya sendiri.

Contoh:

```text
Admin A
  |
  +---- mengubah role User B
  |         -> allowed
  |
  +---- mengubah role dirinya sendiri
            -> 422
```

Perbedaannya penting:

```text
403 = tidak memiliki hak

422 = memiliki hak, tetapi request melanggar business rule
```

---

## Pencegahan `owner_id` Spoofing

Client tidak diperbolehkan menentukan owner melalui request body.

Contoh request berbahaya:

```json
{
  "nim": "X001",
  "name": "Spoof",
  "email": "spoof@test.com",
  "grade": 80,
  "password": "Rahasia123",
  "owner_id": 999
}
```

Field tersebut tidak digunakan untuk menentukan ownership.

Server mengambil owner dari authenticated user:

```text
JWT
  |
  v
CurrentUser.StudentID
  |
  v
owner_id
```

Dengan demikian client tidak dapat mengaku sebagai pemilik record lain.

---

## API Contract

Base URL:

```text
http://localhost:3000/api/v1
```

### Student Endpoints

| Method   | Endpoint             | Permission / Rule                | Fungsi                   |
| :------- | :------------------- | :------------------------------- | :----------------------- |
| `GET`    | `/students`          | `student:list`                   | Melihat seluruh student  |
| `GET`    | `/students/:id`      | Ownership / `student:read:any`   | Melihat student tertentu |
| `POST`   | `/students`          | `student:create`                 | Membuat student          |
| `PUT`    | `/students/:id`      | Ownership / `student:update:any` | Mengubah student         |
| `PATCH`  | `/students/:id`      | Ownership / `student:update:any` | Mengubah sebagian data   |
| `DELETE` | `/students/:id`      | `student:delete` + business rule | Menghapus student        |
| `PATCH`  | `/students/:id/role` | `role:assign` + business rule    | Mengubah role            |

Semua endpoint student tetap melewati authentication dari Modul 5.

---

## Matriks Akses

Hasil pengujian utama:

| Request                                   | `admin` | `staff` | `user` |
| :---------------------------------------- | :-----: | :-----: | :----: |
| `GET /students`                           |  `200`  |  `200`  |  `403` |
| `GET /students/:id` — self                |  `200`  |  `200`  |  `200` |
| `GET /students/:id` — milik orang lain    |  `200`  |  `200`  |  `403` |
| `POST /students`                          |  `201`  |  `201`  |  `403` |
| `PUT /students/:id` — milik orang lain    |  `200`  |  `403`  |  `403` |
| `DELETE /students/:id` — orang lain       |  `204`  |  `403`  |  `403` |
| `DELETE /students/:id` — diri sendiri     |  `403`  |  `403`  |  `403` |
| `PATCH /students/:id/role` — orang lain   |  `200`  |  `403`  |  `403` |
| `PATCH /students/:id/role` — diri sendiri |  `422`  |  `403`  |  `403` |
| Tanpa Authorization                       |  `401`  |  `401`  |  `401` |

---

## Status HTTP

Authorization dibedakan dari authentication dan business validation:

| Status                     | Makna                                                     |
| :------------------------- | :-------------------------------------------------------- |
| `401 Unauthorized`         | Request tidak memiliki authentication yang valid          |
| `403 Forbidden`            | User sudah terautentikasi tetapi tidak memiliki hak akses |
| `422 Unprocessable Entity` | Request melanggar business rule                           |

Contoh:

```text
Tidak ada token
    -> 401

User mencoba membaca data milik orang lain
    -> 403

Admin mencoba mengubah role dirinya sendiri
    -> 422
```

---

## Pengujian Negatif

Modul ini tidak hanya menguji happy path.

Beberapa skenario keamanan yang diuji:

### Owner ID Spoofing

Client mencoba mengirim:

```text
owner_id = 999
```

Hasil:

```text
owner_id tetap berasal dari authenticated user.
```

---

### User Mengubah Data Orang Lain

User mencoba:

```text
PUT /students/1
```

terhadap student yang bukan miliknya.

Hasil:

```text
403 Forbidden
```

---

### Staff Create

Staff membuat student baru.

Database kemudian diverifikasi untuk memastikan:

```text
owner_id = id staff
```

bukan nilai yang dikirim oleh client.

---

### Role Berubah tetapi Token Lama

JWT bersifat stateless.

Jika:

```text
user
```

diubah menjadi:

```text
staff
```

token lama masih membawa:

```json
{
  "role": "user"
}
```

hingga token tersebut kedaluwarsa.

Setelah login ulang, token baru membawa role terbaru.

Ini merupakan karakteristik desain token stateless, bukan bug implementasi.

---

## Unit Testing

Business rules authorization diuji secara terisolasi.

Test utama:

```text
TestCanAccessStudent
TestValidateAssignRole
```

`TestCanAccessStudent` memverifikasi:

* Self-access.
* Ownership.
* Permission-based access.
* Access ditolak jika seluruh kondisi gagal.

`TestValidateAssignRole` memverifikasi:

* Admin dapat mengubah role student lain.
* Admin tidak dapat mengubah role dirinya sendiri.
* Role yang tidak dikenal ditolak.

Menjalankan test:

```bash
cd tugas06
go test ./app/service -v
```

---

## Struktur Folder

Struktur aplikasi tetap mengikuti Clean Architecture dari modul sebelumnya.

```text
tugas06/
├── app/
│   ├── model/
│   ├── repository/
│   └── service/
│       ├── student_service.go
│       ├── student_rules.go
│       └── student_authz_rules.go
├── helper/
│   └── authz.go
├── middleware/
│   └── authz.go
├── route/
│   └── route.go
├── config/
├── database/
├── migrations/
├── main.go
├── go.mod
└── README.md
```

Nama file dapat berkembang mengikuti implementasi, tetapi tanggung jawab setiap layer tetap dipisahkan.

---

## Menjalankan Aplikasi

### 1. Masuk ke folder project

```bash
cd tugas06
```

### 2. Install dependency

```bash
go mod tidy
```

### 3. Konfigurasi environment

Gunakan `.env.example` sebagai referensi konfigurasi.

Pastikan PostgreSQL dan konfigurasi JWT dari Modul 5 tersedia.

### 4. Jalankan migration

Jalankan migration Modul 6 sesuai urutan migration pada project.

Migration mencakup:

```text
roles
permissions
role_permissions
owner_id
foreign key
index
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

## Analisis Desain

### Mengapa ownership tidak ditaruh di middleware?

Middleware tidak mengetahui isi row:

```text
/students/7
```

hanya memberi informasi bahwa target ID adalah:

```text
7
```

Middleware tidak mengetahui:

```text
owner_id dari student 7
```

Jika middleware melakukan query hanya untuk authorization, service kemungkinan akan melakukan query yang sama lagi.

Karena itu pembagian yang digunakan adalah:

```text
Middleware
    |
    +-- keputusan berbasis role/permission

Service
    |
    +-- keputusan berbasis data/ownership
```

---

### Risiko Dua Tempat Pemeriksaan

Pembagian ini memiliki konsekuensi: endpoint baru dapat lupa memasang permission middleware.

Mitigasi yang digunakan:

1. Semua endpoint student tetap berada di bawah `RequireAuth`.
2. Matriks authorization diperlakukan sebagai kontrak.
3. Setiap endpoint diuji terhadap role yang relevan.

Dengan demikian route juga berfungsi sebagai peta hak akses aplikasi.

---

### RBAC untuk Relasi Dosen-Mahasiswa

RBAC murni tidak cukup untuk aturan seperti:

```text
"Dosen hanya boleh melihat mahasiswa bimbingannya."
```

Pertanyaan tersebut tidak hanya bergantung pada role:

```text
role == dosen
```

tetapi juga hubungan:

```text
dosen X membimbing mahasiswa Y?
```

Untuk kebutuhan tersebut, pendekatan yang dipertimbangkan adalah tabel relasi:

```text
dosen_mahasiswa
├── dosen_id
└── mahasiswa_id
```

atau pendekatan ABAC.

Untuk konteks proyek ini, relasi tabel dianggap lebih sederhana dan lebih mudah diaudit.

---

## Keterbatasan

Karena role berada di JWT yang bersifat stateless, perubahan role tidak langsung mengubah access token yang sudah diterbitkan.

Contoh:

```text
Token lama:
role = user

Database:
role = staff
```

Token lama tetap membawa role `user` sampai expiration.

Untuk sistem yang membutuhkan pencabutan hak secara langsung, alternatifnya adalah:

* Mengecek role ke database pada setiap request.
* Menggunakan token revocation / blocklist.
* Menggunakan access token dengan lifetime yang lebih pendek.

Untuk modul ini, perilaku tersebut didokumentasikan sebagai trade-off desain.

---

## Repositori

Kode sumber Modul 6:

`https://github.com/vxpal3n/go-workspace/tree/main/tugas06`

---

## Sumber Bantuan

* Dokumentasi Go.
* Dokumentasi Fiber v2.
* Dokumentasi PostgreSQL.
* Dokumentasi JWT.
* Dokumentasi bcrypt.
* Referensi HTTP status code.
* Materi praktikum Pemrograman Backend Lanjut (SIP375).

Beberapa bagian implementasi dan dokumentasi dibantu oleh alat bantu AI untuk debugging dan penyusunan struktur, sedangkan logika utama disesuaikan dengan kebutuhan praktikum.
