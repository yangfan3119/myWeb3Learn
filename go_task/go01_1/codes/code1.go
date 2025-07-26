package codes

import "fmt"

/*
136. 只出现一次的数字：给定一个非空整数数组，除了某个元素只出现一次以外，
其余每个元素均出现两次。找出那个只出现了一次的元素。
可以使用 for 循环遍历数组，结合 if 条件判断和 map 数据结构来解决，
例如通过 map 记录每个元素出现的次数，然后再遍历 map 找到出现次数为1的元素。
*/
var C_hello = "hello go World!"

var C1_testNum1 []int = []int{1, 2, 3, 4, 5, 6, 7, 2, 5, 4, 6}

func OnceNum(ns []int) {
	nCount := make(map[int]int)
	for _, x := range ns {
		v, exists := nCount[x]
		if exists {
			// 值存在则+1
			nCount[x] = v + 1
		} else {
			// 值不存在则写入map，值为1
			nCount[x] = 1
		}
	}

	for k, v := range nCount {
		if v == 1 {
			fmt.Println("num ", k, "appears only once")
		}
	}

}
