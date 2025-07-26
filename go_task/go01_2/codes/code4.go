package codes

import (
	"fmt"
	"sync"
)

/*
Channel
题目1 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，
另一个协程从通道中接收这些整数并打印出来。
考察点 ：通道的基本使用、协程间通信。
题目2 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
考察点 ：通道的缓冲机制。
*/
//题目2
func C4_ChannelPassNum(n int, chanSize int) {
	var wg sync.WaitGroup
	var pcChan = make(chan int, chanSize)

	producer := func(ch chan<- int) {
		defer wg.Done()
		for i := 1; i <= n; i++ {
			fmt.Printf("P%d\n", i)
			ch <- i
		}
		close(ch)
	}

	consumer := func(ch <-chan int) {
		defer wg.Done()
		for v := range ch {
			fmt.Printf("C%d, ch len=%d\n", v, len(ch))
		}
	}

	wg.Add(2)
	go producer(pcChan)
	go consumer(pcChan)

	wg.Wait()
}

// 题目1
func ChanTransmitNum(n int) string {
	var trans = make(chan int, 1)
	var resStr = make(chan string)
	var wg sync.WaitGroup

	go func(ch chan<- int) {
		for i := 1; i <= n; i++ {
			ch <- i
		}
		close(trans)
	}(trans)

	wg.Add(1)

	go func(ch <-chan int) {
		defer wg.Done()
		var resFmt string
		for v := range ch {
			resFmt = resFmt + fmt.Sprintf("%d,", v)
		}
		resStr <- resFmt
	}(trans)

	go func() {
		wg.Wait()
		close(resStr)
	}()

	res, ok := <-resStr
	if ok {
		return res
	}

	return "error"
}
