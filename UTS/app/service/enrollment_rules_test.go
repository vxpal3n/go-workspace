package service

import "testing"

func TestParseTahunAkademik_Valid(t *testing.T) {
	valid := []string{
		"2026/2027-Ganjil",
		"2026/2027-Genap",
		"2025/2026-Ganjil",
		"2024/2025-Genap",
		"2100/2101-Ganjil",
	}
	for _, v := range valid {
		if err := ParseTahunAkademik(v); err != nil {
			t.Errorf("harap valid: %q, dapat error: %v", v, err)
		}
	}
}

func TestParseTahunAkademik_Invalid(t *testing.T) {
	cases := []struct {
		name  string
		value string
	}{
		{"tanpa pemisah tahun", "2026-2027-Ganjil"},
		{"tanpa semester", "2026/2027"},
		{"semester salah", "2026/2027-Summer"},
		{"tahun kedua bukan +1", "2026/2028-Ganjil"},
		{"tahun pertama bukan angka", "abcd/2027-Ganjil"},
		{"tahun kedua bukan angka", "2026/abcd-Ganjil"},
		{"tahun di luar rentang", "1999/2000-Ganjil"},
		{"kosong", ""},
		{"hanya separator", "/-"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := ParseTahunAkademik(tc.value); err == nil {
				t.Errorf("harap error untuk %q, tapi lolos", tc.value)
			}
		})
	}
}

func TestCanEnrollSKS(t *testing.T) {
	cases := []struct {
		name     string
		total    int
		baru     int
		batas    int
		wantErr  bool
	}{
		{"kosong + 3 dari 24", 0, 3, 24, false},
		{"21 + 3 dari 24", 21, 3, 24, false},   // tepat di batas
		{"22 + 3 dari 24", 22, 3, 24, true},    // kelebihan 1
		{"18 + 3 dari 21", 18, 3, 21, false},
		{"19 + 3 dari 21", 19, 3, 21, true},
		{"15 + 3 dari 18", 15, 3, 18, false},
		{"16 + 3 dari 18", 16, 3, 18, true},
		{"4 sks di batas 24", 21, 4, 24, true}, // +4 dari 21 > 24
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := CanEnrollSKS(tc.total, tc.baru, tc.batas)
			if tc.wantErr && err == nil {
				t.Errorf("harap error, tapi lolos (total=%d, baru=%d, batas=%d)",
					tc.total, tc.baru, tc.batas)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("harap lolos, tapi error: %v", err)
			}
		})
	}
}