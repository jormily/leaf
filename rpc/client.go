package rpc

import (
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/network"
	"math"
	"time"
)


type Client struct {
	*Agent
	rpcClient 	*network.TCPClient
}

func (this *Client) Init(addr string,pendingWriteNum int,) bool {
	if addr == "" {
		log.Error("rpc client addr is error")
		return false
	}

	this.rpcClient = new(network.TCPClient)
	this.rpcClient.Addr = addr
	this.rpcClient.ConnNum = 1
	this.rpcClient.ConnectInterval = 1 * time.Second
	this.rpcClient.PendingWriteNum = pendingWriteNum
	this.rpcClient.LenMsgLen = 4
	this.rpcClient.MaxMsgLen = math.MaxUint32
	this.rpcClient.AutoReconnect = true
	this.rpcClient.NewAgent = func(conn *network.TCPConn) network.Agent {
		agent := newAgent(conn, func(a *Agent) {
			this.Agent = a

		})

		return agent
	}


	this.rpcClient.Start()

	return true
}

func (this *Client) Close() {

}

