package pool_test

import (
	"fmt"
	"testing"

	"github.com/ntfox0001/svrLib/pool"
)

type TestRes struct {
	A int
}

func (t TestRes) Close() {
	t.A = 0
}
func Test1(t *testing.T) {
	pool := pool.NewPool(func() (pool.Resource, error) {
		fmt.Println("new")
		return TestRes{}, nil
	}, 10)

	for i := 0; i < 100; i++ {
		res, _ := pool.Get()
		pool.Put(res)
	}

}
