package rpcx

import (
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/log"
	. "github.com/name5566/leaf/msg"
	"github.com/name5566/leaf/network"
	"reflect"
	"sync"
)

var server *RpcServer

type RpcServer struct {
	*network.TCPServer
}

type RpcServerClient struct {
	sync.Mutex
	conn 		network.Conn

	calls 		map[uint32]*RpcCall
}

func NewRpcSlient(conn network.Conn) *RpcServerClient {
	rs := new(RpcServerClient)
	rs.conn = conn
	rs.calls = make(map[uint32]*RpcCall)
	return rs
}


func (rs *RpcServerClient) Run() {
	for {
		data, err := rs.conn.ReadMsg()
		if err != nil {
			log.Error("read message: %v", err)
			break
		}
		msg, err := Unmarshal(data)
		if err != nil {
			log.Error("unmarshal message error: %v", err)
			break
		}

		processer.ServerExec(rs,msg)
	}
}

func (rs *RpcServerClient) OnClose() {

}

func (rs *RpcServerClient) SendMessage(data proto.Message) {
	msg := NewMessage(0,data)
	if msg == nil {
		log.Error("msgid of %v not find",reflect.TypeOf(data).Name())
		return
	}

	rs.SendPbMessage(msg)
}

func (rs *RpcServerClient) SendPbMessage(msg IMessage) {
	bytes,err := Marshal(msg)
	if err != nil {
		log.Debug("notify err: %v", err)
		return
	}

	err = rs.conn.WriteMsg(bytes)
	if err != nil {
		log.Error("write message error: %v", err)
	}
}

func (rs *RpcServerClient) MsgProcess(msg IMessage) {
	return
}
