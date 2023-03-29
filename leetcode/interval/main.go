package main

import "fmt"

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// https://leetcode.com/problems/insert-interval/
//PRO: https://leetcode.com/problems/insert-interval/submissions/924232610/

func insert(intervals [][]int, newInterval []int) [][]int {
	start, end := newInterval[0], newInterval[1]
	length := len(intervals)
	result := [][]int{}
	i := 0

	for ; i < length && intervals[i][1] < start; i++ {
		result = append(result, intervals[i])
	}

	for ; i < length && intervals[i][0] <= end; i++ {
		start = min(start, intervals[i][0])
		end = max(end, intervals[i][1])
	}

	return append(append(result, []int{start, end}), intervals[i:]...)
}

func insertBoss(intervals [][]int, newInterval []int) [][]int {

	res := make([][]int, 0)

	for i, interval := range intervals {
		curL, curR := interval[0], interval[1]
		newL, newR := newInterval[0], newInterval[1]

		if newR < curL {
			res = append(res, newInterval)
			res = append(res, intervals[i:]...)
			return res
		}

		if newL > curR {
			res = append(res, interval)
		} else {
			newInterval[0], newInterval[1] = min(curL, newL), max(curR, newR)
		}
	}
	res = append(res, newInterval)
	return res
}

func main() {

	fmt.Print(insert([][]int{{1, 3}, {6, 7}, {9, 11}, {13, 16}}, []int{4, 8}))
}
