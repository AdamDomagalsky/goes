package main

import "fmt"

func isValid(s string) bool {
	stack := []int{}
	for i := 0; i < len(s); i++ {
		switch elo := s[i]; elo {
		case 40:
			stack = append(stack, 40)
		case 41:
			stack = stack[:len(stack)-1]
		case 91:
			stack = append(stack, 91)
		case 123:
			stack = append(stack, 123)
		default:
			// fmt.Print("eloo", elo)
		}
	}
	fmt.Print(stack)

	return true
}

func main() {
	fmt.Print(isValid("()[]{}"))

	// stack := []string{}
	// stack = append(stack, "e")
	// fmt.Print(stack)
	// stack = stack[:len(stack)-1]
	// fmt.Print(stack)

}
