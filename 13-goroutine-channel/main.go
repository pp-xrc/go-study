/*
第 13 课：Goroutine 与 Channel

- go f() 启动一个新协程，和 main 并发跑
- channel 用来在协程之间传数据
- 无缓冲：发送和接收必须同时准备好
- 有缓冲：相当于带容量的队列
- Worker Pool：固定几个工人 + 任务 channel，避免 goroutine 开太多

运行：
  go run ./13-goroutine-channel
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== 第 13 课：Goroutine + Channel ===")
	demoGoroutine()
	demoUnbuffered()
	demoBuffered()
	demoWorkerPool()
}

func demoGoroutine() {
	fmt.Println("\n--- go 关键字 ---")
	var wg sync.WaitGroup // 等待组，用来等待协程结束
	wg.Add(2) // 设置等待的协程数
	go func() {
		defer wg.Done() // 协程结束时调用 Done
		for i := 1; i <= 3; i++ {
			fmt.Println("  A", i)
			time.Sleep(30 * time.Millisecond) // 模拟协程工作
		}
	}()
	go func() {
		defer wg.Done()
		for i := 1; i <= 3; i++ {
			fmt.Println("  B", i)
			time.Sleep(30 * time.Millisecond)
		}
	}()
	wg.Wait() // 等两个协程都结束；否则 main 一退出程序就没了
}

func demoUnbuffered() {
	fmt.Println("\n--- 无缓冲 channel：一边发、一边收 ---")
	ch := make(chan string)

	go func() { // 启动一个协程，往 channel 发数据
		ch <- "hello" // 没人收就会卡住
		close(ch)     // 发完由发送方关闭
	}()

	msg := <-ch
	fmt.Println("收到:", msg)
}

func demoBuffered() {
	fmt.Println("\n--- 有缓冲 channel：容量 2，没满就能发 ---")
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	for v := range ch { // range 直到 close
		fmt.Println("  读到", v)
	}
}

func demoWorkerPool() {
	fmt.Println("\n--- Worker Pool：2 个工人消化 5 个任务 ---")
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	var wg sync.WaitGroup
	for w := 1; w <= 2; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range jobs {
				fmt.Println("  worker", id, "做任务", j)
				results <- j * 2
			}
		}(w)
	}

	for i := 1; i <= 5; i++ {
		jobs <- i
	}
	close(jobs) // 不再发任务，工人的 range 才会结束
	wg.Wait()
	close(results)
	for r := range results {
		fmt.Println("  结果", r)
	}
}
