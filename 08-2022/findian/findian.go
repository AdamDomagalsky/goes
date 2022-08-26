package main

import (
	"fmt"
	"strings"
)

// Write a program which prompts the user to enter a string.
// The program searches through the entered string for the characters ‘i’, ‘a’, and ‘n’.
// The program should print “Found!” if the entered string starts with the character ‘i’, ends with the character ‘n’, and contains the character ‘a’.
// The program should print “Not Found!” otherwise.
// The program should not be case-sensitive, so it does not matter if the characters are upper-case or lower-case.

func main() {

	var inputString string
	fmt.Print("Enter a string: ")
	fmt.Scan(&inputString)
	inputString = strings.ToLower(inputString)

	if strings.HasPrefix(inputString, "i") && strings.HasSuffix(inputString, "n") && strings.Contains(inputString, "a") {
		fmt.Println("Found!")
	} else {
		fmt.Println("Not Found!")
	}

}
