package codes

/*
给定一个表示 大整数 的整数数组 digits，其中 digits[i] 是整数的第 i 位数字。
这些数字按从左到右，从最高位到最低位排列。这个大整数不包含任何前导 0。
将大整数加 1，并返回结果的数字数组。

示例 1：
输入：digits = [1,2,3]
输出：[1,2,4]
解释：输入数组表示数字 123。
加 1 后得到 123 + 1 = 124。
因此，结果应该是 [1,2,4]。
*/
var C5_t_digits = []int{9, 9, 9, 9}

func PlusOne(digits []int) []int {
	var base int = 10
	var PlusNum int = 1
	lenDig := len(digits) - 1
	digits[lenDig] = digits[lenDig] + PlusNum
	for ; lenDig >= 0; lenDig-- {
		if digits[lenDig] >= base {
			digits[lenDig] = 0
			if lenDig == 0 {
				digits = append([]int{1}, digits...)
			} else {
				digits[lenDig-1] = digits[lenDig-1] + 1
			}
		} else { // 无需进位直接退出
			break
		}
	}
	return digits
}
