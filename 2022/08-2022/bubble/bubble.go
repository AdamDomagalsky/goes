package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {

	var userInput []int
	fmt.Print("Enter up to 10 numbers.\nX - to stop entering and start bubble sort\n")

	for {
		if len(userInput) == 10 {
			break
		}
		var inString string
		fmt.Print("Enter a number: ")
		fmt.Scan(&inString)
		if strings.ToLower(inString) == "x" {
			break
		}
		number, err := strconv.Atoi(inString)
		if err != nil {
			fmt.Println(err.Error())
		}
		userInput = append(userInput, number)
	}

	BubbleSort(userInput)
	fmt.Println(userInput)
}

func BubbleSort(aSlice []int) {
	for j := 0; j < len(aSlice)-1; j++ {
		for i := 0; i < len(aSlice)-1; i++ {
			if aSlice[i] > aSlice[i+1] {
				Swap(aSlice, i)
			}
		}
	}
}

func Swap(aSlice []int, i int) {
	tmp := aSlice[i]
	aSlice[i] = aSlice[i+1]
	aSlice[i+1] = tmp
}
