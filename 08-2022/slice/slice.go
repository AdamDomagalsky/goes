package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Write a program which prompts the user to enter integers and stores the integers in a sorted slice.
// The program should be written as a loop.
// Before entering the loop, the program should create an empty integer slice of size (length) 3.
// During each pass through the loop, the program prompts the user to enter an integer to be added to the slice.
// The program adds the integer to the slice, sorts the slice, and prints the contents of the slice in sorted order.
// The slice must grow in size to accommodate any number of integers which the user decides to enter.
// The program should only quit (exiting the loop) when the user enters the character ‘X’ instead of an integer.

func main() {
	slice := make([]int, 0, 3)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter an integer: ")
	for scanner.Scan() {
		text := scanner.Text()
		if strings.ToLower(text) == "x" {
			fmt.Println("‘X’ program exit")
			return
		}

		number, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		slice = insertNumberToSortedSlice(&slice, number)
		fmt.Println(slice)
		fmt.Print("Enter an integer: ")
	}

}

func insertNumberToSortedSlice(slice *[]int, number int) []int {

	sliceLength := len(*slice)
	destinationIndex := 0
	for i := 0; i < sliceLength; i++ {
		// fmt.Print(" slice[", i, "]:", (*slice)[i], " number:", number, "\n")
		if number > (*slice)[i] {
			destinationIndex = i + 1
		} else {
			break
		}
	}

	if destinationIndex == 0 { // head of slice
		return append([]int{number}, (*slice)...)
	}

	if destinationIndex == sliceLength { // tail of slice
		return append((*slice), number)
	}

	// insert in the middle of the slice
	a := append((*slice)[:destinationIndex+1], (*slice)[destinationIndex:]...)
	a[destinationIndex] = number
	return a
}
