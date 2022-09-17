package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

// Write a program to sort2 an array of integers. The program should partition the array into 4 parts, each of which is sorted by a different goroutine. Each partition should be of approximately equal size. Then the main goroutine should merge the 4 sorted subarrays into one large sorted array.
// The program should prompt the user to input a series of integers. Each goroutine which sorts ¼ of the array should print the subarray that it will sort2. When sorting is complete, the main goroutine should print the entire sorted list.
func IsSorted(data []int) bool {
	for i := 0; i < len(data)-1; i++ {
		if data[i] > data[i+1] {
			return false
		}
	}
	return true
}

func InsertSort(slice []int) {
	for i := 1; i < len(slice); i++ {
		key := slice[i]
		j := i - 1
		for j >= 0 && slice[j] > key {
			slice[j+1] = slice[j]
			j -= 1
		}
		slice[j+1] = key
	}
}

func sortWG(sortAlgo func([]int), slice []int, wg *sync.WaitGroup) {
	sortAlgo(slice)
	wg.Done()
}

func generateRandomSlice(size uint, mod int) []int {
	arr := make([]int, size)
	for i, _ := range arr {
		arr[i] = rand.Int() % mod
	}
	return arr
}

func SortWithGoRoutinePartition(data []int, partitions uint, sortAlgo func([]int)) {
	wg := &sync.WaitGroup{}
	for pStart, pEnd, i := uint(0), uint(0), uint(1); i <= partitions; pStart, i = pEnd, i+1 {
		pEnd = uint(len(data)) / partitions * i
		wg.Add(1)
		fmt.Println("sorting [", pStart, ",", pEnd, "]: ", data[pStart:pEnd])
		go sortWG(sortAlgo, data[pStart:pEnd], wg)
		pStart = pEnd
	}
	wg.Wait()
	sortAlgo(data)
}

func main() {

	var numbers string
	fmt.Println("Provide coma separated numbers: ")
	fmt.Scan(&numbers)
	var data []int
	for _, element := range strings.Split(numbers, ",") {
		number, err := strconv.Atoi(element)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, number)
	}
	fmt.Println(data)
	SortWithGoRoutinePartition(data, 4, InsertSort)
	fmt.Println(data)
}
