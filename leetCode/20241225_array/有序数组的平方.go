package main

import "fmt"

func sortedSquares(nums []int) []int { //暴力法和双指针法均自解，如果涉及中间元素未处理的话直接用左闭右闭
	left, right, l := 0, len(nums)-1, len(nums)-1
	res := make([]int, l+1)
	for left <= right { //要考虑取中间值的话就左闭右闭
		if nums[left]*nums[left] >= nums[right]*nums[right] {
			res[l] = nums[left] * nums[left]
			left++
		} else {
			res[l] = nums[right] * nums[right]
			right--
		}

		l--
	}
	return res
}

func main() {
	nums := []int{-4, -1, 0, 3, 10}
	fmt.Println(sortedSquares(nums))
}
