package rpcx

import (
	"reflect"

	. "github.com/name5566/leaf/msg"
)

type RpcHandler struct {
	Method 			*reflect.Method
	RequestId		uint16
	ReplayId		uint16
	RequestType 	reflect.Type
	ReplyType 		reflect.Type
	Service 		interface{}
}

//func FullMessagex(msg Messagex,responnse proto.Message) {
//	request := msg.Request
//	if request == nil {
//		return
//	}
//
//	if request.GetMessageType() == MessageType_Request {
//		if responnse == nil {
//			//todo:打印。。。。
//			return
//		}
//		msg.Request = NewPbResponseMessage(responnse,request.GetId())
//	}
//}
//
//func (handler *RpcHandler) CreateMessagex(req IMessage) (*Messagex,error) {
//	reqId := req.GetMessageId()
//
//	msg := &Messagex{}
//	msg.Request = req
//	if req.GetMessageType() != MessageType_Request && req.GetMessageType() != MessageType_Push {
//		return nil,fmt.Errorf("messageid-%v messageType err",reqId)
//	}
//
//	if req.GetMessageType() == MessageType_Notify {
//		return msg,nil
//	}
//
//	if handler.ReplyType != nil {
//		msg.Response = NewPbResponseMessage(nil,req.GetId())
//	}
//	return msg,nil
//}

func (handler *RpcHandler) Exec(s interface{},msg *Messagex) {
	//handler.Func(s,msg)
}
