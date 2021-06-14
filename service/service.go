package service

import (
	"github.com/name5566/leaf/rpcx"
	"github.com/name5566/leaf/timer"
	"reflect"
	"time"
)

const (
	TimerDispatcherLen = 100
)

type Service struct {
	serviceName 	string
	OnStart 		func()
	rpcChannel 		*rpcx.RpcChannel
	dispatcher      *timer.Dispatcher
}

func NewService() *Service {
	s := new(Service)
	s.rpcChannel = new(rpcx.RpcChannel)
	s.dispatcher = timer.NewDispatcher(TimerDispatcherLen)
	return s
}

func (s *Service) GetServiceType() string {
	return reflect.TypeOf(s).Name()
}

func (s *Service) GetServiceName() string {
	return s.serviceName
}

func (s *Service) SetServiceName(serviceName string) {
	s.serviceName = serviceName
}

func (s *Service) Init(serv interface{}) {
	//s.rpcHandler.Init(Service,100)
	s.rpcChannel.Init(serv)
}

func (s *Service) Destroy() {
	s.rpcChannel.OnClose()
}

func (this *Service) Run(closeSig chan bool) {
	if this.OnStart != nil {
		this.OnStart()
	}
	for {
		select {
		case <-closeSig:
			return
		case c := <- this.rpcChannel.GetChannel():
			this.rpcChannel.Cb(c)
		//case c := <- this.rpcHandler.GetRpcCallChan():
		//	this.rpcHandler.Exec(c)
		//case c := <- this.rpcHandler.GetRpcCastChan():
		//	this.rpcHandler.Exec(c)
		case t := <-this.dispatcher.ChanTimer:
			t.Cb()
		}
	}
}

func (s *Service) GetRpcChannel() *rpcx.RpcChannel {
	return s.rpcChannel
}


func (s *Service) AfterFunc(d time.Duration, cb func(),ct...int32) *timer.Timer {
	return s.dispatcher.AfterFunc(d, cb, ct...)
}

func (s *Service) CronFunc(cronExpr *timer.CronExpr, cb func()) *timer.Cron {
	return s.dispatcher.CronFunc(cronExpr, cb)
}
