package codes

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/*
锁机制
题目1 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ： sync.Mutex 的使用、并发数据安全。
题目2 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ：原子操作、并发数据安全。
*/
// 题目2
func GoAtomicCounter(nLoop int) uint64 {
	var wg sync.WaitGroup
	var mCounter uint64

	start := time.Now()
	defer func() {
		fmt.Printf("GoAtomicCounter elapsed: %v\n", time.Since(start))
	}()

	counter := func() {
		defer wg.Done()
		for range 10000 {
			atomic.AddUint64(&mCounter, 1)
		}
	}

	wg.Add(nLoop)
	for range nLoop {
		go counter()
	}

	wg.Wait()
	return mCounter
}

// 题目1
func GoMutexCounter(nLoop int) uint64 {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var mCounter uint64
	start := time.Now()
	defer func() {
		fmt.Printf("GoMutexCounter elapsed: %v\n", time.Since(start))
	}()

	counter := func() {
		defer wg.Done()
		for range 10000 {
			mu.Lock()
			mCounter++
			mu.Unlock()
		}
	}

	wg.Add(nLoop)
	for range nLoop {
		go counter()
	}

	wg.Wait()
	return mCounter
}
