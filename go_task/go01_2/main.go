package main

import (
	c "go01_2/codes"
)

func main() {
	c.Hello()

	// 2.1 go 关键字启动两个协程,分别输出奇数和偶数
	// c.C2_PrintGo2()
	// 2.2 任务调度
	c.C2_MainRun()

	// 4.1 channel 测试
	// t := c.ChanTransmitNum(5)
	// fmt.Println("Res:", t)
	// 4.2 缓冲通道测试
	// c.C4_ChannelPassNum(50, 2)
}
