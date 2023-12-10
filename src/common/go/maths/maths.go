package maths

import "errors"

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(integers ...int) int {
	if len(integers) < 2 {
		panic(errors.New("must provide at least 2 integer params to find LCM"))
	}

	a := integers[0]
	b := integers[1]

	result := a * b / GCD(a, b)

	for _, val := range integers[2:] {
		result = LCM(result, val)
	}

	return result
}
