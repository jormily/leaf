package common

//var mapMessage2Id = map[reflect.Type]uint16{}
//var mapId2Message = map[uint16]reflect.Type{}
//var registerOnce sync.Once
//
//
//func init() {
//	mapMessage2Id[reflect.TypeOf(&pb.HeartRequest{})] = 1
//	mapMessage2Id[reflect.TypeOf(&pb.HeartRespose{})] = 2
//
//	mapId2Message[1] = reflect.TypeOf(&pb.HeartRequest{})
//	mapId2Message[2] = reflect.TypeOf(&pb.HeartRespose{})
//}
//
//func RegisterMessage(msg2id map[reflect.Type]uint16,id2msg map[uint16]reflect.Type)  {
//	registerOnce.Do(func() {
//		for k,v := range msg2id {
//			mapMessage2Id[k] = v
//		}
//
//		for k,v := range id2msg {
//			mapId2Message[k] = v
//		}
//	})
//}
//
//func CheckMessage(typ reflect.Type) bool {
//	if _,ok := mapMessage2Id[typ];ok {
//		return true
//	}
//	return false
//}
//
//func GetMessageId(typ reflect.Type) uint16 {
//	if id,ok := mapMessage2Id[typ];ok {
//		return id
//	}
//	return 0
//}
//
//func GetMessageType(msgid uint16) reflect.Type {
//	if typ,ok := mapId2Message[msgid];ok {
//		return typ
//	}
//	return nil
//}