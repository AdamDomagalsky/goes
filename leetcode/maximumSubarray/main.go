package main

import "fmt"

// https://leetcode.com/problems/maximum-subarray/submissions/924264272/
// kadane's algorithm
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxSubArray(nums []int) int {
	overallMaxSum := nums[0]
	currentMaxSum := nums[0]

	for i := 1; i < len(nums); i++ {
		currentMaxSum = max(nums[i], currentMaxSum+nums[i])
		overallMaxSum = max(currentMaxSum, overallMaxSum)
	}
	return overallMaxSum
}

func main() {
	fmt.Print(maxSubArray([]int{5, 4, -1, 7, 8}))
	// fmt.Print(maxSubArray([]int{1}))
	// fmt.Print(maxSubArray([]int{-2, 1, -3, 4, -1, 2, 1, -5, 4}))
}
