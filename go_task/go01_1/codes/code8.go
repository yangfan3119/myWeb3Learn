package codes

import (
	"fmt"
	"math/rand"
)

/*
两数之和：
给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出和为目标值 target  的那两个整数，
并返回它们的数组下标。你可以假设每种输入只会对应一个答案，并且你不能使用两次相同的元素。
你可以按任意顺序返回答案。

示例 1：
输入：nums = [2,7,11,15], target = 9
输出：[0,1]
解释：因为 nums[0] + nums[1] == 9 ，返回 [0, 1] 。

示例 2：
输入：nums = [3,2,4], target = 6
输出：[1,2]

示例 3：
输入：nums = [3,3], target = 6
输出：[0,1]

提示：
2 <= nums.length <= 104
-109 <= nums[i] <= 109
-109 <= target <= 109
只会存在一个有效答案
*/
var C8_nums = []int{2, 7, 11, 15}
var C8_target = 9

func C8_getTestNum(max int) ([]int, int) {
	// 设置随机数种子（必须在生成随机数前调用，建议只设置一次）
	// rand.Seed(time.Now().UnixNano()) // 以当前纳秒时间为种子

	if max < 20 {
		max = 20
	}
	// 1. 生成Target随机数
	target := (rand.Intn(max*100) + max*25) / 100 // 生成区间的随机整数
	fmt.Println("Target随机数:", target)

	// 2. 生成 [0, max) 范围内的随机整数（自定义范围）,且去重
	resNums := []int{}
	seen := make(map[int]bool)

	for i := max / 4; i > 0; {
		x := rand.Intn(max)
		if !seen[x] {
			seen[x] = true
			resNums = append(resNums, x)
			i--
		}
	}
	fmt.Println("生成 [0, ", max, ") 范围内的随机整数数组:", resNums)
	return resNums, target
}

func TwoSum(nums []int, target int) []int {
	fmt.Println("输入: nums = ", nums, "target = ", target)

	tar := target/2 + 1
	p := make(map[int]int)
	var res [][]int
	for index, x := range nums {
		if x < tar {
			p[target-x] = index
		} else {
			v, exist := p[x]
			if exist {
				fmt.Println("找到一组解", v, index)
				fmt.Println("因为nums[", v, "]+nums[", index, "] = target ", target)
				fmt.Println("即", nums[v], "+ ", nums[index], " = ", target)
				res = append(res, []int{v, index})
			}
		}
	}
	if len(res) > 0 {
		return res[0]
	}
	fmt.Print("当前nums中没有和为target的两个数。")
	return []int{}
}
