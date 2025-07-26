package codes

import "fmt"

/*
编写一个函数来查找字符串数组中的最长公共前缀。
如果不存在公共前缀，返回空字符串 ""。

示例 1：
输入：strs = ["flower","flow","flight"]
输出："fl"
示例 2：
输入：strs = ["dog","racecar","car"]
输出：""
解释：输入不存在公共前缀。

提示：
1 <= strs.length <= 200
0 <= strs[i].length <= 200
strs[i] 如果非空，则仅由小写英文字母组成
*/
var C4_str1 = []string{"flower", "flow", "flight"}

// var str1 = []string{"adddogsss", "addracecar", "addcar"}

func MaxPrefix(strInfo []string) {
	if len(strInfo) < 2 {
		fmt.Println("输入字符串数组长度小于2,无效数据.")
		return
	}

	maxL := len(strInfo[0])
	fmt.Println("strInfo[0]", strInfo[0], "长度为:", maxL)
strLoop:
	for maxL > 0 {
		maxL--
		for i := 1; i < len(strInfo); i++ {
			if len(strInfo[i]) > maxL {
				if strInfo[i-1][maxL] != strInfo[i][maxL] { //当前和前一个字符串的maxL位比较
					continue strLoop
				}
			} else {
				continue strLoop
			}
		}
		break
	}
	fmt.Println("最大公共前缀为：", strInfo[0][:maxL+1])
}
