package gs

import (
	"fmt"
	"sync"
)

func RunBulk(arg ...*Task) {
	var wg sync.WaitGroup
	for _, task := range arg {
		wg.Add(1)
		go func(task *Task) {
			task.Run()
			wg.Done()
		}(task)
	}
	wg.Wait()
	fmt.Println("所有任务全部完成！")
}
