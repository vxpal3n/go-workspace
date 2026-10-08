package service

func BatasSKS(ipk float64) int {
	switch {
	case ipk >= 3.00:
		return 24
	case ipk >= 2.50:
		return 21
	default:
		return 18
	}
}

func TotalSKS(sks []int) int {
	total := 0
	for _, v := range sks {
		total += v
	}
	return total
}