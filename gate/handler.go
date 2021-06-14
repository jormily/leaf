package gate

import (
	. "github.com/name5566/leaf/common"
	. "github.com/name5566/leaf/msg"
	"reflect"
)

//
type MsgHandler struct {
	Func 			func(s interface{},msg *Messagex) *Error
	RequestId		uint16
	ReplayId		uint16
	RequestType 	reflect.Type
	ReplyType 		reflect.Type
}

