package rpc

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
	"strings"
	"sync"
)

var rpcHandlerMap = make(map[string]*RpcHandler)
var rpcHandlerMapLocker sync.RWMutex

var rpcRouter = &RpcRouter{Items: make([]*RpcRouterItem,0,1000)}
var rpcRouterLocker sync.RWMutex
var Clients = []*Agent{}

type CallInfo struct {
	cid 		uint32
	requestMsg 	proto.Message
	replyMsg	proto.Message
	err 		error
	done 		chan *CallInfo
}

func GetRpcHandler(name string) (*RpcHandler,*RpcMethod) {
	methods := strings.Split(name,".")
	if len(methods) != 2 {
		return nil,nil
	}

	rpcHandlerMapLocker.RLock()
	defer rpcHandlerMapLocker.RUnlock()
	if r,ok := rpcHandlerMap[methods[0]];ok {
		if m,ok := r.rpcMethodMap[methods[1]];ok {
			return r,m
		}
	}
	return nil,nil
}

func CallLocal(name string,requestMsg interface{}) (interface{},error) {
	rpcHandler,rpcMethod := GetRpcHandler(name)
	if rpcMethod == nil {
		return nil,fmt.Errorf("rpc call local [%v] method not exist",name)
	}

	if rpcMethod.RequestType != reflect.TypeOf(requestMsg) {
		return nil,fmt.Errorf("rpc call local [%v] request type error",name)
	}

	if rpcMethod.ReplyType == nil {
		return nil,fmt.Errorf("rpc call local [%v] have not reply",name)
	}

	replyMsg := reflect.New(rpcMethod.ReplyType.Elem())
	rpcCall := &RpcCall{
		method: rpcMethod.Method,
		requestMsg: reflect.ValueOf(requestMsg),
		replyMsg: replyMsg,
		err: nil,
		done: make(chan *RpcCall,1),
	}
	rpcHandler.GetRpcCallChan()<-rpcCall
	<-rpcCall.done

	return rpcCall.replyMsg,rpcCall.err
}


func CastLocal(name string,requestMsg interface{}) (error) {
	rpcHandler,rpcMethod := GetRpcHandler(name)
	if rpcMethod == nil {
		return fmt.Errorf("rpc cast local [%v] method not exist",name)
	}

	if rpcMethod.RequestType != reflect.TypeOf(requestMsg) {
		return fmt.Errorf("rpc cast local [%v] request type error",name)
	}

	if rpcMethod.ReplyType != nil {
		return fmt.Errorf("rpc cast local [%v] have reply",name)
	}

	replyMsg := reflect.New(rpcMethod.ReplyType.Elem())
	rpcCall := &RpcCall{
		method: rpcMethod.Method,
		requestMsg: reflect.ValueOf(requestMsg),
		replyMsg: replyMsg,
		err: nil,
		done: make(chan *RpcCall,1),
	}
	rpcHandler.GetRpcCallChan()<-rpcCall
	<-rpcCall.done

	return rpcCall.err
}
