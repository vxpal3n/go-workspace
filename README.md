# Go Workspace – Pemrograman Backend Lanjut (SIP375)

Repositori ini merupakan kumpulan tugas praktikum mata kuliah **Pemrograman Backend Lanjut (SIP375)**.

Setiap modul dikerjakan secara bertahap dengan menerapkan konsep backend development menggunakan Go, mulai dari pengenalan sintaks, REST API, database, Clean Architecture, hingga authentication dan security.

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

---

## Modul 1 – Persiapan & Sintaks Go

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

## Modul 2 – REST API & HTTP Deep Dive

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

## Modul 3 – Database & Repository Pattern

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

## Modul 4 – Clean Architecture

Folder:

```text
tugas04/
```

Modul keempat melakukan restrukturisasi API Students ke dalam pendekatan **Clean Architecture**.

Struktur aplikasi dipisahkan menjadi beberapa bagian:

```text
app/
├── model/
├── repository/
└── service/

helper/
middleware/
route/
config/
database/
migrations/
```

Konsep utama yang diterapkan:

* Dependency Rule.
* Separation of Concerns.
* Business Rules.
* Repository sebagai gateway database.
* Helper sebagai presenter dan request reader.
* Middleware untuk cross-cutting concerns.
* Structured logging.
* Unit testing untuk business rules.

Business rules dipisahkan dari Fiber dan database sehingga dapat diuji secara independen. Struktur ini menjadi fondasi untuk pengembangan authentication pada Modul 5.

Dokumentasi modul:

[`tugas04/README.md`](./tugas04/README.md)

---

## Modul 5 – Authentication & Security

Folder:

```text
tugas05/
```

Modul kelima mengembangkan hasil Modul 4 dengan menjadikan `Student` sebagai **entitas authentication**.

Fitur authentication:

* Register.
* Login.
* Logout.
* Refresh token.
* Profile `/auth/me`.
* JWT access token.
* bcrypt password hashing.
* Refresh token rotation.
* Refresh token hashing menggunakan SHA-256.

Fitur security:

* Authentication middleware.
* Rate limiter pada login.
* CORS policy.
* Body limit 1 MB.
* Validasi JWT secret.
* Mitigasi timing attack.
* Pencegahan JWT algorithm confusion.
* Pencegahan mass assignment pada `role`.
* Password tidak pernah dikirimkan melalui JSON response.

Endpoint authentication:

| Method | Endpoint                | Keterangan           |
| :----- | :---------------------- | :------------------- |
| `POST` | `/api/v1/auth/register` | Registrasi student   |
| `POST` | `/api/v1/auth/login`    | Login                |
| `POST` | `/api/v1/auth/refresh`  | Refresh access token |
| `POST` | `/api/v1/auth/logout`   | Logout               |
| `GET`  | `/api/v1/auth/me`       | Profil pengguna      |

Endpoint student pada Modul 5 dilindungi authentication:

| Method   | Endpoint               | Keterangan      |
| :------- | :--------------------- | :-------------- |
| `GET`    | `/api/v1/students`     | Daftar student  |
| `GET`    | `/api/v1/students/:id` | Detail student  |
| `PUT`    | `/api/v1/students/:id` | Replace student |
| `PATCH`  | `/api/v1/students/:id` | Partial update  |
| `DELETE` | `/api/v1/students/:id` | Hapus student   |

Pembuatan student tidak lagi dilakukan melalui `POST /students`. Registrasi akun dilakukan melalui:

```text
POST /api/v1/auth/register
```

Dokumentasi lengkap:

[`tugas05/README.md`](./tugas05/README.md)

---

## Perkembangan Arsitektur

Repositori ini dikembangkan secara bertahap dari modul ke modul:

```text
┌─────────────────────────────────────┐
│  Modul 1 — Sintaks Go & Fiber       │
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
│            Pattern                  │
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
```

Setiap modul mempertahankan hasil modul sebelumnya dan mengembangkannya dengan konsep baru.

---

## Teknologi Utama

| Teknologi     | Penggunaan                    |
| :------------ | :---------------------------- |
| Go            | Bahasa pemrograman utama      |
| Fiber v2      | Web framework                 |
| PostgreSQL    | Database                      |
| pgx/v5        | PostgreSQL driver             |
| JWT           | Access token                  |
| bcrypt        | Password hashing              |
| SHA-256       | Refresh token hashing         |
| Fiber Limiter | Rate limiting                 |
| CORS          | Cross-Origin Resource Sharing |
| Helmet        | Security headers              |
| `log/slog`    | Structured logging            |
| lumberjack    | Log rotation                  |
| Go testing    | Unit testing                  |

---

## Prasyarat

Pastikan perangkat pengembangan telah memiliki:

* **Go** ≥ 1.27
* **Git**
* **PostgreSQL**
* **Postman** atau `curl`
* **Visual Studio Code** dengan Go extension

Verifikasi:

```bash
go version
git --version
psql --version
```

---

## Clone Repository

```bash
git clone https://github.com/vxpal3n/go-workspace.git
cd go-workspace
```

---

## Menjalankan Modul

Setiap modul memiliki dependency dan konfigurasi masing-masing.

Contoh menjalankan Modul 5:

```bash
cd tugas05
go mod tidy
go run .
```

Server berjalan pada:

```text
http://localhost:3000
```

Untuk konfigurasi database dan environment variable, lihat README pada masing-masing modul.

---

## Dokumentasi

Dokumentasi teknis tersedia pada README masing-masing modul:

* [`tugas01/README.md`](./tugas01/README.md)
* [`tugas02/README.md`](./tugas02/README.md)
* [`tugas03/README.md`](./tugas03/README.md)
* [`tugas04/README.md`](./tugas04/README.md)
* [`tugas05/README.md`](./tugas05/README.md)

README pada setiap modul berisi penjelasan implementasi, struktur kode, cara menjalankan, serta pengujian yang relevan dengan modul tersebut.

---

## Workflow Git

Pengerjaan modul menggunakan commit bertahap untuk mendokumentasikan perkembangan implementasi.

Konvensi commit yang digunakan:

```text
feat     → fitur baru
fix      → perbaikan bug
refactor → restrukturisasi kode
test     → testing
docs     → dokumentasi
chore    → dependency, konfigurasi, tooling
```

Mulai Modul 5, workflow commit menggunakan pendekatan yang lebih terstruktur sehingga setiap perubahan logis dapat ditelusuri melalui Git history.

---

## Repository

GitHub:

https://github.com/vxpal3n/go-workspace

---

## Sumber Bantuan

Dokumentasi dan referensi utama yang digunakan selama pengerjaan:

* Dokumentasi resmi Go.
* Dokumentasi Fiber v2.
* Dokumentasi PostgreSQL.
* Dokumentasi pgx.
* Dokumentasi JWT.
* Dokumentasi bcrypt.

Beberapa bagian kode dan dokumentasi dibantu oleh alat bantu AI untuk debugging dan penyusunan struktur, sedangkan implementasi disesuaikan dengan kebutuhan praktikum.
