package common

import "sync"

type Error struct {
	Id  uint16
	Str string
}

func (r *Error) Error() string {
	return r.Str
}

func (r *Error) Code() uint16 {
	return r.Id
}


var idErrMap = map[uint16]error{}
var errMapLock sync.RWMutex
var errIdMap = map[error]uint16{}
var errIdMapLock sync.RWMutex

func NewError(str string, id uint16) *Error {
	err := &Error{id, str}
	errMapLock.Lock()
	idErrMap[id] = err
	errMapLock.Unlock()

	errIdMapLock.Lock()
	errIdMap[err] = id
	errIdMapLock.Unlock()

	return err
}

var (
	Error_RpcExecErr = NewError("method exec panic",1)
	Error_RpcNotFind = NewError("method not find",2)
	Error_RpcRespNil = NewError("result and error are nil",3)
	Error_RpcUnknowErr = NewError("errorcode not find",4)
	Error_RpcCallErr = NewError("method can not call",5)
	Error_RpcRespType = NewError("requset type is err",6)

	Error_Unknown = NewError("unknow error", 255)
	Error_OK 	  = NewError("ok",256)
)

var MinUserError = 257

func GetError(id uint16) error {
	errMapLock.RLock()
	if e, ok := idErrMap[id]; ok {
		errMapLock.RUnlock()
		return e
	}
	errMapLock.RUnlock()
	return Error_Unknown
}

func GetErrId(err error) uint16 {
	errIdMapLock.RLock()
	defer errIdMapLock.RUnlock()
	if id, ok := errIdMap[err]; ok {
		return id
	}
	return errIdMap[Error_Unknown]
}
