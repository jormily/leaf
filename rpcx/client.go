package rpcx

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/log"
	. "github.com/name5566/leaf/msg"
	"github.com/name5566/leaf/network"
	"github.com/name5566/leaf/pb"
	"reflect"
	"sync"
	"time"
)

var clients []*RpcClient
var clientsLock sync.Mutex

type RpcClient struct {
	sync.Mutex
	conn 		network.Conn

	client 		*network.TCPClient
	callId 		uint32
	calls 		map[uint32]*RpcCall
	handlers 	map[string]*pb.RpcHandler
}

func NewRpcClient(conn network.Conn,client *network.TCPClient) *RpcClient {
	rc := new(RpcClient)
	rc.conn = conn
	rc.client = client
	rc.calls = make(map[uint32]*RpcCall)
	rc.handlers = make(map[string]*pb.RpcHandler)
	return rc
}


func (rc *RpcClient) GetCallId() uint32 {
	rc.Lock()
	rc.callId++
	if rc.callId == 0 {
		rc.callId = 1
	}
	rc.Unlock()
	return rc.callId
}

func (rc *RpcClient) Run() {
	defer func() {
		clientsLock.Lock()
		for index,client := range clients {
			if client.conn == rc.conn {
				clients = append(clients[:index],clients[index+1:]...)
			}
		}
		clientsLock.Unlock()
	}()

	for {
		data, err := rc.conn.ReadMsg()
		if err != nil {
			log.Error("read message: %v", err)
			break
		}
		msg, err := Unmarshal(data)
		if err != nil {
			log.Error("unmarshal message error: %v", err)
			break
		}

		processer.ClientExec(rc,msg)
	}
}

func (rc *RpcClient) OnClose() {

}

func (rc *RpcClient) SendMessage(data proto.Message) {
	msg := NewMessage(0,data)
	if msg == nil {
		log.Error("msgid of %v not find",reflect.TypeOf(data).Name())
		return
	}

	rc.SendPbMessage(msg)
}

func (rc *RpcClient) SendPbMessage(msg IMessage) {
	bytes,err := Marshal(msg)
	if err != nil {
		log.Debug("notify err: %v", err)
		return
	}

	err = rc.conn.WriteMsg(bytes)
	if err != nil {
		log.Error("write message error: %v", err)
	}
}

func (rc *RpcClient) MsgProcess(msg IMessage) {
	if msg.GetMessageId() == uint16(pb.RpcMsgType_Handlers) {
		for _, handler := range msg.GetMessage().(*pb.RpcHandlers).GetHandlers() {
			rc.handlers[handler.GetMethod()] = handler
		}

		clientsLock.Lock()
		clients = append(clients, rc)
		clientsLock.Unlock()
	} else if msg.GetMessageId() == uint16(pb.RpcMsgType_Heart) {

	}
}

func (rc *RpcClient) Call(method string,request interface{}) (interface{},error) {
	handler,ok := rc.handlers[method]
	if !ok {
		return nil,fmt.Errorf("call handler [%v] not find",method)
	}

	requestId := GetMessageId(reflect.TypeOf(request))
	if requestId == 0 {
		return nil,fmt.Errorf("call handler [%v] request not registed",method)
	}

	if uint16(handler.GetRequestId()) != requestId {
		requestType := GetMessageType(uint16(handler.GetRequestId()))
		return nil,fmt.Errorf("call handler [%v] request must be [%v]",method, requestType.String())
	}

	if handler.GetReplyId() == 0 {
		return nil,fmt.Errorf("call handler [%v] responseid is zero",method)
	}

	cid := rc.GetCallId()
	c := NewRpcCall(method,cid,request.(proto.Message))
	rc.SendPbMessage(c.Msg.Request)
	rc.calls[cid] = c

	t := time.NewTimer(1*time.Second)
	for {
		select {
			case <- c.Done:
				delete(rc.calls, cid)
				if err := c.Msg.Response.GetError(); err != nil {
					return nil,fmt.Errorf("call handler [%v] response err-%v",method, err.Error())
				}

				response := c.Msg.Response.GetMessage()
				if response == nil {
					return nil,fmt.Errorf("call handler [%v] response is nil",method)
				}

				responseId := GetMessageId(reflect.TypeOf(response))
				if responseId != uint16(handler.GetReplyId()) {
					return nil,fmt.Errorf("call handler [%v] response type error",method)
				}

				return response,nil
			case <- t.C:
				t.Stop()
				delete(rc.calls, cid)
				return nil,fmt.Errorf("call handler [%v] response time out",method)

		}
	}
}

func (rc *RpcClient) Cast(method string,request interface{}) error {
	handler,ok := rc.handlers[method]
	if !ok {
		return fmt.Errorf("cast handler [%v] not find",method)
	}

	requestId := GetMessageId(reflect.TypeOf(request))
	if requestId == 0 {
		return fmt.Errorf("cast handler [%v] request not registed",method)
	}

	if uint16(handler.GetRequestId()) != requestId {
		requestType := GetMessageType(uint16(handler.GetRequestId()))
		return fmt.Errorf("cast handler [%v] request must be [%v]",method, requestType.Name())
	}

	if handler.GetReplyId() != 0 {
		return fmt.Errorf("cast handler [%v] responseid is not zero",method)
	}

	c := NewRpcCall(method,0,request.(proto.Message))
	rc.SendPbMessage(c.Msg.Request)

	return nil
}
