package service

import "testing"

func TestBatasSKS(t *testing.T) {
	cases := []struct {
		ipk  float64
		want int
	}{
		{4.00, 24},  // batas atas
		{3.50, 24},  // di atas 3.00
		{3.00, 24},  // tepat di batas
		{2.99, 21},  // turun satu tier
		{2.75, 21},
		{2.50, 21},  // tepat di batas bawah
		{2.49, 18},  // turun ke tier terbawah
		{2.00, 18},
		{0.00, 18},
	}
	for _, tc := range cases {
		if got := BatasSKS(tc.ipk); got != tc.want {
			t.Errorf("BatasSKS(%.2f): harap %d, dapat %d", tc.ipk, tc.want, got)
		}
	}
}

func TestTotalSKS(t *testing.T) {
	cases := []struct {
		name string
		sks  []int
		want int
	}{
		{"kosong", []int{}, 0},
		{"satu elemen", []int{3}, 3},
		{"beberapa elemen", []int{3, 3, 4, 2}, 12},
		{"banyak elemen", []int{3, 3, 3, 3, 3, 3, 3, 3}, 24},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := TotalSKS(tc.sks); got != tc.want {
				t.Errorf("harap %d, dapat %d", tc.want, got)
			}
		})
	}
}