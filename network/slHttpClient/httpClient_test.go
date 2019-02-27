package slHttpClient_test

import (
	"testing"

	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/network/slHttpClient"
	"github.com/ntfox0001/svrLib/selectCase"
)

func TestHttpClient1(t *testing.T) {

	sc := selectCase.NewSelectChannel()
	slHttpClient.Instance().Initial(10, 10)
	slHttpClient.Instance().HttpPost(nil, "http://www.guhuozaiol.com/php/common/wxLogin/test.php", "aa=aa", network.ContentTypeFrom)
	slHttpClient.Instance().HttpPost(nil, "https://www.guhuozaiol.com/php/common/wxLogin/test.php", "aa=aa", network.ContentTypeFrom)

	s1 := sc.GetReturn().(slHttpClient.HttpClientResult)
	s2 := sc.GetReturn().(slHttpClient.HttpClientResult)

	if s1.Body != s2.Body {
		t.Fail()
	}

}

func TestHttpClient2(t *testing.T) {

	sc := selectCase.NewSelectChannel()
	slHttpClient.Instance().Initial(10, 10)
	slHttpClient.Instance().HttpPostByHeader(nil, "http://www.guhuozaiol.com/php/common/wxLogin/test.php", "aa=aa", network.ContentTypeFrom, nil)
	slHttpClient.Instance().HttpPostByHeader(nil, "https://www.guhuozaiol.com/php/common/wxLogin/test.php", "aa=aa", network.ContentTypeFrom, nil)

	s1 := sc.GetReturn().(slHttpClient.HttpClientResult)
	s2 := sc.GetReturn().(slHttpClient.HttpClientResult)

	if s1.Body != s2.Body {
		t.Fail()
	}

}
