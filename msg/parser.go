package msg

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
)

type IParser interface {
	Unmarshal(data []byte) (IMessage,uint16,error)
	Marshal(msg IMessage) ([][]byte,error)
}


func Unmarshal(data []byte) (IMessage,error) {
	if len(data) < MessageHead_Flow_Len {
		return nil,fmt.Errorf("length of data is error")
	}

	msg := &PbMessage{}
	index := 0
	msg.MsgHead.Type = MessageType(data[index])
	index += 1
	msg.MsgHead.ID = binary.LittleEndian.Uint32(data[index:index+4])
	index += 4
	msg.MsgHead.Code = binary.LittleEndian.Uint16(data[index:index+2])
	index += 2
	if msg.GetCode() != 0 {
		return msg,nil
	}
	msg.MsgHead.MsgId = binary.LittleEndian.Uint16(data[index:index+2])
	index += 2
	msg.RpcHead.Len = binary.LittleEndian.Uint16(data[index:index+2])
	index += 2
	if msg.GetRpcLen() > 0 {
		if msg.GetMessageType() != MessageType_Rpc {
			return nil,fmt.Errorf("msg type error")
		}
		msg.RpcHead.Method = string(data[index:index+int(msg.GetRpcLen())])
	}
	index += int(msg.GetRpcLen())

	msgType := GetMessageType(msg.MsgHead.MsgId)
	if msgType == nil {
		return nil,fmt.Errorf("not find messageid of %v",msg.MsgHead.MsgId)
	}
	
	msg.Data = reflect.New(msgType.Elem()).Interface().(proto.Message)
	if err := proto.Unmarshal(data[index:],msg.Data);err != nil {
		return nil,err
	}
	return msg,nil
}

func Marshal(msg IMessage) ([]byte,error) {
	if msg.GetRpcLen() > 0 && msg.GetMessageType() != MessageType_Rpc {
		return nil, fmt.Errorf("msg rpc err")
	}

	var temp []byte
	var err error
	if msg.GetCode() == 0 {
		temp, err = proto.Marshal(msg.GetMessage().(proto.Message))
		if err != nil {
			return nil, err
		}
	}

	data := make([]byte, MessageHead_Flow_Len + int(msg.GetRpcLen()) + len(temp))
	index := 0
	data[0] = byte(msg.GetMessageType())
	index += 1
	binary.LittleEndian.PutUint32(data[index:index+4],msg.GetId())
	index += 4
	binary.LittleEndian.PutUint16(data[index:index+2],msg.GetCode())
	index += 2
	if msg.GetCode() != 0 {
		return data,nil
	}

	binary.LittleEndian.PutUint16(data[index:index+2],msg.GetMessageId())
	index += 2
	binary.LittleEndian.PutUint16(data[index:index+2],msg.GetRpcLen())
	index += 2
	copy(data[index:index+int(msg.GetRpcLen())],[]byte(msg.GetRpcMethod()))
	index += int(msg.GetRpcLen())

	copy(data[index:],temp)
	return data, nil
}
