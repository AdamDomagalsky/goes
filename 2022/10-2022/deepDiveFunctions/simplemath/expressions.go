package simplemath

import "fmt"

func Divide(p1, p2 float64) (result float64, err error) {
	if p2 == 0 {
		err = fmt.Errorf("cannot divide by zero")
	}

	result = p1 / p2
	return
}
