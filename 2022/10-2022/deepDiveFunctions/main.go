package main

import "goes/10-2022/deepDiveFunctions/simplemath"

func main() {
	s := []float64{12.5, 12.5, 14.3, -5.0}
	total := sum(s...)
	println(total)

	sv := simplemath.NewSematicVersion(1, 2, 3)
	println(sv.String())

}

func sum(value ...float64) float64 {
	total := 0.0
	for _, value := range value {
		total += value
	}
	return total
}
