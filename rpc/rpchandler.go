package rpc

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/common"
	. "github.com/name5566/leaf/pb"
	"github.com/name5566/leaf/log"
	"reflect"
	"strings"
)

var err error
var ErrorType = reflect.TypeOf(err)

type IRpcHandler interface {
	Init()
	GetModel() interface{}
	GetRpcCallChan() chan *RpcCall
	GetRpcCastChan() chan *RpcCast
}

type RpcCall struct {
	method 		reflect.Method
	requestMsg	reflect.Value
	replyMsg 	reflect.Value
	err 		error
	done 		chan *RpcCall
}

type RpcCast struct {
	method 		reflect.Method
	requestMsg	reflect.Value
}

type RpcMethod struct {
	Method    		reflect.Method
	RequestId 		uint16
	ReplyId 		uint16
	RequestType		reflect.Type
	ReplyType 		reflect.Type
}

type RpcHandler struct {
	model				interface{}
	rpcMethodMap 		map[string]*RpcMethod
	rpcCallChan 		chan *RpcCall
	rpcCastChan			chan *RpcCast
}

func (r *RpcHandler) GetRpcMethod(name string) *RpcMethod {
	if m,ok := r.rpcMethodMap[name];ok {
		return m
	}
	return nil
}


func (r *RpcHandler) IsRpcHandler(m reflect.Method) bool {
	methodType := m.Type
	if !strings.HasPrefix(m.Name,"RPC_") {
		return false
	}

	if methodType.NumIn() != 2 {
		return false
	}

	if methodType.In(0).Kind() != reflect.Ptr || !common.CheckMessage(methodType.In(1))  {
		return false
	}

	if methodType.NumOut() == 1 && methodType.Out(0).Kind() == reflect.Interface {
		return true
	}

	if methodType.NumOut() == 2 && common.CheckMessage(methodType.Out(0)) && methodType.Out(1).Kind() == reflect.Interface {
		return true
	}

	return false
}

func (r *RpcHandler) register(model interface{}) {
	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)
	modelName := reflect.Indirect(modelValue).Type().Name()
	if modelType.Kind() != reflect.Ptr {
		log.Error("RpcHandler register error: model-[%v] not ptr",modelName)
		return
	}

	for i:=0;i<modelType.NumMethod();i++{
		method := modelType.Method(i)
		if r.IsRpcHandler(method) {
			rpcMethod := new(RpcMethod)
			rpcMethod.Method = method
			rpcMethod.RequestType = method.Type.In(1)
			rpcMethod.RequestId = common.GetMessageId(rpcMethod.RequestType)
			if method.Type.NumOut() == 1 {
				rpcMethod.ReplyType = nil
				rpcMethod.ReplyId = 0
			}else{
				rpcMethod.ReplyType = method.Type.Out(0)
				rpcMethod.ReplyId = common.GetMessageId(rpcMethod.ReplyType)
			}
			methodName := strings.TrimPrefix(method.Name,"RPC_")
			r.rpcMethodMap[methodName] = rpcMethod

			rpcRouterLocker.Lock()
			rpcRouter.Items = append(rpcRouter.Items,&RpcRouterItem{
				Method: proto.String(modelName + "."+ methodName),
				ReplyMsgId: proto.Uint32(uint32(rpcMethod.ReplyId)),
				ReqMsgId: proto.Uint32(uint32(rpcMethod.RequestId)),
			})
			rpcRouterLocker.Unlock()
		}
	}

	rpcHandlerMapLocker.Lock()
	rpcHandlerMap[modelName] = r
	rpcHandlerMapLocker.Unlock()
}


func (r *RpcHandler) Init(model interface{},rpcChanLen int) {
	r.rpcMethodMap = make(map[string]*RpcMethod)
	r.rpcCallChan =  make(chan *RpcCall,rpcChanLen)
	r.rpcCastChan =  make(chan *RpcCast,rpcChanLen)
	r.model = model
	r.register(model)
}

func (r *RpcHandler) OnClose() {

}

func (r *RpcHandler) GetRpcCallChan() chan *RpcCall {
	return r.rpcCallChan
}

func (r *RpcHandler) GetRpcCastChan() chan *RpcCast {
	return r.rpcCastChan
}

func (r *RpcHandler) Exec(data interface{}) {
	if c,ok := data.(*RpcCall);ok {
		ret := c.method.Func.Call([]reflect.Value{reflect.ValueOf(r.model),c.requestMsg})
		if len(ret) != 2 {
			c.err = errors.New("RpcHandler Exec error:ret length error")
		}

		c.replyMsg = ret[0]
		if err,ok := ret[1].Interface().(error);ok {
			c.err = err
		}

		c.done<-c
	}

	if c,ok := data.(*RpcCast);ok {
		ret := c.method.Func.Call([]reflect.Value{reflect.ValueOf(r.model),c.requestMsg})
		if len(ret) != 1 {
			log.Error("RpcHandler Exec error:ret length error")
		}
		if err,ok := ret[0].Interface().(error);ok {
			log.Error("RpcHandler Exec error:%v",err.Error())
		}
	}
}

func (r *RpcHandler) Call(c *RpcCall) {
	ret := c.method.Func.Call([]reflect.Value{reflect.ValueOf(r.model),c.requestMsg})
	if len(ret) != 2 {
		c.err = errors.New("RpcHandler Exec error:ret length error")
	}

	c.replyMsg = ret[0]
	if err,ok := ret[1].Interface().(error);ok {
		c.err = err
	}

	c.done<-c
}

func (r *RpcHandler) Cast(c *RpcCast) {
	ret := c.method.Func.Call([]reflect.Value{reflect.ValueOf(r.model),c.requestMsg})
	if len(ret) != 1 {
		log.Error("RpcHandler Exec error:ret length error")
	}
	if err,ok := ret[0].Interface().(error);ok {
		log.Error("RpcHandler Exec error:%v",err.Error())
	}
}


