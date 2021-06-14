package common

import (
	"fmt"
	"github.com/name5566/leaf/log"
	"reflect"
	"runtime"
)

func IsNil(inf interface{}) bool {
	val := reflect.ValueOf(inf)
	if val.Kind() == reflect.Ptr {
		return val.IsNil()
	} else {
		return inf == nil
	}
}


func Try(fun func(), handler func(interface{})) {
	defer func() {
		if r := recover(); r != nil {
			if handler == nil {
				buf := make([]byte, 4096)
				l := runtime.Stack(buf, true)
				err := fmt.Errorf("%v: %s", r, buf[:l])
				log.Error("error stack:%v \n",err)
			} else {
				handler(r)
			}
		}
	}()
	fun()
}
