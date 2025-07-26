package codes

import "fmt"

/*
字符串:
有效的括号 ,
考察：字符串处理、栈的使用
题目：给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效
*/
var C3_testStr []string = []string{"[][{()}](){}", "{([])}{})"}

type ByteStack []byte

func (b *ByteStack) Push(a byte) {
	*b = append(*b, a)
	fmt.Println("push", string(a), "Stack = ", string(*b))
}

func (b *ByteStack) Pop() {
	if n := len(*b); n == 0 {
		return
	} else {
		ele := (*b)[n-1]
		*b = (*b)[:n-2]
		fmt.Println("pop byte ", string(ele), " then Stack = ", string(*b))
		return
	}
}

var jud = map[string]string{"}": "{", "]": "[", ")": "("}

func (b *ByteStack) ValidJudge() {
	if n := len(*b); n >= 2 {
		v, exist := jud[string((*b)[n-1])]
		if exist && v == string((*b)[n-2]) {
			(*b).Pop()
		}
	}
}

func ValidBracket(tests []string) {
	if tests == nil {
		tests = C3_testStr
	}

	for _, x := range tests {
		fmt.Println("Start verifying the string: ", x)
		b1 := []byte(x)
		var bStack ByteStack
		for _, b := range b1 {
			bStack.Push(b)
			bStack.ValidJudge()
		}
		if len(bStack) == 0 {
			fmt.Println(x, " is Valid Bracket!")
		} else {
			fmt.Println(x, " is Not Valid Bracket!")
		}
	}
}
