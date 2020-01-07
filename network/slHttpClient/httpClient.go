package slHttpClient

import (
	"github.com/ntfox0001/svrLib/goroutinePool"
	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"

	"github.com/ntfox0001/svrLib/log"
)

var _self *HttpClientManager

type HttpClientManager struct {
	goPool *goroutinePool.GoroutinePool
}

type HttpClientResult struct {
	Body string
	Err  error
}

func Instance() *HttpClientManager {
	if _self == nil {
		_self = &HttpClientManager{}
	}
	return _self
}

func (*HttpClientManager) Initial(goPoolSize, execSize int) {
	_self.goPool = goroutinePool.NewGoPool("HttpClientManager", goPoolSize, execSize)
}
func (*HttpClientManager) Release() {
	_self.goPool.Release(0)

	log.Debug("HttpClientManager release")
}

func (*HttpClientManager) HttpGet(cb *selectCaseInterface.CallbackHandler, url string) {
	hp := func(data interface{}) {
		rtStr, err := network.SyncHttpGet(url)
		rt := HttpClientResult{
			Body: rtStr,
			Err:  err,
		}
		if cb != nil {
			cb.SendReturnMsgNoReturn(rt)
		}

	}
	_self.goPool.Go(hp)
}

func (*HttpClientManager) HttpPost(cb *selectCaseInterface.CallbackHandler, url string, content string, contentType string) {
	hp := func(data interface{}) {
		rtStr, err := network.SyncHttpPost(url, content, contentType)
		rt := HttpClientResult{
			Body: rtStr,
			Err:  err,
		}
		if cb != nil {
			cb.SendReturnMsgNoReturn(rt)
		}

	}
	_self.goPool.Go(hp)
}

func (*HttpClientManager) HttpPostByHeader(cb *selectCaseInterface.CallbackHandler, url string, content string, contentType string, header map[string]string) {
	hp := func(data interface{}) {
		rtStr, err := network.SyncHttpPostByHeader(url, content, contentType, header)
		rt := HttpClientResult{
			Body: rtStr,
			Err:  err,
		}
		if cb != nil {
			cb.SendReturnMsgNoReturn(rt)
		}

	}
	_self.goPool.Go(hp)
}
