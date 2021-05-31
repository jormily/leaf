package cluster

import (
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/pb"
	"github.com/name5566/leaf/rpc"
	"github.com/name5566/leaf/service"
)

type Master struct {
	*service.Service
	ServerStatuseMap 	map[int32]*pb.RpcServerStatus
	version 			uint32

	heartResponse 		*pb.RpcHeartRespose
}

//func init() {
//	var s service.IService = NewMaster()
//	var ss service.IService = service.NewService()
//}

func NewMaster() *Master {
	m := new(Master)
	m.Service = service.NewService()
	m.version = 1
	m.ServerStatuseMap = map[int32]*pb.RpcServerStatus{}
	m.heartResponse = &pb.RpcHeartRespose{
		Version: proto.Uint32(m.version),
		StatusList: []*pb.RpcServerStatus{},
	}

	return m
}

func (m *Master) UpdateVersion() {
	m.version++
	if m.version == 0 {
		m.version++
	}

	for _,v := range m.ServerStatuseMap {
		m.heartResponse.StatusList = append(m.heartResponse.StatusList, v)
	}
}

func (m *Master) RPC_Heart(req *pb.RpcHeartRequest) (error,*pb.RpcHeartRespose) {
	res := &pb.RpcHeartRespose{}
	newStatus := req.GetStatus()
	if status,ok := m.ServerStatuseMap[newStatus.GetSid()];ok {
		if status.GetAddr() != newStatus.GetAddr() {
			m.ServerStatuseMap[newStatus.GetSid()] = newStatus
			m.UpdateVersion()
		}else{
			if status.GetLoad() != newStatus.GetLoad() {
				status.Load = proto.Int32(newStatus.GetLoad())
				m.UpdateVersion()
			}
		}
	}else{
		m.ServerStatuseMap[newStatus.GetSid()] = newStatus
		m.UpdateVersion()
	}

	res.Version = proto.Uint32(m.version)
	if m.version == req.GetVersion() {
		return nil,res
	}else{
		return nil,m.heartResponse
	}
}


type MasterClient struct {
	*service.Service
	version 			uint32
	ServerStatuseMap 	map[int32]*pb.RpcServerStatus

	rpcClient			rpc.Client
}

func NewMasterClient() *MasterClient {
	m := new(MasterClient)
	m.Service = service.NewService()
	m.version = 0
	m.ServerStatuseMap = map[int32]*pb.RpcServerStatus{}
	return m
}

func (m *MasterClient) RPC_GetNode() error {
	return nil
}



