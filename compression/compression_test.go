package compression_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ntfox0001/svrLib/compression"
)

func Test1(t *testing.T) {
	buf1 := new(bytes.Buffer)
	buf2 := bytes.Buffer{}
	var buf3 bytes.Buffer
	var buf4 *bytes.Buffer

	fmt.Println(buf1, buf2, buf3, buf4)
}

func Test2(t *testing.T) {
	c := compression.Compression{}

	s := []byte("afdsafdsajgieroqjgi4weqjgi4o3qjgjlkdfasjhogfdjsahklfdsjgkr'e jytio45bw3qj6y950vw43u6t930 wau6t9450w3 y6tu59bgy4y0g645bh09 iuj465w0y g9vmuj9wabh6tjmu45wesh904y65nb uy90465twmby u54mb 54nbw y;u654w9y 8o;56bmw4uy95o;m4wsub jyitrosm;ub98y5mo;b4u9y5roesbm;uy98rom;bsu8y;obm uy590ws;b miuy954e;swmb uyio65j4qiv6 4j3qv 7lj45bq7 j59043pq ujy658 puqa4t8j459p3qvut48903qnpvutj4i3qov;tj43koqb64e3bfdsafsegrehreshreshrshrajiogjeriao;gu84oevut8ov4nu3qt9ov;4jqiwejtviero;qjgviroe;qjagieryaghuierahguierhalguekrhlajedfkahgu4ayvt849ayut8vio")
	fmt.Println(len(s))
	cb := c.CompressGzip(s)
	fmt.Println(len(cb))

	fmt.Println(string(c.DecompressGzip(cb)))

}
