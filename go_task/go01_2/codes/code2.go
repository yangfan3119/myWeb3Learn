package codes

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
Goroutine:
题目1 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
考察点 ： go 关键字的使用、协程的并发执行。
题目2 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
考察点 ：协程原理、并发任务调度。
*/
// 题目2
type Task func(tDelay int) error

func setTask() Task {
	return func(tDelay int) error {
		time.Sleep(time.Second * time.Duration(tDelay))
		if tDelay == 3 {
			return fmt.Errorf("error task 3")
		}
		return nil
	}
}

// TaskResult 存储任务执行结果
type TaskResult struct {
	TaskID    int
	Err       error
	DlayTime  int
	Duration  time.Duration
	StartTime time.Time
	EndTime   time.Time
}
type Scheduler struct {
	tasks []Task
}

func (s *Scheduler) AddTask(task Task) {
	s.tasks = append(s.tasks, task)
}

func (s *Scheduler) Run() []TaskResult {
	var wg sync.WaitGroup
	l_sTasks := len(s.tasks)
	results := make([]TaskResult, l_sTasks)
	resultChan := make(chan TaskResult, l_sTasks)
	wg.Add(l_sTasks)
	for i, task := range s.tasks {
		go func(id int, t Task) {
			defer wg.Done()

			start := time.Now()

			dT := (rand.Intn(50) + 10) / 10
			err := t(dT)

			end := time.Now()
			duration := end.Sub(start)

			resultChan <- TaskResult{
				TaskID:    id,
				DlayTime:  dT,
				Duration:  duration,
				Err:       err,
				StartTime: start,
				EndTime:   end,
			}
		}(i, task)
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for res := range resultChan {
		results[res.TaskID] = res
	}

	return results
}

func C2_MainRun() {
	scheduler := &Scheduler{tasks: []Task{}}

	for range 5 {
		scheduler.AddTask(setTask())
	}

	start := time.Now()
	defer func() {
		fmt.Println("程序总耗时：", time.Since(start))
	}()

	errStart := func(e error) string {
		if e == nil {
			return "执行成功"
		}
		return fmt.Sprintf("%v", e)
	}
	resDatas := scheduler.Run()
	for _, res := range resDatas {
		fmt.Printf("taskId: %d, Err: %s, DlayTime: %d, Duration: %v, StartTime: %s, EndTime: %s \n",
			res.TaskID, errStart(res.Err), res.DlayTime, res.Duration,
			res.StartTime.Format("15:04:05.000"), res.EndTime.Format("15:04:05.000"))
	}
}

// 题目1
func printOddNum() {
	for i := 1; i <= 10; i += 2 {
		fmt.Println("printOddNum: ", i)
	}
}

func printEvenNum() {
	for i := 2; i <= 10; i += 2 {
		fmt.Println("printEvenNum: ", i)
	}
}

func C2_PrintGo() {
	go printOddNum()
	go printEvenNum()

	time.Sleep(time.Second * 2)
}

func C2_PrintGo2() {
	var wg sync.WaitGroup
	oddDone := make(chan bool, 1)
	evenDone := make(chan bool, 1)

	wg.Add(2)

	// 奇数协程
	go func() {
		defer wg.Done()
		for i := 1; i <= 9; i += 2 { // 修改：循环到9为止
			<-evenDone
			fmt.Println(i)
			oddDone <- true
		}
		// 最后一次不需要发送信号，直接退出
	}()

	// 偶数协程
	go func() {
		defer wg.Done()
		evenDone <- true // 启动第一个循环
		for i := 2; i <= 10; i += 2 {
			<-oddDone
			fmt.Println(i)
			if i < 10 { // 修改：最后一次不发送信号
				evenDone <- true
			}
		}
		close(evenDone)
	}()

	wg.Wait()
}
