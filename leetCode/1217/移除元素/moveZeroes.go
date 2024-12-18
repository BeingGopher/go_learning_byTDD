package main

import "fmt"

// 自己的解法
func moveZeroes(nums []int) []int {
	if len(nums) == 1 {
		return nums //定义了返回值的情况下，可以使用return

	}

	slow := 0
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != 0 {
			nums[slow], nums[fast] = nums[fast], nums[slow]
			slow++
		}
	}
	return nums

}

func main() {
	nums := []int{0, 1, 0, 3, 12}
	fmt.Println(moveZeroes(nums))
}

/*注：打印和return的区别如下
     1.打印限制了函数的用途，因为输出是一次性的，并且不能被程序的其他部分重用。
     2.返回值可以被程序的其他部分重用，并且可以返回多个值。
	 3.如果需要使用return，需要定义返回值类型，不然在GO中会报错。
	 4.并且，return在性能上优于打印，打印操作涉及到字符串格式化和系统调用，这可能会稍微慢一些，尤其是在打印大量数据时。
*/
