package gs

import (
	"sync"
	"time"

	"github.com/mereithhh/gs/logger"
)

type Task struct {
	// 工作者映射
	workers    map[string]Worker
	chanLength int
	wg         sync.WaitGroup
	name       string
	// 记录日志的 writer
	logger *logger.Logger
}
type PrintLogFunc func(format string, v ...interface{})

type Worker interface {
	// 对于一个 worker，他必须有这些方法。
	// 开始运行
	Start()
	// 获取名称
	GetName() string
	// 获取输出通道
	GetOutputChan() chan interface{}
	// 工作者类型
	GetWorkerType() string
	// 输出日志
	PrintLog(format string, v ...interface{})
}

func NewTask(chanLength int, taskName string) *Task {
	return &Task{
		workers:    make(map[string]Worker),
		chanLength: chanLength,
		wg:         sync.WaitGroup{},
		name:       taskName,
		logger:     logger.New(),
	}
}

func (t *Task) addWorker(worker Worker) {
	// 创建工作者
	// 将工作者添加到工作者映射中
	t.workers[worker.GetName()] = worker
	t.wg.Add(1)
	// 设置日志输出函数
	t.logger.RegisterWorkerLine(t.name, worker.GetWorkerType(), worker.GetName())
}

func (t *Task) UseInputer(name string, inputFunc InputFunc) {
	// 创建一个 inputer
	inputer := t.NewInputer(name, inputFunc, t.chanLength, &t.wg, t.name, t.logger)
	// 将 inputer 添加到工作者映射中
	t.addWorker(inputer)
}

func (t *Task) UseOutputer(name string, outputFunc OutputFunc, inputWorkerName string) {
	// 创建一个 outputer
	outputer := t.NewOutputer(name, outputFunc, inputWorkerName, &t.wg, t.name, t.logger)
	// 将 outputer 添加到工作者映射中
	t.addWorker(outputer)
}

func (t *Task) UseWasher(name string, washFunc WashFunc, inputWorkerNames []string) {
	// 创建一个 washer
	washer := t.NewWasher(name, washFunc, inputWorkerNames, &t.wg, t.name, t.logger)
	// 将 washer 添加到工作者映射中
	t.addWorker(washer)
}

func (t *Task) getWorker(name string) Worker {
	return t.workers[name]
}

func (t *Task) CloseWorker(name string) {
	delete(t.workers, name)
}

func (t *Task) GetName() string {
	return t.name
}

func (t *Task) Run() {
	go t.logger.Run()
	t.logger.RegisterTaskLine(t.name)
	t.logger.SetTaskLine(t.name, "开始任务！")

	start := time.Now()
	// 多线程开始所有工作者
	for _, worker := range t.workers {
		t.logger.SetWorkerLine(t.name, worker.GetWorkerType(), worker.GetName(), "准备执行！")
		go worker.Start()
	}

	t.wg.Wait()
	t.logger.SetTaskLine(t.name, "任务执行完成！用时：%v", time.Since(start))
	t.logger.Close()

}
