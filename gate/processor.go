package gate

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/common"
	"github.com/name5566/leaf/log"
	"reflect"
)

//type gateHandler func(a Agent,msg proto.Message) (proto.Message,error)

type Handler struct {
	Cmd 			uint16
	Func 			func(a Agent,msg *Message)
	RequestId		uint16
	ReplyId			uint16
	RequestType 	reflect.Type
	ReplyType 		reflect.Type
}
var handlerMap = map[uint16]*Handler{}

func Register(requestMsg proto.Message,replyMsg proto.Message,handler func(a Agent,msg *Message)) {
	requestType := reflect.TypeOf(requestMsg)
	if !common.CheckMessage(requestType) {
		log.Error("Register requestType:%v check message err",requestType)
		return
	}

	var replyType reflect.Type
	if replyMsg != nil {
		replyType = reflect.TypeOf(replyMsg)
		if !common.CheckMessage(replyType) {
			log.Error("Register replyType:%v check message err",replyType)
			return
		}
	}

	requestId := common.GetMessageId(requestType)
	if _,ok := handlerMap[requestId];ok {
		log.Error("Register handler of %v exist",requestId)
		return
	}

	h := new(Handler)
	h.Cmd = requestId
	h.Func = handler
	h.RequestType = requestType
	h.ReplyType = replyType
	h.RequestId = requestId
	if replyType != nil {
		h.ReplyId = common.GetMessageId(replyType)
	}

	handlerMap[requestId] = h
}


func (h *Handler) Route(a Agent,msg *Message) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Route [%v-%v] error~",msg.head.MsgId,msg.head.ID)
		}
	}()

	h.Func(a,msg)
}


func Unmarshal(data []byte) (*Message,*Handler) {
	if len(data) < 11 {
		return nil,nil
	}

	msg := &Message{}
	msg.head.Type = MessageType(data[0])
	msg.head.ID = binary.LittleEndian.Uint32(data[1:5])
	msg.head.MsgId = binary.LittleEndian.Uint16(data[9:11])
	if msg.head.Type != MessageType_Request && msg.head.Type != MessageType_Push {
		return nil,nil
	}

	handler,ok := handlerMap[msg.head.MsgId]
	if !ok {
		return nil,nil
	}

	if handler.RequestType != nil {
		msg.RequestMsg = reflect.New(handler.RequestType.Elem()).Interface().(proto.Message)
		if err := proto.Unmarshal(data[11:],msg.RequestMsg);err != nil {
			return nil,nil
		}
	}
	if handler.ReplyType != nil {
		msg.ReplyMsg = reflect.New(handler.ReplyType.Elem()).Interface().(proto.Message)
	}

	return msg,handler
}


func Marshal(msgType MessageType,id uint32,errCode uint32,msg proto.Message) ([][]byte,error) {
	headArray := make([]byte,11)
	headArray[0] = byte(msgType)
	binary.LittleEndian.PutUint32(headArray[1:5],id)
	if errCode != 0 {
		binary.LittleEndian.PutUint32(headArray[5:9],errCode)
		return [][]byte{headArray},nil
	}

	if msgId := common.GetMessageId(reflect.TypeOf(msg));msgId == 0 {
		return nil, fmt.Errorf("ss")
	}else {
		binary.LittleEndian.PutUint32(headArray[5:9], uint32(msgId))

		msgArray, err := proto.Marshal(msg)
		if err != nil {
			return nil, err
		}
		return [][]byte{headArray, msgArray}, nil
	}
}
