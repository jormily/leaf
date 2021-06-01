package service

import (
	"github.com/name5566/leaf/rpc"
	"github.com/name5566/leaf/timer"
	"time"
)

const (
	TimerDispatcherLen = 100
)

type Service struct {
	OnRun 			func()
	rpcHandler 		*rpc.RpcHandler
	dispatcher      *timer.Dispatcher
}

func NewService() *Service {
	s := new(Service)
	s.rpcHandler = new(rpc.RpcHandler)
	s.dispatcher = timer.NewDispatcher(TimerDispatcherLen)
	return s
}

func (s *Service) Init(Service interface{}) {
	s.rpcHandler.Init(Service,100)

}

func (s *Service) Destroy() {
	s.rpcHandler.OnClose()
}

func (this *Service) Run(closeSig chan bool) {
	if this.OnRun != nil {
		this.OnRun()
	}
	for {
		select {
		case <-closeSig:
			return
		case c := <- this.rpcHandler.GetRpcCallChan():
			this.rpcHandler.Exec(c)
		case c := <- this.rpcHandler.GetRpcCastChan():
			this.rpcHandler.Exec(c)
		case t := <-this.dispatcher.ChanTimer:
			t.Cb()
		}
	}
}


func (s *Service) AfterFunc(d time.Duration, cb func(),ct...int32) *timer.Timer {
	return s.dispatcher.AfterFunc(d, cb, ct...)
}

func (s *Service) CronFunc(cronExpr *timer.CronExpr, cb func()) *timer.Cron {
	return s.dispatcher.CronFunc(cronExpr, cb)
}
