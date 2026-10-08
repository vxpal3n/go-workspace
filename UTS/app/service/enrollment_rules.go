package service

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseTahunAkademik(value string) error {
	parts := strings.SplitN(value, "-", 2)
	if len(parts) != 2 {
		return fmt.Errorf("format tahun akademik tidak valid")
	}

	semester := parts[1]
	if semester != "Ganjil" && semester != "Genap" {
		return fmt.Errorf("semester harus Ganjil atau Genap")
	}

	years := strings.SplitN(parts[0], "/", 2)
	if len(years) != 2 {
		return fmt.Errorf("format tahun tidak valid")
	}
	y1, err := strconv.Atoi(years[0])
	if err != nil {
		return fmt.Errorf("tahun pertama bukan angka")
	}
	y2, err := strconv.Atoi(years[1])
	if err != nil {
		return fmt.Errorf("tahun kedua bukan angka")
	}
	if y1 < 2000 || y1 > 2100 {
		return fmt.Errorf("tahun di luar rentang wajar")
	}
	if y2 != y1+1 {
		return fmt.Errorf("tahun kedua harus tahun pertama + 1")
	}

	return nil
}

func CanEnrollSKS(totalSKS, sksBaru, batas int) error {
	if totalSKS+sksBaru > batas {
		sisa := batas - totalSKS
		if sisa < 0 {
			sisa = 0
		}
		return fmt.Errorf(
			"total SKS melebihi batas %d. Sisa SKS yang dapat diambil: %d",
			batas, sisa)
	}
	return nil
}