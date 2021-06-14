package rpcx

import (
	"fmt"
	. "github.com/name5566/leaf/common"
	"github.com/name5566/leaf/log"
	. "github.com/name5566/leaf/msg"
	"reflect"
	"runtime"
)


type IRpcChannel interface {
	Init()
	GetService() interface{}
	GetChannel() chan *RpcCallx
	Exec()
}

type RpcChannel struct {
	service				interface{}
	rpcCallChannel 		chan *RpcCallx
}

func (r *RpcChannel) Init(s interface{}) {
	r.service = s
	r.rpcCallChannel = make(chan *RpcCallx,100)
}

func (r *RpcChannel) OnClose() {

}

func (r *RpcChannel) GetChannel() chan *RpcCallx {
	return r.rpcCallChannel
}

func (r *RpcChannel) Exec(method *reflect.Method,msg *Messagex) {
	cx := NewRpcCallx(msg,method)
	r.rpcCallChannel<-cx
	if msg.Response != nil {
		<- cx.Done
	}
}

func (r *RpcChannel) Cb(cx *RpcCallx) {
	var (
		method = cx.Method
		msgx = cx.Msg
	)

	defer func() {
		if r := recover(); r != nil {
			if msgx.Response != nil {
				msgx.Response.SetError(Error_RpcExecErr.Code())
			}

			buf := make([]byte, 4096)
			l := runtime.Stack(buf, true)
			err := fmt.Errorf("%v: %s", r, buf[:l])
			log.Error("error stack:%v \n",err)
		}
	}()

	res := method.Func.Call([]reflect.Value{reflect.ValueOf(r.service),reflect.ValueOf(msgx.Request.GetMessage())})
	if len(res) == 2 {
		response := res[0].Interface()
		msgx.Response.SetMessage(response)

		err := res[1].Interface()
		if IsNil(err) {
			if IsNil(response) {
				msgx.Response.SetError(Error_RpcRespNil.Code())
			}
		}else{
			if _,ok := err.(*Error);ok {
				msgx.Response.SetError(err.(*Error).Code())
			}else{
				msgx.Response.SetError(Error_RpcUnknowErr.Code())
			}
		}

		cx.Done<-cx.RpcCall
	}
}
