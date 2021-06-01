package cluster

import (
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/conf"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/pb"
	"github.com/name5566/leaf/service"
	"time"
)

type Master struct {
	*service.Service
	heartResponse 		*pb.HeartRespose
}

func NewMaster() *Master {
	m := new(Master)
	m.Service = service.NewService()
	m.heartResponse = &pb.HeartRespose{
		Version: proto.Uint32(statusVersion),
		StatusList: []*pb.ServerStatus{},
	}

	return m
}

func updateStatusVersion() {
	statusVersion++
	if statusVersion == 0 {
		statusVersion++
	}
}

func (m *Master) RPC_Heart(req *pb.HeartRequest) (*pb.HeartRespose,error) {
	res := &pb.HeartRespose{}
	newStatus := req.GetStatus()
	statusLock.Lock()
	defer statusLock.Unlock()

	var findFlag bool = false
	var changeFlag bool = false
	for index,status := range NodeInfos {
		if status.GetSid() == newStatus.GetSid() {
			if status.GetAddr() != newStatus.GetAddr() || status.GetLoad() != newStatus.GetLoad() {
				NodeInfos[index].ServerStatus = newStatus
				updateStatusVersion()
				changeFlag = true
			}
			findFlag = true
		}
	}

	if !findFlag {
		NodeInfos = append(NodeInfos,NewNodeInfo(newStatus))
		updateStatusVersion()
		changeFlag = true
	}

	statusArray := []*pb.ServerStatus{}
	for _,v := range NodeInfos {
		statusArray = append(statusArray,v.ServerStatus)
	}

	res.Version = proto.Uint32(statusVersion)
	if statusVersion != req.GetVersion() {
		res.StatusList = statusArray
	}

	if changeFlag {
		BroadCastAll("MasterClient.Update",res,req.Status.GetSid())
	}

	return res,nil
}

type MasterClient struct {
	*service.Service
	version 			uint32
	ServerStatuseMap 	map[int32]*pb.ServerStatus
}

func NewMasterClient() *MasterClient {
	m := new(MasterClient)
	m.Service = service.NewService()
	m.version = 0
	m.ServerStatuseMap = map[int32]*pb.ServerStatus{}
	return m
}

func (m *MasterClient) OnInit() {
	timerFunc := func() {
		res, err := Call("master", "Master.Heart", &pb.HeartRequest{
			Version: proto.Uint32(statusVersion),
			Status:  statusLocal,
		})
		if err != nil {
			log.Release("MasterClient OnInit call Master.Heart err:%v", err)
		}else{
			m.RPC_Update(res.(*pb.HeartRespose))
		}
	}

	timerFunc()
	m.AfterFunc(3*time.Second, timerFunc)
}

func (m *MasterClient) RPC_Update(req *pb.HeartRespose) error {
	statusLock.Lock()
	defer statusLock.Unlock()

	for _,status := range req.StatusList {
		if status.GetSid() == conf.ServerId {
			continue
		}

		findFlag := false
		for _,node := range NodeInfos {
			if node.GetSid() == status.GetSid() {
				node.ServerStatus = status
				findFlag = true
			}
		}

		if !findFlag {
			for _,stype := range conf.RpcStypes {
				if stype == status.GetStype() {
					NodeInfos = append(NodeInfos,NewNodeInfo(status))
				}
			}
		}
	}
	statusVersion = req.GetVersion()

	return nil
}
