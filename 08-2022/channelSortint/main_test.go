package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

const TEST_SLICE_LENGTH = 500_00
const RANDOM_MODULO_RANGE = 300_000

type args struct {
	data       []int
	partitions uint
	sortAlgo   func([]int)
	isSorted   bool
}

func allAvailableSortingFunctions() map[string]func([]int) {
	return map[string]func([]int){
		"Insert": InsertSort,
		"Bubble": BubbleSort,
		"Select": SelectSort,
	}
}

func getSortingFunctionBaseOnName(name string) (func([]int), bool) {
	elem, ok := allAvailableSortingFunctions()[name]
	return elem, ok
}

var testCases = map[string]args{
	"SampleSort5": {
		data:       generateRandomSlice(TEST_SLICE_LENGTH, RANDOM_MODULO_RANGE),
		partitions: 333,
		sortAlgo:   InsertSort,
		isSorted:   false,
	},
}

func addTestCase(benchmarkFunctionName string) args {

	if !strings.Contains(benchmarkFunctionName, "Sort") {
		panic("benchmarkFunctionName not following convention \nHaving: " + benchmarkFunctionName + "\nExpected: <SortAlgoName>Sort<partitionNumber>")
	}
	typeAndNumber := strings.Split(benchmarkFunctionName, "Sort")
	chosenSortAlgo, ok := getSortingFunctionBaseOnName(typeAndNumber[0])
	if !ok {
		keys := make([]string, 0, len(allAvailableSortingFunctions()))
		for k := range allAvailableSortingFunctions() {
			keys = append(keys, k)
		}
		fmt.Println(keys)
		panic("Chosen sorting algorithm not implemented.\n\tHaving: " + typeAndNumber[0] + "\n\tExpected: <SortAlgoName>Sort<partitionNumber> where <SortAlgoName> OneOf(" + strings.Join(keys, ",") + ")")
	}

	atoi, err := strconv.Atoi(typeAndNumber[1])
	if err != nil {
		panic("\n\tCould not parse partitionNumber from benchmarkFunctionName \n\tHaving: " + typeAndNumber[1] + "\n\tExpected: <SortAlgoName>Sort<partitionNumber>")
	}

	testCases[benchmarkFunctionName] = args{
		data:       generateRandomSlice(TEST_SLICE_LENGTH, RANDOM_MODULO_RANGE),
		partitions: uint(atoi),
		sortAlgo:   chosenSortAlgo,
		isSorted:   false,
	}
	return testCases[benchmarkFunctionName]
}

func benchmarkSortWithGoRoutinePartition(args args, b *testing.B) {
	for n := 0; n < b.N; n++ {
		args.data = generateRandomSlice(TEST_SLICE_LENGTH, RANDOM_MODULO_RANGE)
		SortWithGoRoutinePartition(args.data, args.partitions, args.sortAlgo)
	}
}

func BenchmarkInsertSort1(b *testing.B) {
	benchmarkFunctionName := "InsertSort1"
	addTestCase(benchmarkFunctionName)
	benchmarkSortWithGoRoutinePartition(addTestCase(benchmarkFunctionName), b)
}

func BenchmarkInsertSort4(b *testing.B) {
	benchmarkFunctionName := "InsertSort4"
	addTestCase(benchmarkFunctionName)
	benchmarkSortWithGoRoutinePartition(addTestCase(benchmarkFunctionName), b)
}

func BenchmarkInsertSort10(b *testing.B) {
	benchmarkFunctionName := "InsertSort10"
	addTestCase(benchmarkFunctionName)
	benchmarkSortWithGoRoutinePartition(addTestCase(benchmarkFunctionName), b)
}

func BenchmarkInsertSort20(b *testing.B) {
	benchmarkFunctionName := "InsertSort20"
	addTestCase(benchmarkFunctionName)
	benchmarkSortWithGoRoutinePartition(addTestCase(benchmarkFunctionName), b)
}

func BenchmarkBubbleSort1(b *testing.B) {
	benchmarkFunctionName := "BubbleSort1"
	addTestCase(benchmarkFunctionName)
	benchmarkSortWithGoRoutinePartition(addTestCase(benchmarkFunctionName), b)
}

func BenchmarkBubbleSort4(b *testing.B) {
	benchmarkFunctionName := "BubbleSort4"
	addTestCase(benchmarkFunctionName)
	benchmarkSortWithGoRoutinePartition(addTestCase(benchmarkFunctionName), b)
}

func BenchmarkBubbleSort10(b *testing.B) {
	benchmarkFunctionName := "BubbleSort10"
	addTestCase(benchmarkFunctionName)
	benchmarkSortWithGoRoutinePartition(addTestCase(benchmarkFunctionName), b)
}

func BenchmarkBubbleSort20(b *testing.B) {
	benchmarkFunctionName := "BubbleSort20"
	addTestCase(benchmarkFunctionName)
	benchmarkSortWithGoRoutinePartition(addTestCase(benchmarkFunctionName), b)
}

func BenchmarkSelectSort1(b *testing.B) {
	benchmarkFunctionName := "SelectSort1"
	addTestCase(benchmarkFunctionName)
	benchmarkSortWithGoRoutinePartition(addTestCase(benchmarkFunctionName), b)
}

func BenchmarkSelectSort4(b *testing.B) {
	benchmarkFunctionName := "SelectSort4"
	addTestCase(benchmarkFunctionName)
	benchmarkSortWithGoRoutinePartition(addTestCase(benchmarkFunctionName), b)
}

func BenchmarkSelectSort10(b *testing.B) {
	benchmarkFunctionName := "SelectSort20"
	addTestCase(benchmarkFunctionName)
	benchmarkSortWithGoRoutinePartition(addTestCase(benchmarkFunctionName), b)
}

func BenchmarkSelectSort20(b *testing.B) {
	benchmarkFunctionName := "SelectSort20"
	addTestCase(benchmarkFunctionName)
	benchmarkSortWithGoRoutinePartition(addTestCase(benchmarkFunctionName), b)
}

//func BenchmarkInsertSortX(b *testing.B) {
//	benchmarkFunctionName := "InsessrtSort455"
//	addTestCase(benchmarkFunctionName)
//	benchmarkSortWithGoRoutinePartition(addTestCase(benchmarkFunctionName), b)
//}

//func BenchmarkInsertSortY(b *testing.B) {
//	benchmarkFunctionName := "InsertSort4X55"
//	addTestCase(benchmarkFunctionName)
//	benchmarkSortWithGoRoutinePartition(addTestCase(benchmarkFunctionName), b)
//}
