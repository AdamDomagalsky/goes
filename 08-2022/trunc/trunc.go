package main

import (
	"fmt"
	"strconv"
)

// Write a program which prompts the user to enter a floating point number and prints the integer which is a truncated version of the floating point number that was entered.
// Truncation is the process of removing the digits to the right of the decimal place.

func main() {
	var input string
	fmt.Print("Enter a floating point number: ")
	fmt.Scan(&input)
	fl, err := strconv.ParseFloat(input, 2)
	if err != nil {
		fmt.Print(err.Error())
	}

	fmt.Print(int(fl))
}
