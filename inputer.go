package gs

import (
	"fmt"
	"sync"

	"github.com/mereithhh/gs/logger"
)

// inputer 只是单纯的提取数据到指定的地方。不同的通道之间通过带缓冲的通道相连，通道的数量可以配置
// 对于 intper，只有一个输入函数就行了，因为它只有一个输入通道
type Inputer struct {
	// 名称
	name string
	// 输入函数
	inputFunc InputFunc
	// 输出通道
	outputChan chan interface{}
	// 等待池子
	wg *sync.WaitGroup
}

type InputFunc func(intputChan chan interface{}, printLogFunc PrintLogFunc)

func (t *Task) NewInputer(name string, inputFunc InputFunc, chanLength int, wg *sync.WaitGroup) *Inputer {
	return &Inputer{
		name:       name,
		inputFunc:  inputFunc,
		outputChan: make(chan interface{}, chanLength),
		wg:         wg,
	}
}

func (i *Inputer) GetName() string {
	return i.name
}

func (i *Inputer) Start() {
	i.inputFunc(i.outputChan, i.PrintLog)
	i.PrintLog("完成！")
	close(i.outputChan)
	i.wg.Done()
}

func (i *Inputer) GetOutputChan() chan interface{} {
	return i.outputChan
}
func (i *Inputer) GetWorkerType() string {
	return "Inputer"
}

func (i *Inputer) PrintLog(format string, v ...interface{}) {
	prefix := fmt.Sprintf("[%s][%s]:", i.GetWorkerType(), i.GetName())
	f := prefix + format
	logger.Printf(f, v...)
}
