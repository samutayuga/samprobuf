package decomposer

import (
	"math"
)

func IsPrime(aNumber int) bool {
	if aNumber == 2 {
		return true
	}
	//if even
	if aNumber > 2 && aNumber%2 == 0 {
		return false
	}
	//if odd
	for i := 3; float64(i) <= math.Sqrt(float64(aNumber)); i += 2 {
		if aNumber%i == 0 {
			return false
		}
	}
	return true
}
