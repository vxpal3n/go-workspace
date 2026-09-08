# Go Workspace — Pemrograman Backend Lanjut

Repositori ini berisi kumpulan tugas praktikum mata kuliah **Pemrograman Backend Lanjut (SIP375)**.

Setiap modul disusun secara bertahap untuk membangun pemahaman mengenai bahasa pemrograman Go, framework Fiber, HTTP, REST API, hingga pengembangan backend yang lebih terstruktur.

## Daftar Modul

| Modul                 | Topik                     | Status    |
| --------------------- | ------------------------- | --------- |
| [Tugas 01](./tugas01) | Persiapan & Sintaks Go    | Selesai |
| [Tugas 02](./tugas02) | REST API & HTTP Deep Dive | Selesai |

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
└── tugas02/
    ├── README.md
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── model.go
    ├── helper.go
    └── handler.go
```

## Teknologi

* **Go**
* **Fiber v2**
* **REST API**
* **HTTP**
* **Git & GitHub**
* **Postman / cURL**

## Getting Started

### Prasyarat

Pastikan beberapa tools berikut sudah terpasang:

```bash
go version
git --version
```

Tools yang digunakan:

* Go
* Git
* Postman atau cURL
* Visual Studio Code + Go Extension *(opsional)*

### Clone Repository

```bash
git clone https://github.com/vxpal3n/go-workspace.git
cd go-workspace
```

## Menjalankan Modul

### Modul 1

Masuk ke salah satu direktori program:

```bash
cd tugas01/fiber
go run main.go
```

Untuk contoh sintaks lainnya:

```bash
cd tugas01/syntax/variabel
go run main.go
```

```bash
cd tugas01/syntax/pointer
go run main.go
```

```bash
cd tugas01/syntax/struct
go run main.go
```

### Modul 2

```bash
cd tugas02
go mod tidy
go run .
```

REST API tersedia pada:

```text
http://localhost:3000
```

## REST API — Modul 2

Modul 2 menyediakan operasi CRUD untuk entitas `Student`.

| Method   | Endpoint               | Deskripsi                         |
| -------- | ---------------------- | --------------------------------- |
| `GET`    | `/api/v1/students`     | Mendapatkan daftar student        |
| `GET`    | `/api/v1/students/:id` | Mendapatkan detail student        |
| `POST`   | `/api/v1/students`     | Membuat student baru              |
| `PUT`    | `/api/v1/students/:id` | Mengganti seluruh data student    |
| `PATCH`  | `/api/v1/students/:id` | Memperbarui sebagian data student |
| `DELETE` | `/api/v1/students/:id` | Menghapus student                 |

Dokumentasi lengkap tersedia di:

**[Dokumentasi Tugas 2](./tugas02/README.md)**

## Contoh Request

Membuat student:

```bash
curl -X POST http://localhost:3000/api/v1/students -H "Content-Type: application/json" -d '{"nim":"S001","name":"StudentA","grade":85.5}'
```

Mengambil daftar student:

```bash
curl "http://localhost:3000/api/v1/students?page=1&limit=5"
```

Menghapus student:

```bash
curl -X DELETE http://localhost:3000/api/v1/students/1 -i
```

## Catatan

Pada Modul 2, data student masih disimpan **di memory menggunakan slice**. Oleh karena itu, seluruh data akan hilang ketika server dihentikan atau di-restart.

Persistensi database direncanakan untuk modul berikutnya.

## Repository

[GitHub — go-workspace](https://github.com/vxpal3n/go-workspace)

## Referensi

* [Go Documentation](https://go.dev/)
* [A Tour of Go](https://go.dev/tour/)
* [Fiber Documentation](https://docs.gofiber.io/)

---

**Pemrograman Backend Lanjut — SIP375**
