package rpcx

import (
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/pb"
	"strings"

	"github.com/name5566/leaf/log"
	. "github.com/name5566/leaf/common"
	. "github.com/name5566/leaf/msg"
	"reflect"
)

type IServiceFilter interface {
	GetService(sid uint16) interface{}
}

type IRpcService interface {
	GetServiceName() string
}

type IGetRpcChannel interface {
	GetRpcChannel() *RpcChannel
}

type RpcProcessor struct {
	handlers 		map[string]map[string]*RpcHandler
	services 		map[string]interface{}
	pb				*pb.RpcHandlers
	serviceFilter 	IServiceFilter
}

var processer = NewRpcProcessor()

func NewRpcProcessor() *RpcProcessor {
	p := new(RpcProcessor)
	p.handlers = make(map[string]map[string]*RpcHandler)
	p.services = make(map[string]interface{})
	p.pb = &pb.RpcHandlers{
		Handlers: []*pb.RpcHandler{},
	}
	return p
}

//func (p *RpcProcessor) GetServiceMsgHandlers(sname string) map[string]*RpcHandler {
//	if msgHandlers,ok := p.handlers[sname];ok {
//		return msgHandlers
//	}else{
//		msgHandlers = map[string]*RpcHandler{}
//		p.handlers[sname] = msgHandlers
//		return msgHandlers
//	}
//}

func  (p *RpcProcessor) IsRpcHandler(m reflect.Method) bool {
	methodType := m.Type
	if !strings.HasPrefix(m.Name,"RPC_") {
		return false
	}

	if methodType.NumIn() != 2 {
		return false
	}

	if methodType.In(0).Kind() != reflect.Ptr || !CheckMessage(methodType.In(1))  {
		return false
	}

	if methodType.NumOut() == 1 && methodType.Out(0).Kind() == reflect.Interface {
		return true
	}

	if methodType.NumOut() == 2 && methodType.Out(1).Kind() == reflect.Interface {
		if CheckMessage(methodType.Out(0)) {
			return true
		}
	}

	return false
}


func (p *RpcProcessor) Register(s IRpcService) {
	serviceType := reflect.TypeOf(s)
	serviceName := s.GetServiceName()
	if serviceType.Kind() != reflect.Ptr {
		log.Error("RpcProcessor register error: IRpcService-[%v] not ptr",serviceName)
		return
	}
	p.services[serviceName] = s

	rpcHandlerMap :=  make(map[string]*RpcHandler)
	for i:=0;i<serviceType.NumMethod();i++{
		method := serviceType.Method(i)
		if p.IsRpcHandler(method) {
			rpcHandler := new(RpcHandler)
			rpcHandler.Method = &method
			rpcHandler.RequestType = method.Type.In(1)
			rpcHandler.RequestId = GetMessageId(rpcHandler.RequestType)
			if method.Type.NumOut() == 1 {
				rpcHandler.ReplyType = nil
				rpcHandler.ReplayId = 0
			}else{
				rpcHandler.ReplyType = method.Type.Out(0)
				rpcHandler.ReplayId = GetMessageId(rpcHandler.ReplyType)
			}
			rpcHandler.Service = s
			methodName := strings.TrimPrefix(method.Name,"RPC_")
			rpcHandlerMap[methodName] = rpcHandler


			p.pb.Handlers = append(p.pb.Handlers,&pb.RpcHandler{
				Method: proto.String(serviceName + "."+ methodName),
					ReplyId: proto.Uint32(uint32(rpcHandler.ReplayId)),
					RequestId: proto.Uint32(uint32(rpcHandler.RequestId)),
			})
		}
	}

	p.handlers[serviceName] = rpcHandlerMap
}

func (p *RpcProcessor) GetHandler(sname string) (interface{},*RpcHandler) {
	slist := strings.Split(sname,".")
	if len(slist) != 2 {
		return nil,nil
	}

	serviceName := slist[0]
	methodName := slist[1]

	rpcHandlerMap,ok := p.handlers[serviceName]
	if !ok {
		return nil,nil
	}

	rpcHandler,ok := rpcHandlerMap[methodName]
	if !ok {
		return nil,nil
	}

	service,ok := p.services[serviceName]
	if !ok {
		return nil,nil
	}

	return service,rpcHandler
}


func (p *RpcProcessor) ClientExec(rc *RpcClient,msg IMessage) {
	if msg == nil {
		return
	}

	if msg.GetMessageType() == MessageType_Msg {
		rc.MsgProcess(msg)
		return
	}

	if msg.GetMessageType() == MessageType_Rpc {
		c,ok := rc.calls[msg.GetId()]
		if !ok {
			log.Error("rpc msg id-%v not find")
			return
		}

		c.Msg.Response  = msg
		c.Done <- c
		return
	}

	log.Error("msg type err")
}


func (p *RpcProcessor) ServerExec(rs *RpcServerClient,request IMessage) {
	if request == nil {
		log.Error("msg request is nil")
		return
	}

	if request.GetMessageType() == MessageType_Msg {
		return
	}

	if request.GetMessageType() == MessageType_Rpc {
		var response *PbMessage
		if request.GetId() != 0 {
			response = new(PbMessage)
			response.MsgHead.Type = MessageType_Rpc
			response.MsgHead.ID = request.GetId()
		}

		service,rpcHandler := p.GetHandler(request.GetRpcMethod())
		if rpcHandler == nil || rpcHandler.Method == nil {
			if response != nil {
				response.MsgHead.Code = Error_RpcNotFind.Code()
				rs.SendPbMessage(response)
			}

			return
		}

		if rpcHandler.RequestId != request.GetMessageId() {
			if response != nil {
				response.MsgHead.Code = Error_RpcRespType.Code()
				rs.SendPbMessage(response)
			}
			return
		}

		if rpcHandler.RequestId == 0 {
			if response != nil {
				response.MsgHead.Code = Error_RpcCallErr.Code()
				rs.SendPbMessage(response)
			}
			return
		}

		if response != nil {
			response.MsgHead.MsgId = rpcHandler.ReplayId
		}

		msg := &Messagex{
			Request: request,
			Response: response,
		}

		service.(IGetRpcChannel).GetRpcChannel().Exec(rpcHandler.Method,msg)
		if msg.Response != nil {
			rs.SendPbMessage(msg.Response)
		}
	}

}


