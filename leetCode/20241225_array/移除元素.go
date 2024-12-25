package main

import "fmt"

func removeElement(nums []int, val int) int { //暴力法个人要理解size的定义（size--，return size）；双指针两种均没有问题（快慢指针以及左闭右闭双向）
	slow := 0
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
	nums := []int{3, 2, 2, 3}
	fmt.Println(removeElement(nums, 3))
}
