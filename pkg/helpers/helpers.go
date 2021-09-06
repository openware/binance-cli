package helpers

import "encoding/json"

func ValuePrecision(number json.Number) int {
	n := 0
	if v, _ := number.Float64(); v == 0 {
		return 0
	} else if v < 1 {
		n = 1
		for v * 10 < 1 {
			n +=1
			v *= 10
		}

	} else if v < 10 {
		n = 0
	} else {
		n = -1
		for v / 10 > 1 {
			n-=1
		}
	}
	return n
}
