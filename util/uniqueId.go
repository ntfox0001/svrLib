package util

import (
	"crypto/sha256"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/oklog/ulid"
)

var id uint64 = 0
var entropy *rand.Rand

func init() {
	entropy = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// 获得一个本次运行以来的唯一值
func GetNextId() uint64 {
	// 返回一个唯一值
	return atomic.AddUint64(&id, 1)
}

// 获得一个全局唯一id
func GetUniqueId() string {

	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	return id.String()
}

// 继承这个类，可以得到这个类实例化以来的累进id
type UnqueIdWrap struct {
	_unqueId uint64
}

// 创建一个可以累进id的类
func NewUnqueIdWarp() UnqueIdWrap {
	return UnqueIdWrap{_unqueId: 0}
}

func (u *UnqueIdWrap) GetNextId() uint64 {
	return atomic.AddUint64(&u._unqueId, 1)
}

// 生成一个token
func NewToken(str string) string {
	hashstr := "uguess?215fuf!!hhaf.xXHgh4$%&2jsg" + GetUniqueId() + str

	h := sha256.New()
	h.Write([]byte(hashstr))
	rt := base58.CheckEncode(h.Sum(nil), 0)

	return rt
}
