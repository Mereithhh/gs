package gs

import (
	"fmt"
	"sync"

	"github.com/mereithhh/gs/logger"
)

// outputer 负责输出数据，只有一个输入通道
type Outputer struct {
	name string
	// 输入通道
	inputChan chan interface{}
	// 输出函数
	outputFunc OutputFunc
	// 等待池子
	wg *sync.WaitGroup
}
type OutputFunc func(inputChan chan interface{}, printLogFunc PrintLogFunc)

func (t *Task) NewOutputer(name string, outputFunc OutputFunc, inputWorkerName string, wg *sync.WaitGroup) *Outputer {
	// 找 input worker 的 name
	inputWorker := t.getWorker(inputWorkerName)
	return &Outputer{
		name:       name,
		inputChan:  inputWorker.GetOutputChan(),
		outputFunc: outputFunc,
		wg:         wg,
	}
}

func (o *Outputer) GetName() string {
	return o.name
}
func (o *Outputer) Start() {
	o.outputFunc(o.inputChan, o.PrintLog)
	o.PrintLog("完成！")
	o.wg.Done()
}
func (o *Outputer) GetOutputChan() chan interface{} {
	return nil
}
func (o *Outputer) GetWorkerType() string {
	return "Outputer"
}
func (o *Outputer) PrintLog(format string, v ...interface{}) {
	prefix := fmt.Sprintf("[%s][%s]:", o.GetWorkerType(), o.GetName())
	f := prefix + format
	logger.Printf(f, v...)
}
