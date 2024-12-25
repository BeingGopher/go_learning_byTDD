package main

import "fmt"

func search(nums []int, target int) int { //双指针解决二分查找，注意左闭右开和左闭右闭！
	left := 0
	right := len(nums)
	for left < right {
		mid := (right + left) / 2
		if nums[mid] < target {
			left = mid + 1
		} else if nums[mid] > target {
			right = mid
		} else {
			return mid
		}
	}
	return -1
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(search(nums, 9))
}
