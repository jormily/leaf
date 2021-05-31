package rpc

import (
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/network"
	"math"
	"time"
)


type Client struct {
	rpcClient *network.TCPClient
}

func (this *Client) Init(addr string,pendingWriteNum int) {
	if addr == "" {
		log.Error("rpc client addr is error")
		return
	}

	this.rpcClient = new(network.TCPClient)
	this.rpcClient.Addr = addr
	this.rpcClient.ConnNum = 1
	this.rpcClient.ConnectInterval = 3 * time.Second
	this.rpcClient.PendingWriteNum = pendingWriteNum
	this.rpcClient.LenMsgLen = 4
	this.rpcClient.MaxMsgLen = math.MaxUint32
	this.rpcClient.NewAgent = func(conn *network.TCPConn) network.Agent {
		agent := newAgent(conn)
		return agent
	}
	this.rpcClient.AutoReconnect = true

	this.rpcClient.Start()
}

func (this *Client) Close() {

}

