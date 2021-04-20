package gate

import "github.com/golang/protobuf/proto"

type MessageType byte

const (
	MessageType_Request  MessageType = 0x00
	MessageType_Notify               = 0x01
	MessageType_Response             = 0x02
	MessageType_Push                 = 0x03
)

type MessageHead struct {
	Type       	MessageType
	ID         	uint32
	Error     	uint32
	MsgId      	uint16
}

type Message struct {
	head 		MessageHead
	RequestMsg 	proto.Message
	ReplyMsg   	proto.Message
	Error 		uint32
}