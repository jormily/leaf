package cluster

import (
	"fmt"
	"github.com/name5566/leaf/conf"
	"github.com/name5566/leaf/pb"
	"github.com/name5566/leaf/rpc"
)

type NodeInfo struct {
	*pb.ServerStatus
	Client 		*rpc.Client
}

func NewNodeInfo(status *pb.ServerStatus) *NodeInfo {
	node := new(NodeInfo)
	node.ServerStatus = status
	node.Client = new(rpc.Client)
	node.Client.Init(node.GetAddr(),1)
	return node
}

func GetNodeByType(stype string) []*NodeInfo {
	statusLock.Lock()
	defer statusLock.Unlock()

	nodes := []*NodeInfo{}
	for _,node := range NodeInfos {
		if node.GetStype() == stype {
			nodes = append(nodes, node)
		}
	}

	return nodes
}


func (node *NodeInfo) Call(method string,requestMsg interface{}) (interface{},error) {
	if node.GetSid() == conf.ServerId {
		return rpc.CallLocal(method,requestMsg)
	}

	if node.Client.Agent == nil {
		return nil,fmt.Errorf("node call method:%v agent is nil",method)
	}

	return node.Client.Agent.Call(method,requestMsg)
}

func (node *NodeInfo) Cast(method string,requestMsg interface{}) error {
	if node.GetSid() == conf.ServerId {
		return rpc.CastLocal(method,requestMsg)
	}

	if node.Client.Agent == nil {
		return fmt.Errorf("node call method:%v agent is nil",method)
	}

	return node.Client.Agent.Cast(method,requestMsg)
}
