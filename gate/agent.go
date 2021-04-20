package gate

import (
	"net"
)

type Agent interface {
	WriteMsg(data [][]byte)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
}
