package msg


import (
	"github.com/name5566/leaf/pb"
	"reflect"
	"sync"
)

var mapMessage2Id = map[reflect.Type]uint16{}
var mapId2Message = map[uint16]reflect.Type{}
var registerOnce sync.Once


func init() {
	mapMessage2Id[reflect.TypeOf(&pb.RpcNil{})] = 1
	mapMessage2Id[reflect.TypeOf(&pb.RpcHeart{})] = 2
	mapMessage2Id[reflect.TypeOf(&pb.RpcHandlers{})] = 3

	mapId2Message[1] = reflect.TypeOf(&pb.RpcNil{})
	mapId2Message[2] = reflect.TypeOf(&pb.RpcHeart{})
	mapId2Message[3] = reflect.TypeOf(&pb.RpcHandlers{})
}

func RegisterMessage(msg2id map[reflect.Type]uint16,id2msg map[uint16]reflect.Type)  {
	registerOnce.Do(func() {
		for k,v := range msg2id {
			mapMessage2Id[k] = v
		}

		for k,v := range id2msg {
			mapId2Message[k] = v
		}
	})
}

func CheckMessage(typ reflect.Type) bool {
	if _,ok := mapMessage2Id[typ];ok {
		return true
	}
	return false
}

func GetMessageId(typ reflect.Type) uint16 {
	if id,ok := mapMessage2Id[typ];ok {
		return id
	}
	return 0
}

func GetMessageType(msgid uint16) reflect.Type {
	if typ,ok := mapId2Message[msgid];ok {
		return typ
	}
	return nil
}