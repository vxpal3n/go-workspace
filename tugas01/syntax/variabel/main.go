package main

import "fmt"

func cetakMap(m map[string]int) {
	if len(m) == 0 {
		fmt.Println("  (kosong)")
		return
	}
	for nama, nilai := range m {
		fmt.Printf("  %s : %d\n", nama, nilai)
	}
}

func main() {
	fmt.Println("||| PROGRAM MANAJEMEN DATA MAHASISWA |||")
	fmt.Println("")

	var nama string = "Valentino Chandra" // string
	umur := 19                            // int (short declaration)
	ipk := 3.85                           // float64
	lulus := false                        // bool
	hobi := []string{"Coding", "Gaming"}  // slice

	fmt.Println("~~~ Data Pribadi ~~~")
	fmt.Printf("Nama  : %s\n", nama)
	fmt.Printf("Umur  : %d\n", umur)
	fmt.Printf("IPK   : %.2f\n", ipk)
	fmt.Printf("Lulus : %v\n", lulus)
	fmt.Printf("Hobi  : %v\n\n", hobi)


	dataMahasiswa := map[string]int{
		"Alice":   85,
		"Bob":     90,
		"Charlie": 75,
	}

	fmt.Println("~~~ Keadaan Awal Map ~~~")
	cetakMap(dataMahasiswa)
	fmt.Println()

	fmt.Println("~~~ Operasi TAMBAH ~~~")
	fmt.Println(">> Menambahkan: David = 88")
	dataMahasiswa["David"] = 88
	fmt.Println(">> Setelah penambahan:")
	cetakMap(dataMahasiswa)
	fmt.Println()


	fmt.Println("~~~ Operasi BACA ~~~")
	namaCari := "Alice"
	if nilai, ada := dataMahasiswa[namaCari]; ada {
		fmt.Printf(">> %s ditemukan, nilainya %d\n", namaCari, nilai)
	} else {
		fmt.Printf(">> %s tidak ditemukan\n", namaCari)
	}

	namaCari = "Eve"
	if nilai, ada := dataMahasiswa[namaCari]; ada {
		fmt.Printf(">> %s ditemukan, nilainya %d\n", namaCari, nilai)
	} else {
		fmt.Printf(">> %s tidak ditemukan\n", namaCari)
	}
	fmt.Println()


	fmt.Println("~~~ Operasi HAPUS ~~~")
	fmt.Println(">> Menghapus: Charlie")
	delete(dataMahasiswa, "Charlie")
	fmt.Println(">> Setelah penghapusan:")
	cetakMap(dataMahasiswa)
	fmt.Println()


	fmt.Println("~~~ Operasi TELUSURI ~~~")
	fmt.Println(">> Seluruh data mahasiswa saat ini:")
	cetakMap(dataMahasiswa)
	fmt.Println()


	fmt.Println("~~~ Operasi UPDATE ~~~")
	fmt.Println(">> Mengubah nilai Bob menjadi 95")
	dataMahasiswa["Bob"] = 95
	fmt.Println(">> Setelah update:")
	cetakMap(dataMahasiswa)
}