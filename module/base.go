package module

import (
	//"reflect"
	"github.com/name5566/leaf/rpc"
)


type IModel interface {
	Run()
}

type Model struct {
	*rpc.RpcHandler
}

func NewModel() *Model {
	m := new(Model)
	m.RpcHandler = new(rpc.RpcHandler)
	return m
}

func (m *Model) Init(model interface{},rpcChanLen int) {
	m.RpcHandler.Init(model,rpcChanLen)
}

func (m *Model) OnClose() {
	m.RpcHandler.OnClose()
}

func (this *Model) Run(closeSig chan bool) {
	for {
		select {
		case <-closeSig:
			return
		case c := <-this.GetRpcCallChan():
			this.Exec(c)
		case c := <- this.GetRpcCastChan():
			this.Exec(c)
		}
	}
}

