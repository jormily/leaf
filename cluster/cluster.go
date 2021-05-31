package cluster

import (
	"github.com/name5566/leaf/conf"
	"github.com/name5566/leaf/network"
	"github.com/name5566/leaf/rpc"
	"github.com/name5566/leaf/service"
)

var (
	rpcServer  		*rpc.Server
	rpcMasterClient	*rpc.Client
	rpcClients 		[]*rpc.Client
)

func Init() {
	if conf.ListenAddr != "" {
		rpcServer = new(rpc.Server)
		rpcServer.Init(conf.ListenAddr,conf.PendingWriteNum)
	}

	if conf.IsMaster {
		m := NewMaster()
		service.Register(m)
	}else{
		rpcMasterClient = new(rpc.Client)
		rpcMasterClient.Init(conf.MasterAddr,conf.PendingWriteNum)
		m := NewMasterClient()
		service.Register(m)
	}

	//for _, addr := range conf.ConnAddrs {
	//	client := new(network.TCPClient)
	//	client.Addr = addr
	//	client.ConnNum = 1
	//	client.ConnectInterval = 3 * time.Second
	//	client.PendingWriteNum = conf.PendingWriteNum
	//	client.LenMsgLen = 4
	//	client.MaxMsgLen = math.MaxUint32
	//	client.NewAgent = newAgent
	//
	//	client.Start()
	//	clients = append(clients, client)
	//}

}

func Destroy() {
	if rpcServer != nil {
		rpcServer.Close()
	}

	for _, client := range rpcClients {
		client.Close()
	}
}

type Agent struct {
	conn *network.TCPConn
}

func newAgent(conn *network.TCPConn) network.Agent {
	a := new(Agent)
	a.conn = conn
	return a
}

func (a *Agent) Run() {}

func (a *Agent) OnClose() {}
