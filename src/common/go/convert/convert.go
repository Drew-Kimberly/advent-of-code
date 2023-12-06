package convert

import "strconv"

func MustBeInt(val string) int {
	num, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return num
}
