package main

import "fmt"

// STRUCT STUDENT

type Student struct {
	ID       int
	Name     string
	Grade    float64
	IsActive bool
}

// METHOD-METHOD

func (s Student) GetInfo() string {
	status := "Tidak Aktif"
	if s.IsActive {
		status = "Aktif"
	}
	return fmt.Sprintf("ID: %d | Nama: %s | Nilai: %.2f | Status: %s",
		s.ID, s.Name, s.Grade, status)
}

func (s *Student) UpdateGrade(grade float64) {
	s.Grade = grade
}

func (s *Student) Activate() {
	s.IsActive = true
}

func (s *Student) Deactivate() {
	s.IsActive = false
}


func tryUpdateWithValue(s Student, newGrade float64) {
	s.Grade = newGrade
	fmt.Printf("(Di dalam fungsi value) grade diubah menjadi %.2f\n", s.Grade)
}

// MAIN PROGRAM

func main() {
	fmt.Println("||| STRUCT & METHOD |||")

	mhs := Student{
		ID: 2024,
		Name: "Valentino",
		Grade: 80.5,
		IsActive: false,
	}

	fmt.Println("~~~ Data Awal ~~~")
	fmt.Println(mhs.GetInfo())
	fmt.Printf("Alamat struct di memori: %p\n\n", &mhs)

	fmt.Println("~~~ Mengaktifkan dan mengupdate nilai (pointer receiver) ~~~")
	mhs.Activate()
	mhs.UpdateGrade(92.0)
	fmt.Println("Setelah Activate() dan UpdateGrade(92.0):")
	fmt.Println(mhs.GetInfo())
	fmt.Println()

	fmt.Println("~~~ Menonaktifkan (Deactivate) ~~~")
	mhs.Deactivate()
	fmt.Println("Setelah Deactivate():")
	fmt.Println(mhs.GetInfo())
	fmt.Println()

	fmt.Println("~~~ Perbandingan: Value Receiver vs Pointer Receiver ~~~")
	fmt.Println("")
	fmt.Println(">> Memanggil tryUpdateWithValue (pass by value):")
	tryUpdateWithValue(mhs, 100.0)
	fmt.Println("Setelah fungsi selesai, grade asli tetap:", mhs.Grade)
	fmt.Println("(Karena fungsi hanya bekerja pada salinan)")

	fmt.Println("\n>> Memanggil UpdateGrade (pointer receiver):")
	mhs.UpdateGrade(100.0)
	fmt.Println("Setelah UpdateGrade, grade asli menjadi:", mhs.Grade)
	fmt.Println("(Karena pointer receiver mengubah instance asli)")

	fmt.Println("\n~~~ GetInfo menggunakan value receiver ~~~")
	info := mhs.GetInfo()
	fmt.Println("Info:", info)
	fmt.Println("Struct asli tetap sama:", mhs)
}