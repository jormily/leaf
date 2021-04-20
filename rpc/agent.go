package rpc

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/common"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/network"
	"reflect"
	"sync"
)

type RouterItem struct {
	Method 		string
	ReqMsgId 	uint16
	ReplyMsgId 	uint16
	ReqType   	reflect.Type
	ReplyType 	reflect.Type
}

type Agent struct {
	conn 		network.Conn
	server 		*Server
	userData 	interface{}

	sync.Mutex

	rlock 		sync.Mutex
	cid  		uint32

	requests 	map[uint32]*CallInfo
	router 		map[string]*RouterItem
}


func newAgent(conn *network.TCPConn) network.Agent {
	agent := new(Agent)
	agent.conn = conn
	agent.requests = make(map[uint32]*CallInfo)
	agent.router = make(map[string]*RouterItem)
	return agent
}

func (a *Agent) getRid() uint32 {
	a.rlock.Lock()
	a.cid++
	if a.cid == 0 {
		a.cid++
	}
	
	cid := a.cid
	a.rlock.Unlock()

	return cid
}

func (a *Agent) makeRpcRequest(route string,msg proto.Message,call bool) *RpcRequest {
	req := new(RpcRequest)
	req.Method = proto.String(route)
	req.Data,_ = proto.Marshal(msg)
	if call {
		req.Rid = proto.Uint32(a.getRid())
	}
	return req
}


func (a *Agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		msg, msgtype, err := Unmarshal(data)
		if err != nil {
			log.Debug("unmarshal message error: %v", err)
			break
		}

		Route(msgtype, msg, a)
	}
}

func (a *Agent) WriteMsg(msg proto.Message) {
	data, err := Marshal(msg)
	if err != nil {
		log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
		return
	}
	err = a.conn.WriteMsg(data...)
	if err != nil {
		log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
	}
}

func (a *Agent) makeCallInfo(cid uint32,method string,reqMsg proto.Message,replyMsg proto.Message) *CallInfo {
	callInfo := new(CallInfo)
	callInfo.cid = cid
	callInfo.requestMsg = reqMsg
	callInfo.replyMsg = replyMsg
	callInfo.done = make(chan *CallInfo,1)

	if cid != 0 {
		if _,ok := a.requests[cid];ok {
			log.Error("MakeRpcCall cid = %v is exist",cid)
		}

		a.requests[cid] = callInfo
	}

	return callInfo
}

func (a *Agent) Call(method string,reqMsg proto.Message) (proto.Message,error) {
	//a.Lock()
	//defer a.Unlock()

	route,ok := a.router[method]
	if !ok {
		return nil, fmt.Errorf("rpc call %v error",method)
	}
	if route.ReqType != reflect.TypeOf(reqMsg) {
		return nil, fmt.Errorf("rpc call %v reqMsg error",method)
	}
	if route.ReqType == nil || route.ReplyType == nil {
		return nil, fmt.Errorf("rpc call ReqType = %v ReplyType = %v",route.ReqType, route.ReplyType)
	}
	replyMsg := reflect.New(route.ReplyType.Elem()).Interface()
	msg := a.makeRpcRequest(method,reqMsg,true)
	c := a.makeCallInfo(msg.GetRid(),method,reqMsg,replyMsg.(proto.Message))

	a.WriteMsg(msg)
	done := <- c.done

	return done.replyMsg,nil
}


func (a *Agent) Cast(method string,reqMsg proto.Message) error {
	//a.Lock()
	//defer a.Unlock()

	route,ok := a.router[method]
	if !ok {
		return fmt.Errorf("rpc cast %v error",method)
	}
	if route.ReqType != reflect.TypeOf(reqMsg) {
		return fmt.Errorf("rpc cast %v reqMsg error",method)
	}
	if route.ReqType == nil {
		return fmt.Errorf("rpc cast ReqType = %v", route.ReplyType)
	}

	msg := a.makeRpcRequest(method,reqMsg,false)
	a.makeCallInfo(msg.GetRid(),method,reqMsg,nil)

	a.WriteMsg(msg)
	
	return nil
}

func (a *Agent) OnClose() {

}

func (a *Agent) OnHeart(msg proto.Message) {

}

func (a *Agent) OnRequset(msg proto.Message) {
	rpcRequest := msg.(*RpcRequest)
	rpcResponse := &RpcResponse{
		Rid: rpcRequest.Rid,
	}

	rpcHandler,rpcMethod := GetRpcHandler(rpcRequest.GetMethod())
	if rpcMethod == nil {
		rpcResponse.Err = proto.String(fmt.Sprintf("rpc %v not exist",rpcRequest.GetMethod()))
		a.WriteMsg(rpcResponse)
		return
	}

	requestMsg := reflect.New(rpcMethod.RequestType.Elem())
	if err := proto.Unmarshal(rpcRequest.Data,requestMsg.Interface().(proto.Message));err != nil {
		rpcResponse.Err = proto.String(fmt.Sprintf("rpc %v data error",rpcRequest.GetMethod()))
		a.WriteMsg(rpcResponse)
		return
	}

	if rpcResponse.GetRid() == 0 {
		rpcHandler.GetRpcCastChan()<-&RpcCast{
			method: rpcMethod.Method,
			requestMsg: requestMsg,
		}
		return
	}else{
		if rpcMethod.ReplyType == nil {
			rpcResponse.Err = proto.String(fmt.Sprintf("rpc %v replytype is nil",rpcRequest.GetMethod()))
			a.WriteMsg(rpcResponse)
			return
		}

		replyMsg := reflect.New(rpcMethod.ReplyType.Elem())
		rpcCall := &RpcCall{
			method: rpcMethod.Method,
			requestMsg: requestMsg,
			replyMsg: replyMsg,
			err: nil,
			done: make(chan *RpcCall,1),
		}
		rpcHandler.GetRpcCallChan()<-rpcCall
		<-rpcCall.done

		if rpcCall.err != nil {
			rpcResponse.Err = proto.String(rpcCall.err.Error())
			a.WriteMsg(rpcResponse)
			return
		}

		rpcResponse.Data,_ = proto.Marshal(rpcCall.replyMsg.Interface().(proto.Message))
		a.WriteMsg(rpcResponse)

		return
	}
}

func (a *Agent) OnResponse(msg proto.Message) {
	rpcResponse := msg.(*RpcResponse)
	if rpcResponse.GetRid() == 0 {
		log.Error("rpc response rid zero")
		return
	}

	c,ok := a.requests[rpcResponse.GetRid()]
	if !ok {
		log.Error("rpc response rid %v not find",rpcResponse.GetRid())
		return
	}

	delete(a.requests,rpcResponse.GetRid())

 	err := proto.Unmarshal(rpcResponse.Data,c.replyMsg)
 	if err != nil {
 		c.err = err
 		c.done <- c
		return
	}

	if rpcResponse.Err != nil {
		c.err = errors.New(rpcResponse.GetErr())
	}
	c.done <- c

	return
}

func (a *Agent) OnRouter(msg proto.Message)  (proto.Message,error) {
	rpcRouter := msg.(*RpcRouter)
	for _,v := range rpcRouter.Items {
		item := &RouterItem{}
		item.Method = v.GetMethod()
		item.ReqMsgId = uint16(v.GetReqMsgId())
		item.ReplyMsgId = uint16(v.GetReplyMsgId())
		item.ReqType = common.GetMessageType(item.ReqMsgId)
		item.ReplyType = common.GetMessageType(item.ReplyMsgId)
		a.router[v.GetMethod()] = item
	}

	Clients = append(Clients,a)
	return nil,nil
}
