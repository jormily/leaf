package rpc

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/log"
	"reflect"
)

var (
	rpcId2Msg = map[RpcMsgType]reflect.Type{
		RpcMsgType_Heart: reflect.TypeOf(&RpcHeart{}),
		RpcMsgType_Request: reflect.TypeOf(&RpcRequest{}),
		RpcMsgType_Response: reflect.TypeOf(&RpcResponse{}),
		RpcMsgType_Router: reflect.TypeOf(&RpcRouter{}),
	}

	rpcMsg2Id = map[string]RpcMsgType{}
)

func init() {
	for rpctype,typ := range rpcId2Msg {
		rpcMsg2Id[typ.Elem().Name()] = rpctype
	}
}

func Route(msgtype RpcMsgType,msg proto.Message, userData interface{}) {
	agent := userData.(*Agent)
	switch msgtype {
	case RpcMsgType_Heart:
		agent.OnHeart(msg)
	case RpcMsgType_Request:
		agent.OnRequset(msg)
	case RpcMsgType_Response:
		agent.OnResponse(msg)
	case RpcMsgType_Router:
		agent.OnRouter(msg)
	default:
		log.Error("Route unknow msgtype = %v",msgtype)
	}
}

func Unmarshal(data []byte) (proto.Message, RpcMsgType, error) {
	msgtype := RpcMsgType(data[0])
	typ,ok := rpcId2Msg[msgtype]
	if !ok {
		return nil,msgtype,fmt.Errorf("Unmarshal msgtype [%v] is error",msgtype)
	}

	msg := reflect.New(typ.Elem()).Interface()
	err := proto.Unmarshal(data[1:],msg.(proto.Message))
	if err != nil {
		return nil,msgtype,fmt.Errorf("Unmarshal msgtype [%v] unmarshal error",msgtype)
	}

	return msg.(proto.Message),msgtype,nil
}

func Marshal(msg proto.Message) ([][]byte, error) {
	typ := reflect.TypeOf(msg)
	if typ.Kind() != reflect.Ptr {
		return nil,fmt.Errorf("Marshal msg is not ptr")
	}

	msgName := typ.Elem().Name()
	msgtype,ok := rpcMsg2Id[msgName]
	if !ok {
		return nil,fmt.Errorf("Marshal msg [%v] not reg",msgName)
	}
	msgtypes := []byte{byte(msgtype)}

	data,err := proto.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("Marshal msg err:%v",err.Error())
	}

	return [][]byte{msgtypes, data},nil
}