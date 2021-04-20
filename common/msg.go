package common

import (
	"reflect"
	"sync"
)

var mapMessage2Id = map[reflect.Type]uint16{}
var mapId2Message = map[uint16]reflect.Type{}
var registerOnce sync.Once


func Register(msg2id map[reflect.Type]uint16,id2msg map[uint16]reflect.Type)  {
	registerOnce.Do(func() {
		mapMessage2Id = msg2id
		mapId2Message = id2msg
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