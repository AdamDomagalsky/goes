package main

import "fmt"

// https://leetcode.com/problems/binary-search/

func search(nums []int, target int) int {
	left := 0
	right := len(nums) - 1

	for left <= right {
		mid := left + (right-left)/2

		if nums[mid] == target {
			return mid
		}

		if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1
}

func main() {
	fmt.Print(search([]int{-1, 0, 3, 5, 9, 12}, 12))
}
