package gate

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	. "github.com/name5566/leaf/common"
	"github.com/name5566/leaf/log"
	. "github.com/name5566/leaf/msg"
	"reflect"
)


type MsgProcessor struct {
	handlers 		map[uint16]*MsgHandler
}

func NewMsgProcessor() *MsgProcessor {
	p := new(MsgProcessor)
	p.handlers = make(map[uint16]*MsgHandler)
	return p
}

func (p *MsgProcessor) Register(requestMsg interface{},replyMsg interface{},handler func(a interface{},msg *Messagex) *Error) {
	requestType := reflect.TypeOf(requestMsg)
	if !CheckMessage(requestType) {
		log.Error("Register requestType:%v check message err",requestType)
		return
	}

	var replyType reflect.Type
	var replyId uint16
	if replyMsg != nil {
		replyType = reflect.TypeOf(replyMsg)
		replyId = GetMessageId(replyType)
		if replyId == 0 {
			log.Error("Register replyType:%v check message err",replyType)
			return
		}
	}

	requestId := GetMessageId(requestType)
	if _,ok := p.handlers[requestId];ok {
		log.Error("Register handler of %v exist",requestId)
		return
	}

	msgHandler := new(MsgHandler)
	msgHandler.Func = handler
	msgHandler.RequestId = requestId
	msgHandler.RequestType = requestType
	msgHandler.ReplayId = replyId
	msgHandler.ReplyType = replyType

	p.handlers[requestId] = msgHandler
}

func (p *MsgProcessor) GetHandler(msgId uint16) *MsgHandler {
	if msgHandler,ok := p.handlers[msgId];ok {
		return msgHandler
	}
	return nil
}


func (p *MsgProcessor) Exec(a interface{},request IMessage) (*Messagex,error) {
	if request == nil {
		return nil,fmt.Errorf("msg request is nil")
	}

	if request.GetMessageType() != MessageType_Msg {
		return nil,fmt.Errorf("msg type is err")
	}

	msgId := request.GetMessageId()
	msgHandler := p.GetHandler(msgId)
	if msgHandler == nil {
		return nil,fmt.Errorf("handler id-%v not find",msgId)
	}

	var response *PbMessage = nil
	if msgHandler.ReplayId != 0 {
		response = new(PbMessage)
		response.MsgHead.Type = request.GetMessageType()
		response.MsgHead.ID = request.GetId()
		response.MsgHead.MsgId = msgHandler.ReplayId
		response.Data = reflect.New(msgHandler.ReplyType.Elem()).Interface().(proto.Message)
	}
	msg := new(Messagex)
	msg.Request = request
	msg.Response = response

	if err := msgHandler.Func(a,msg);err != nil {
		msg.Response.SetError(err.Code())
	}
	return msg,nil
}
