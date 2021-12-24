package gs

import (
	"sync"

	"github.com/mereithhh/gs/logger"
)

// washer 负责清洗数据，并将清洗后的数据放入到输出通道中
type Washer struct {
	name string
	// 输入通道
	inputChans map[string]chan interface{}
	// 输出通道
	outputChan chan interface{}
	// 清洗函数
	washFunc WashFunc
	// 等待池子
	wg *sync.WaitGroup
	// 任务名称
	taskName string
	// logger
	logger *logger.Logger
}
type WashFunc func(inputChans map[string]chan interface{}, outputChan chan interface{}, printLogFunc PrintLogFunc)

func (t *Task) NewWasher(name string, washFunc WashFunc, inputWorkerNames []string, wg *sync.WaitGroup, taskName string, logger *logger.Logger) *Washer {
	// 找 input worker 的 name
	// var inputChans map[string]chan interface{}
	inputChans := make(map[string]chan interface{})
	for _, inputWorkerName := range inputWorkerNames {
		// fmt.Println("inputWorkerName:", inputWorkerName)
		inputWorker := t.getWorker(inputWorkerName)
		// fmt.Println("inputWorker:", inputWorker)
		inputChans[inputWorkerName] = inputWorker.GetOutputChan()
	}

	return &Washer{
		name:       name,
		inputChans: inputChans,
		outputChan: make(chan interface{}, t.chanLength),
		washFunc:   washFunc,
		wg:         wg,
		taskName:   taskName,
		logger:     logger,
	}
}

func (w *Washer) GetName() string {
	return w.name
}

func (w *Washer) Start() {
	w.washFunc(w.inputChans, w.outputChan, w.PrintLog)
	// w.PrintLog("完成！")
	close(w.outputChan)
	w.wg.Done()
}

func (w *Washer) GetOutputChan() chan interface{} {
	return w.outputChan
}

func (w *Washer) GetWorkerType() string {
	return "Washer"
}
func (i *Washer) PrintLog(format string, v ...interface{}) {
	i.logger.SetWorkerLine(i.taskName, i.GetWorkerType(), i.GetName(), format, v...)
}
