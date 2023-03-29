package main

import (
	"fmt"
)

// https://leetcode.com/submissions/detail/924112130/

func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	var counter [26]int
	for i := 0; i < len(s); i++ {
		counter[int(s[i])-int('a')]++
		counter[int(t[i])-int('a')]--
	}

	for i := 0; i < len(counter); i++ {
		if counter[i] != 0 {
			return false
		}
	}

	return true
}

func main() {
	fmt.Println(isAnagram("anagram", "nagaram"))
	fmt.Println(isAnagram("anagram", "nagarams"))
	fmt.Println(isAnagram("anagramx", "nagarams"))
	fmt.Println(isAnagram("aacc", "ccac"))
}
