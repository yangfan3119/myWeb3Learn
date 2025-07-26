package main

import (
	"fmt"
	c "go01_1/codes"
)

func main() {
	fmt.Println(c.C_hello)

	// 1.
	c.OnceNum(c.C1_testNum1)

	// 2. 回文数
	c.Palindrome(c.C1_test)

	// 3. 有效的括号
	c.ValidBracket(c.C3_testStr)

	// 4. 最大公共前缀,测试数据str1
	c.MaxPrefix(c.C4_str1)

	// 5. 加一，测试数据 t_digits
	fmt.Println("PlusOne Before:", c.C5_t_digits)
	res := c.PlusOne(c.C5_t_digits)
	fmt.Println("PlusOne After:", res)

	// 6. 删除有序数组中的重复项
	c.RemoveDuplicates1(c.C6_t2)

	// 7. 合并区间
	c.MergeNums(c.C7_t3)
	// c.C7_marge(c.C7_t3)
	// c.C7_ft1()

	// 8. 两数之和
	// c.TwoSum(c.C8_nums, c.C8_target)     //使用固定测试数据测试
	nums, target := c.C8_getTestNum(500) //随机生成数列产生测试数据
	c.TwoSum(nums, target)
}
