package main

import "fmt"

func removeElement(nums []int, val int) int {
	// 初始化慢指针 slow
	slow := 0
	// 通过 for 循环移动快指针 fast
	// 当 fast 指向的元素等于 val 时，跳过
	// 否则，将该元素写入 slow 指向的位置，并将 slow 后移一位
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] == val {
			continue
		}
		nums[slow] = nums[fast]
		slow++
	}

	return slow
}

func main() {
	fmt.Println(removeElement([]int{3, 2, 2, 3}, 3))
}
