package main

import "fmt"

func swap(a, b *int) {
	*a, *b = *b, *a
}

func updateSlice(s *[]string, newItem string) {
	*s = append(*s, newItem)
}

func passByValue(x int) {
	x = 999
}

func passByPointer(x *int) {
	*x = 999
}

// MAIN PROGRAM

func main() {
	fmt.Println("||| POINTER |||")
	fmt.Println("")


	fmt.Println("~~~ 1. Fungsi swap ~~~")
	a, b := 10, 20
	fmt.Printf("Sebelum swap : a = %d, b = %d\n", a, b)
	swap(&a, &b)
	fmt.Printf("Setelah swap  : a = %d, b = %d\n\n", a, b)


	fmt.Println("~~~ 2. Fungsi updateSlice ~~~")
	buah := []string{"Apel", "Jeruk"}
	fmt.Printf("Sebelum update: %v (len=%d, cap=%d)\n", buah, len(buah), cap(buah))
	updateSlice(&buah, "Mangga")
	fmt.Printf("Setelah update: %v (len=%d, cap=%d)\n", buah, len(buah), cap(buah))
	updateSlice(&buah, "Durian")
	fmt.Printf("Setelah tambah Durian: %v (len=%d, cap=%d)\n\n", buah, len(buah), cap(buah))


	fmt.Println("~~~ 3. Perbandingan pass by value vs pass by pointer ~~~")
	angka := 50
	fmt.Printf("Nilai awal angka: %d\n", angka)

	passByValue(angka)
	fmt.Printf("Setelah passByValue (salinan): %d \t(tidak berubah)\n", angka)

	passByPointer(&angka)
	fmt.Printf("Setelah passByPointer (alamat): %d \t(berubah)\n\n", angka)


	fmt.Println("~~~ 4. Bonus: Slice tanpa pointer (hanya mengubah elemen) ~~~")
	slice := []int{1, 2, 3}
	fmt.Printf("Sebelum: %v\n", slice)
	func(s []int) {
		s[0] = 100
	}(slice)
	fmt.Printf("Setelah ubah elemen: %v (berubah, walau tanpa pointer!)\n", slice)
	fmt.Println("\nFungsi yang mencoba append tanpa pointer")
	func(s []int) {
		s = append(s, 4)
		fmt.Printf("Di dalam fungsi: %v\n", s)
	}(slice)
	fmt.Printf("Setelah fungsi keluar: %v (tidak berubah, karena append tidak pakai pointer)\n", slice)
	fmt.Println("\nupdateSlice (pakai pointer) berhasil menambah")
	updateSliceString := func(s *[]string, item string) {
		*s = append(*s, item)
	}
	daftar := []string{"A", "B"}
	fmt.Printf("Sebelum: %v\n", daftar)
	updateSliceString(&daftar, "C")
	fmt.Printf("Setelah: %v (berhasil)\n", daftar)
}
