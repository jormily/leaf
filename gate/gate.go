package gate

import (
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/network"
	"net"
	"time"
)

type Gate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool

	// agent
	GoLen              int
	TimerDispatcherLen int
	AsynCallLen        int
	ChanRPCLen         int
	OnAgentInit 	   func(Agent)
	OnAgentDestroy 	   func(Agent)
}

func (gate *Gate) Run(closeSig chan bool) {
	newAgent := func(conn network.Conn) network.Agent {
		a := &agent{conn: conn, gate: gate}
		return a
	}

	var wsServer *network.WSServer
	if gate.WSAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = gate.WSAddr
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.PendingWriteNum = gate.PendingWriteNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.CertFile = gate.CertFile
		wsServer.KeyFile = gate.KeyFile
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			return newAgent(conn)
		}
	}

	var tcpServer *network.TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			return newAgent(conn)
		}
	}

	if wsServer != nil {
		wsServer.Start()
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (gate *Gate) OnDestroy() {

}

type agent struct {
	conn     network.Conn
	gate     *Gate
	userData interface{}
}

func (a *agent) Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("%v",r)
		}
	}()


	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		msg,handler := Unmarshal(data)
		if msg == nil {
			continue
		}

		handler.Route(a,msg)
		if msg.head.Type != MessageType_Request {
			continue
		}

		bytes,err := Marshal(MessageType_Response,msg.head.ID,msg.Error,msg.ReplyMsg)
		if err != nil {
			log.Debug("handle message: %v", err)
			break
		}
		a.WriteMsg(bytes)
	}
}

func (a *agent) OnClose() {

}

func (a agent) Notify(msg proto.Message) {
	bytes,err := Marshal(MessageType_Notify,0,0,msg)
	if err != nil {
		log.Debug("handle message: %v", err)
		return
	}
	a.WriteMsg(bytes)
}

func (a *agent) WriteMsg(data [][]byte) {
	err := a.conn.WriteMsg(data...)
	if err != nil {
		log.Error("write message error: %v", err)
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}


