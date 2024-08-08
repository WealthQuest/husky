package husky

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Signal chan struct{}

func NewSignal() Signal {
	return make(chan struct{}, 1)
}

func (s Signal) Finish() {
	close(s)
}

var programIns ProgramAPI

func InitProgram() {
	programIns = &_Program{
		Mutex:   sync.Mutex{},
		signals: make([]Signal, 0),
		stopCh:  make(chan struct{}),
	}
}

func Program() ProgramAPI {
	return programIns
}

type ProgramAPI interface {
	AddSignal(s Signal)
	Run()
	Stop()
}

type _Program struct {
	sync.Mutex
	signals []Signal
	stopCh  chan struct{}
}

func (p *_Program) AddSignal(s Signal) {
	p.Lock()
	defer p.Unlock()
	p.signals = append(p.signals, s)
}

func (s *_Program) Run() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	LogInfo("程序运行中")
	select {
	case sign := <-ch:
		LogWarn("收到系统停止信号:", sign)
	case <-s.stopCh:
		LogWarn("收到程序停止信号")
	}

	for index := range s.signals {
		ch := s.signals[index]
		ch <- struct{}{}
	}
	for index := range s.signals {
		ch := s.signals[index]
		<-ch
	}
	LogInfo("程序已停止")

}

func (s *_Program) Stop() {
	close(s.stopCh)
}
