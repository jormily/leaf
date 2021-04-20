package rpc

import (
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/network"
	"math"
)

type Server struct {
	rpcServer *network.TCPServer
}


func (this *Server) Init(addr string, pendingWriteNum int) {
	if addr == "" {
		log.Error("rpc Server addr is error")
		return
	}

	this.rpcServer = new(network.TCPServer)
	this.rpcServer.Addr = addr
	this.rpcServer.MaxConnNum = math.MaxInt32
	this.rpcServer.PendingWriteNum = pendingWriteNum
	this.rpcServer.LenMsgLen = 4
	this.rpcServer.MaxMsgLen = math.MaxUint32
	this.rpcServer.NewAgent = func(conn *network.TCPConn) network.Agent {
		agent := newAgent(conn)
		rpcRouterLocker.RLock()
		agent.(*Agent).WriteMsg(rpcRouter)
		rpcRouterLocker.RUnlock()
		return agent
	}

	this.rpcServer.Start()
}
