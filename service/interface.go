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
}

type service struct {
	s 			IService
	closeSig	chan bool
	wg 			sync.WaitGroup
}

var servs []*service

func Register(s IService) {
	m := new(service)
	m.s = s
	m.closeSig = make(chan bool, 1)

	servs = append(servs, m)
}

func Init() {
	for _, s := range servs {
		s.s.Init(s.s)
	}

	for i := 0; i < len(servs); i++ {
		m := servs[i]
		m.wg.Add(1)
		go run(m)
	}
}

func Destroy() {
	for i := len(servs) - 1; i >= 0; i-- {
		m := servs[i]
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