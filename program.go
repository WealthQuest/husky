package husky

import (
	"os"
	"os/signal"
	"syscall"
)

type Program interface {
	AddSignal(chs ...chan os.Signal)
	Run() os.Signal
}

type _Program struct {
	chs []chan os.Signal
}

func (s *_Program) AddSignal(chs ...chan os.Signal) {
	s.chs = append(s.chs, chs...)
}

func (s *_Program) Run() os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	LogInfo("program running")
	sign := <-ch
	LogWarn("program stop signal:", sign)

	for index := range s.chs {
		ch := s.chs[index]
		ch <- sign
	}
	for index := range s.chs {
		ch := s.chs[index]
		<-ch
	}
	LogInfo("program stoped")
	return sign
}

func NewProgram() Program {
	return &_Program{chs: []chan os.Signal{}}
}
