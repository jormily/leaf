package rpcx

import (
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/network"
	"math"
	"time"
)

func Register(s IRpcService) {
	processer.Register(s)
}

func RpcConnect(addr string,pendingWriteNum int) bool {
	if addr == "" {
		log.Error("rpc client addr is error")
		return false
	}

	client := new(network.TCPClient)
	client.Addr = addr
	client.ConnNum = 1
	client.ConnectInterval = 3 * time.Second
	client.PendingWriteNum = pendingWriteNum
	client.LenMsgLen = 4
	client.MaxMsgLen = math.MaxUint32
	client.AutoReconnect = true
	client.NewAgent = func(conn *network.TCPConn) network.Agent {
		rpcClient := NewRpcClient(conn,client)
		return rpcClient
	}

	client.Start()
	return true
}

func RpcListen(addr string, pendingWriteNum int) {
	if addr == "" {
		log.Error("rpc Server addr is error")
		return
	}

	server = new(RpcServer)
	server.TCPServer = new(network.TCPServer)
	server.TCPServer.Addr = addr
	server.TCPServer.MaxConnNum = math.MaxInt32
	server.TCPServer.PendingWriteNum = pendingWriteNum
	server.TCPServer.LenMsgLen = 4
	server.TCPServer.MaxMsgLen = math.MaxUint32
	server.TCPServer.NewAgent = func(conn *network.TCPConn) network.Agent {
		rs := NewRpcSlient(conn)
		rs.SendMessage(processer.pb)
		return rs
	}

	server.TCPServer.Start()
}

func GetRpcClientAll() []*RpcClient {
	rpcClients := []*RpcClient{}
	clientsLock.Lock()
	for _,c := range clients {
		rpcClients = append(rpcClients, c)
	}
	clientsLock.Unlock()
	return rpcClients
}