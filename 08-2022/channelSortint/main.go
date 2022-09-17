package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"
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

func BubbleSort(aSlice []int) {
	for j := 0; j < len(aSlice)-1; j++ {
		for i := 0; i < len(aSlice)-1; i++ {
			if aSlice[i] > aSlice[i+1] {
				//swap(aSlice, i)
				aSlice[i], aSlice[i+1] = aSlice[i+1], aSlice[i]
			}
		}
	}
}

func swap(aSlice []int, i int) {
	tmp := aSlice[i]
	aSlice[i] = aSlice[i+1]
	aSlice[i+1] = tmp
}

func SelectSort(slice []int) {
	for i, _ := range slice {
		minIndex := i
		for j := i + 1; j < len(slice); j++ {
			if slice[minIndex] > slice[j] {
				minIndex = j
			}
		}
		slice[minIndex], slice[i] = slice[i], slice[minIndex]
	}
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
		go sortWG(sortAlgo, data[pStart:pEnd], wg)
		pStart = pEnd
	}
	wg.Wait()
	sortAlgo(data)
}

func timeTrack(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%v\t| %v ns\n", name, time.Since(start).Nanoseconds())
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func performTest(sortAlgo func([]int), size, mod, goRoutinePartitions uint, debug bool) {
	arr := generateRandomSlice(size, int(mod))
	if debug {
		fmt.Println("------start------")
		fmt.Println(IsSorted(arr))
	}
	partition := int(goRoutinePartitions)
	defer timeTrack(fmt.Sprintf("%s %d random numbers modulo %d, %d goRoutinePartitions", getFunctionName(sortAlgo), size, mod, goRoutinePartitions))()
	SortWithGoRoutinePartition(arr, uint(partition), sortAlgo)
	if debug {
		if len(arr) < 100 {
			fmt.Println(arr)
		} else {
			fmt.Println(arr[:10], "...", arr[len(arr)-10:])
		}
		fmt.Println("------end------")
	}
}

func main() {

	//doo := generateRandomSlice(10, 100)
	//fmt.Println(doo)
	//fmt.Println(IsSorted(doo))
	//
	//SortWithGoRoutinePartition(doo, 4, SelectSort)
	//fmt.Println(doo)
	//fmt.Println(IsSorted(doo))
	const TEST_SLICE_LENGTH = 500_00
	const RANDOM_MODULO_RANGE = 300_000
	fmt.Println("InsertSort sortWG:")
	performTest(InsertSort, TEST_SLICE_LENGTH, RANDOM_MODULO_RANGE, 1, false)
	performTest(InsertSort, TEST_SLICE_LENGTH, RANDOM_MODULO_RANGE, 4, false)
	performTest(InsertSort, TEST_SLICE_LENGTH, RANDOM_MODULO_RANGE, 10, false)
	performTest(InsertSort, TEST_SLICE_LENGTH, RANDOM_MODULO_RANGE, 20, false)
	fmt.Println("bubble sortWG:")
	performTest(BubbleSort, TEST_SLICE_LENGTH, RANDOM_MODULO_RANGE, 1, false)
	performTest(BubbleSort, TEST_SLICE_LENGTH, RANDOM_MODULO_RANGE, 4, false)
	performTest(BubbleSort, TEST_SLICE_LENGTH, RANDOM_MODULO_RANGE, 10, false)
	performTest(BubbleSort, TEST_SLICE_LENGTH, RANDOM_MODULO_RANGE, 20, false)
	fmt.Println("SelectSort sortWG:")
	performTest(SelectSort, TEST_SLICE_LENGTH, RANDOM_MODULO_RANGE, 1, false)
	performTest(SelectSort, TEST_SLICE_LENGTH, RANDOM_MODULO_RANGE, 4, false)
	performTest(SelectSort, TEST_SLICE_LENGTH, RANDOM_MODULO_RANGE, 10, false)
	performTest(SelectSort, TEST_SLICE_LENGTH, RANDOM_MODULO_RANGE, 20, false)
}
