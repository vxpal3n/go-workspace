# Tugas 01 — Persiapan & Sintaks Go

## Deskripsi

Tugas ini merupakan tahap awal pembelajaran backend menggunakan bahasa **Go**.

Materi yang dipraktikkan mencakup sintaks dasar Go, penggunaan variabel, slice, map, pointer, struct, method, serta pengenalan framework **Fiber v2** melalui aplikasi sederhana. 

## Struktur

```text
tugas01/
├── README.md
├── fiber/
│   └── main.go
└── syntax/
    ├── variabel/
    │   └── main.go
    ├── pointer/
    │   └── main.go
    └── struct/
        └── main.go
```

## Materi

### 1. Fiber

Berisi implementasi sederhana **Hello World** menggunakan Fiber v2.

### 2. Variabel

Mendemonstrasikan penggunaan:

* Variabel
* Slice
* Map
* Operasi dasar terhadap data

### 3. Pointer

Mendemonstrasikan konsep:

* Pointer
* Pass by value
* Pass by reference
* Fungsi `swap`
* Fungsi `updateSlice`

### 4. Struct & Method

Mendemonstrasikan penggunaan struct `Student` beserta beberapa method:

* `GetInfo`
* `UpdateGrade`
* `Activate`
* `Deactivate`

## Cara Menjalankan

Setiap contoh dapat dijalankan dari direktori masing-masing.

### Fiber

```bash
cd fiber
go run main.go
```

Server berjalan pada:

```text
http://localhost:8080
```

### Variabel

```bash
cd syntax/variabel
go run main.go
```

### Pointer

```bash
cd syntax/pointer
go run main.go
```

### Struct

```bash
cd syntax/struct
go run main.go
```

## Catatan

Modul ini menjadi fondasi sebelum memasuki pengembangan REST API pada **Modul 2**.

Materi dasar seperti struct, pointer, slice, dan method digunakan kembali ketika membangun aplikasi backend yang lebih kompleks.

