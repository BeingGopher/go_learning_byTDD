package main

import "fmt"

// 题解版本，书写可能更规范
// 时间复杂度 O(logn)
func search(nums []int, target int) int {
	// 初始化左右边界
	left := 0
	right := len(nums) - 1

	// 循环逐步缩小区间范围
	for left <= right {
		// 求区间中点
		mid := left + (right-left)>>1

		// 根据 nums[mid] 和 target 的大小关系
		// 调整区间范围
		if nums[mid] == target {
			return mid
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	// 在输入数组内没有找到值等于 target 的元素
	return -1
}

func main() {
	arr := []int{2, 3, 4, 10, 40}
	target := 10

	result := search(arr, target)

	fmt.Println(result)

}
