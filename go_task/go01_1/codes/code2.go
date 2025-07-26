package codes

import (
	"fmt"
	"strconv"
)

/*
回文数
考察：数字操作、条件判断
题目：判断一个整数是否是回文数
*/
var C1_test = 123456654321

func Palindrome(x int) {
	fmt.Println("回文数验证：", x)

	sx := strconv.Itoa(x)
	bx := []byte(sx)
	b1 := make([]byte, len(bx), len(bx)+1)
	j := 0
	for i := len(bx) - 1; i >= 0; i-- {
		b1[j] = bx[i]
		j++
	}
	x2, err := strconv.Atoi(string(b1))
	if err != nil {
		panic(err)
	}
	if x == x2 {
		fmt.Println("数字", x, "是回文数")
	} else {
		fmt.Println("数字", x, "不是回文数")
	}
}
