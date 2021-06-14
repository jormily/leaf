package msg

import (
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/common"
	"github.com/name5566/leaf/log"
	"reflect"
	"sync"
)

var msgRid uint32 = 1
var msgLock sync.Mutex

const (
	MessageHead_Flow_Len = 11
)

func GetMessageRid() uint32 {
	msgLock.Lock()
	defer msgLock.Unlock()
	msgRid++
	if msgRid == 0 {
		msgRid++
	}
	return msgRid
}


type MessageType byte

type IMessage interface {
	GetMessageType() MessageType
	GetMessage() interface{}
	GetId() uint32
	GetCode() uint16
	GetError() error
	GetMessageId() uint16
	GetRpcLen() uint16
	GetRpcMethod() string

	SetMessage(msg interface{})
	SetError(code uint16)
	SetRpcMethod(method string)
}

type Messagex struct {
	Request		IMessage
	Response  	IMessage
}

type MessageHead struct {
	Type       	MessageType
	ID         	uint32
	Code     	uint16
	MsgId      	uint16
}

type RpcHead struct {
	Len 	uint16
	Method 	string
}

type PbMessage struct {
	MsgHead		MessageHead
	RpcHead		RpcHead
	Data 		proto.Message
}

func NewMessage(id uint32,data proto.Message) *PbMessage {
	msg := new(PbMessage)
	msg.MsgHead.Type = MessageType_Msg
	msg.MsgHead.MsgId = GetMessageId(reflect.TypeOf(data))
	if msg.MsgHead.MsgId == 0 {
		log.Error("not find msg %v",reflect.TypeOf(data).Name())
		return nil
	}
	msg.MsgHead.ID = id
	msg.Data = data
	return msg
}

func NewRpcMessage(method string,id uint32,data proto.Message) *PbMessage {
	msg := new(PbMessage)
	msg.MsgHead.Type = MessageType_Rpc
	msg.MsgHead.MsgId = GetMessageId(reflect.TypeOf(data))
	if msg.MsgHead.MsgId == 0 {
		return nil
	}
	msg.MsgHead.ID = id
	msg.Data = data
	msg.RpcHead.Method = method
	msg.RpcHead.Len = uint16(len(method))
	return msg
}

//func NewPbMessage(msgType MessageType,data proto.Message,rid uint32) (*PbMessage) {
//	pbMsg := new(PbMessage)
//	pbMsg.MsgHead.Type =  typ
//	pbMsg.MsgHead.MsgId = GetMessageId(reflect.TypeOf(data))
//	if pbMsg.MsgHead.MsgId == 0 {
//		return nil
//	}
//
//	if typ == MessageType_Request {
//		pbMsg.MsgHead.ID = GetMessageRid()
//	}else if typ == MessageType_Response {
//		pbMsg.MsgHead.ID = rid
//	}
//	pbMsg.Data = data
//
//	return pbMsg
//}
//
//func NewPbReqestMessage(data proto.Message) (*PbMessage) {
//	return NewPbMessage(MessageType_Request,data,0)
//}
//
//func NewPbResponseMessage(data proto.Message,rid uint32) (*PbMessage) {
//	return NewPbMessage(MessageType_Request,data,rid)
//}
//
//func NewPbNotifyMessage(data proto.Message) (*PbMessage) {
//	return NewPbMessage(MessageType_Notify,data,0)
//}
//
//func NewPbPushMessage(data proto.Message) (*PbMessage) {
//	return NewPbMessage(MessageType_Push,data,0)
//}

func (m *PbMessage) GetMessageType() MessageType {
	return m.MsgHead.Type
}

func (m *PbMessage) GetMessage() interface{} {
	return m.Data
}

func (m *PbMessage) GetCode() uint16 {
	return m.MsgHead.Code
}

func (m *PbMessage) GetError() error {
	if m.MsgHead.Code == 0 {
		return nil
	}
	return common.GetError(m.MsgHead.Code)
}

func (m *PbMessage) GetMessageId() uint16 {
	return m.MsgHead.MsgId
}

func (m *PbMessage) GetId() uint32 {
	return m.MsgHead.ID
}

func (m *PbMessage) SetError(code uint16) {
	m.MsgHead.Code = code
}

func (m *PbMessage) GetRpcLen() uint16 {
	return m.RpcHead.Len
}

func (m *PbMessage) GetRpcMethod() string {
	return m.RpcHead.Method
}

func (m *PbMessage) SetRpcMethod(method string) {
	m.RpcHead.Method = method
}

func (m *PbMessage) SetMessage(msg interface{}) {
	m.Data = msg.(proto.Message)
}

