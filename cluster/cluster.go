package cluster

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/conf"
	"github.com/name5566/leaf/pb"
	"github.com/name5566/leaf/rpc"
	"github.com/name5566/leaf/service"
	"sync"
)


var (
	rpcServer  				*rpc.Server
	rpcMasterClient			*rpc.Client
	rpcClients 				map[int32]*rpc.Client

	statusLock				sync.Mutex
	NodeInfos				[]*NodeInfo
	statusVersion 			uint32 = 0
	statusLocal				*pb.ServerStatus

	masterSvr				*Master
	masterClientSvr			*MasterClient
)

func init() {

}

func Register() {
	if conf.ListenAddr != "" {
		rpcServer = new(rpc.Server)
		rpcServer.Init(conf.ListenAddr,conf.PendingWriteNum)
	}

	NodeInfos = []*NodeInfo{}
	statusLocal = &pb.ServerStatus{}
	statusLocal.Addr = proto.String(conf.ListenAddr)
	statusLocal.Sid = proto.Int32(conf.ServerId)
	statusLocal.Stype = proto.String(conf.ServerType)
	NodeInfos = append(NodeInfos, &NodeInfo{
		ServerStatus: statusLocal,
	})

	if conf.ServerType != "master" {
		status := &pb.ServerStatus{}
		status.Addr = proto.String(conf.MasterAddr)
		status.Stype = proto.String("master")
		NodeInfos = append(NodeInfos, NewNodeInfo(status))
	}


	if conf.ServerType == "master" {
		statusVersion = 1
		masterSvr = NewMaster()
		service.Register(masterSvr)
	}else{
		masterClientSvr = NewMasterClient()
		service.Register(masterClientSvr)
	}
}


func Init() {
	//
	if conf.ServerType != "master" {
		masterClientSvr.OnInit()
	}
}


func Destroy() {
	if rpcServer != nil {
		rpcServer.Close()
	}

	for _, client := range rpcClients {
		client.Close()
	}
}

func Call(stype string,method string,requestMsg interface{}) (interface{},error) {
	nodes := GetNodeByType(stype)
	if len(nodes) == 0 {
		return nil,fmt.Errorf("Call method %v.%v,not find node",stype,method)
	}

	return nodes[0].Call(method,requestMsg)
}

func Cast(stype string,method string,requestMsg interface{}) (error) {
	nodes := GetNodeByType(stype)
	if len(nodes) == 0 {
		return fmt.Errorf("Call method %v.%v,not find node",stype,method)
	}

	return nodes[0].Cast(method,requestMsg)
}

func BroadCast(stype string,method string,requestMsg interface{}) {
	nodes := GetNodeByType(stype)
	for _,node := range nodes {
		node.Cast(method,requestMsg)
	}
}

func BroadCastAll(method string,requestMsg interface{},sids...interface{}) {
	for _,node := range NodeInfos {
		for _,sid := range sids {
			if sid != node.GetSid() {
				node.Cast(method,requestMsg)
			}
		}
	}
}