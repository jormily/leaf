package service

import (
	"github.com/name5566/leaf/conf"
	"github.com/name5566/leaf/log"
	"runtime"
	"sync"
)

type IService interface {
	Init(s interface{})
	Destroy()
	Run(closeSig chan bool)
	//GetRpcHandler() *rpcx.RpcHandler
	GetServiceType() string
	GetServiceName() string
}

type service struct {
	s 			IService
	closeSig	chan bool
	wg 			sync.WaitGroup
}

var servlist = []*service{}
var servmap  = map[string]*service{}

func Register(s IService) {
	m := new(service)
	m.s = s
	m.closeSig = make(chan bool, 1)

	servlist = append(servlist, m)
}

func Init() {
	for _, s := range servlist {
		s.s.Init(s.s)
	}
}

func Start() {
	for i := 0; i < len(servlist); i++ {
		m := servlist[i]
		m.wg.Add(1)
		go run(m)
	}
}

func Destroy() {
	for i := len(servlist) - 1; i >= 0; i-- {
		m := servlist[i]
		m.closeSig <- true
		m.wg.Wait()
		destroy(m)
	}
}

func run(m *service) {
	m.s.Run(m.closeSig)
	m.wg.Done()
}

func destroy(m *service) {
	defer func() {
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("%v", r)
			}
		}
	}()

	m.s.Destroy()
}