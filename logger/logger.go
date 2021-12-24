package logger

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gosuri/uilive"
)

type Logger struct {
	writer *uilive.Writer
	// lines 保存每一个 key 对应的数据
	lines map[string]string
	// keys 保存每一行对应的 key
	keys map[int]string
	// 锁
	lock *sync.Mutex
	// 刷新通道，来一个就刷新，不然就不刷新. 通道关了就 stop
	Flush chan bool
	// 开始时间
	startTime time.Time
}

func New() *Logger {
	return &Logger{
		writer:    uilive.New(),
		lines:     make(map[string]string),
		keys:      make(map[int]string),
		lock:      &sync.Mutex{},
		Flush:     make(chan bool),
		startTime: time.Now(),
	}
}

func (l *Logger) RegisterLine(key string) {
	l.lock.Lock()
	l.keys[len(l.keys)] = key
	l.lock.Unlock()
}

func (l *Logger) RegisterTaskLine(taskName string) {
	l.RegisterLine(l.GetTaskKey(taskName))
}
func (l *Logger) RegisterWorkerLine(taskName string, workerType string, workerName string) {
	l.RegisterLine(l.GetWorkerKey(taskName, workerType, workerName))
}

func (l *Logger) Run() {
	// fmt.Println("我运行了")
	l.writer.Start()
	for {
		_, ok := <-l.Flush
		// fmt.Println("输出一个 flush")
		if !ok {
			l.writer.Stop()
			break
		}
		l.Print()
	}
}

// 有一个 setLine 方法，设置每一行的内容
// 有一个输出的方法，输出内容

func (l *Logger) Print() {
	// 第一行单独输出
	l.lock.Lock()
	length := len(l.keys)
	toPrint := ""
	for i := length - 1; i > -1; i-- {
		key := l.keys[i]
		val := l.lines[key]
		toPrint = toPrint + val + "\n"
	}
	// lastLine := fmt.Sprintf("当前时间：%s，已耗时：%s 。", time.Now().Format("2006-01-02 15:04:05"), time.Since(l.startTime))
	// toPrint = toPrint + lastLine
	// fmt.Println(toPrint)
	fmt.Fprintln(l.writer, toPrint)
	l.writer.Flush()
	l.lock.Unlock()
}

func (l *Logger) Close() {
	l.lock.Lock()
	length := len(l.keys)
	toPrint := ""
	for i := 0; i < length; i++ {
		key := l.keys[i]
		val := l.lines[key]
		toPrint = toPrint + val + "\n"
	}
	fmt.Fprintln(l.writer.Bypass(), toPrint)
	l.lock.Unlock()
	l.writer.Stop()
	close(l.Flush)
}

func (l *Logger) SetTaskLine(taskName string, format string, v ...interface{}) {
	key := l.GetTaskKey(taskName)
	prefix := fmt.Sprintf("[%v][%s]", time.Now().Format("2006-01-02 15:04:05"), taskName)
	str := fmt.Sprintf(format, v...)
	l.SetLine(key, prefix+str)
}

func (l *Logger) SetWorkerLine(taskName string, workerType string, workerName string, format string, v ...interface{}) {
	key := l.GetWorkerKey(taskName, workerType, workerName)
	prefix := fmt.Sprintf("[%v][%s][%s][%s]", time.Now().Format("2006-01-02 15:04:05"), taskName, workerType, workerName)
	str := fmt.Sprintf(format, v...)
	l.SetLine(key, prefix+str)
}

func (l *Logger) SetLine(key string, val string) {
	l.lock.Lock()
	l.lines[key] = val
	l.lock.Unlock()
	l.Flush <- true
	// fmt.Println("输入一个 flush")

}

func (l *Logger) GetTaskKey(taskName string) string {
	return taskName
}

func (l *Logger) GetWorkerKey(taskName string, workerType string, workerName string) string {
	key := strings.Join([]string{taskName, workerType, workerName}, "/")
	return key
}
