package rpcx

import (
	"github.com/golang/protobuf/proto"
	. "github.com/name5566/leaf/msg"
	"reflect"
)


type RpcType uint8

type RpcCall struct {
	Msg 		*Messagex
	Done 		chan *RpcCall
}

func NewRpcCall(method string,id uint32,data proto.Message) *RpcCall {
	c := new(RpcCall)
	c.Msg = new(Messagex)
	c.Msg.Request = NewRpcMessage(method,id,data)
	c.Done = make(chan *RpcCall,1)
	return c
}

type RpcCallx struct {
	*RpcCall
	Method 		*reflect.Method
}

func NewRpcCallx(msg *Messagex,f *reflect.Method) *RpcCallx {
	cx := new(RpcCallx)
	cx.RpcCall = new(RpcCall)
	cx.RpcCall.Msg = msg
	cx.RpcCall.Done = make(chan *RpcCall,1)
	cx.Method = f
	return cx
}
