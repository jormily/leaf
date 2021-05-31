package module

import (
	"github.com/name5566/leaf/rpc"
	"github.com/name5566/leaf/timer"
	"time"
)

type Module struct {
	TimerDispatcherLen int
	AsynCallLen        int

	rpcHandler 		*rpc.RpcHandler
	dispatcher      *timer.Dispatcher
}

func NewModule() *Module {
	m := new(Module)
	m.rpcHandler = new(rpc.RpcHandler)
	m.dispatcher = timer.NewDispatcher(m.TimerDispatcherLen)
	return m
}

func (m *Module) Init(Module interface{},rpcChanLen int) {
	m.rpcHandler.Init(Module,rpcChanLen)
}

func (m *Module) OnClose() {
	m.rpcHandler.OnClose()
}

func (this *Module) Run(closeSig chan bool) {
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


func (m *Module) AfterFunc(d time.Duration, cb func()) *timer.Timer {
	if m.TimerDispatcherLen == 0 {
		panic("invalid TimerDispatcherLen")
	}

	return m.dispatcher.AfterFunc(d, cb)
}

func (m *Module) CronFunc(cronExpr *timer.CronExpr, cb func()) *timer.Cron {
	if m.TimerDispatcherLen == 0 {
		panic("invalid TimerDispatcherLen")
	}

	return m.dispatcher.CronFunc(cronExpr, cb)
}
