/*
第 14 课：WaitGroup 与互斥锁

- WaitGroup：等一组 goroutine 全部结束（Add → go → Done → Wait；Go 1.25+ 也可用 wg.Go）
- 多协程改同一变量必须加锁，否则结果是错的（竞态）
- for range n：只重复 n 次、不用下标时（见第 04 课）

运行：
  go run ./14-sync-lock
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("=== 第 14 课：同步锁 ===")
	demoWaitGroup()
	demoRace()
	demoMutex()
}

func demoWaitGroup() {
	fmt.Println("\n--- WaitGroup：Add → go → Done → Wait ---")
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1) // 必须在 go 之前 Add
		go func(id int) {
			defer wg.Done()
			fmt.Println("  worker", id, "完成")
		}(i) // 把 i 当参数传入，避免闭包踩同一个 i
	}
	wg.Wait()
}

func demoRace() {
	fmt.Println("\n--- 不加锁累加（结果经常不是 1000）---")
	var n int
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			n++ // 读改写不是原子的
		}()
	}
	wg.Wait()
	fmt.Println("  期望 1000，实际", n)
}

func demoMutex() {
	fmt.Println("\n--- 加 Mutex 后结果稳定是 1000 ---")
	var n int
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			n++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println("  实际", n)
}
