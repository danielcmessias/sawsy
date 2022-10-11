package histogram

func Min(v []float64) float64 {
	var m float64
	for i, e := range v {
		if i == 0 || e < m {
			m = e
		}
	}
	return m
}

func Max(v []float64) float64 {
	var m float64
	for i, e := range v {
		if i == 0 || e > m {
			m = e
		}
	}
	return m
}
