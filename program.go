package husky

import (
	"os"
	"os/signal"
	"syscall"
)

type Signal chan os.Signal

func NewSignal() Signal {
	return make(chan os.Signal, 1)
}

func (s Signal) Finish() {
	close(s)
}

type Program interface {
	AddSignal(chs ...Signal)
	Run() os.Signal
}

type _Program struct {
	chs []Signal
}

func (s *_Program) AddSignal(chs ...Signal) {
	s.chs = append(s.chs, chs...)
}

func (s *_Program) Run() os.Signal {
	ch := NewSignal()
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
	return &_Program{chs: []Signal{}}
}
