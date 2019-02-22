package network_test

import (
	"fmt"
	"github.com/ntfox0001/svrLib/network"
	"testing"
	"time"
)

func TestHttpClientGet(t *testing.T) {
	for i := 0; i < 2000; i++ {
		go testHttp(i)
	}
	// s1, err1 := network.SyncHttpGet("http://182.92.66.148/YingXiongZhanJi1128/server_list.xml")
	// s2, err2 := network.SyncHttpGet("https://www.guhuozaiol.com/YingXiongZhanJi1128/server_list.xml")

	// if err1 != nil || err2 != nil || s1 != s2 {
	// 	t.Fail()
	// }
	time.Sleep(time.Second * 10)
}

func testHttp(i int) {
	url := fmt.Sprintf("http://bit007.com.cn/test.php?test=%d", i)
	s1, _ := network.SyncHttpGet(url)
	fmt.Println(s1)
}

func TestHttpClientGet2(t *testing.T) {
	xml := `
	<?xml version="1.0"?>
<methodCall>
    <methodName>examples.getStateName</methodName>
    <params>
        <param>
            <value><i4>41</i4></value>
        </param>
    </params>
</methodCall>
	`
	s1, err1 := network.SyncHttpPostByHeader("http://www.guhuozaiol.com/php/common/wxLogin/test.php?t=11", xml, network.ContentTypeText, nil)
	s2, err2 := network.SyncHttpPostByHeader("https://www.guhuozaiol.com/php/common/wxLogin/test.php?t=11", xml, network.ContentTypeText, nil)
	fmt.Println(s1)
	fmt.Println(s2)
	if err1 != nil || err2 != nil || s1 != s2 {
		t.Fail()
	}
}

func TestHttpClientGet3(t *testing.T) {
	xml := `
	fdsgdsgsdfdsfsdgdasgdfagdfagdfagj kdls ajgidfj agkldfj ag
	df 
	aghdfj ahf
	a hgdf
	ah 
	fd
	a ghdfpag k

	`
	// xml := `
	// 	<?xml version="1.0"?>
	// <methodCall>
	//     <methodName>examples.getStateName</methodName>
	//     <params>
	//         <param>
	//             <value><i4>41</i4></value>
	//         </param>
	//     </params>
	// </methodCall>
	// 	`
	s1, err1 := network.SyncHttpPost("http://www.guhuozaiol.com/php/common/wxLogin/test.php?t=111", xml, network.ContentTypeText)
	s2, err2 := network.SyncHttpPost("https://www.guhuozaiol.com/php/common/wxLogin/test.php?t=111", xml, network.ContentTypeText)

	fmt.Println(s1)
	fmt.Println(s2)
	if err1 != nil || err2 != nil || s1 != s2 {
		t.Fail()
	}
}
