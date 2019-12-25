package asyncRedis

import (
	"time"

	"github.com/ntfox0001/svrLib/goroutinePool"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"

	"github.com/go-redis/redis"
)

// 简单的使用多协程实现异步，没有使用pipe
type Client struct {
	RedisClient *redis.Client
	goPool      *goroutinePool.GoroutinePool
}

// 创建一个异步client
func NewClientByRedis(client *redis.Client, poolSize, execSize int) *Client {
	gp := goroutinePool.NewGoPool("RedisClient", poolSize, execSize)

	return &Client{
		RedisClient: client,
		goPool:      gp,
	}
}

// 创建一个异步client
func NewClient(address, password string, db, readTimeout, writeTimeout, minIdleConn, poolSize, execSize int) *Client {
	c := redis.NewClient(&redis.Options{
		Addr:         address,
		Password:     password,
		DB:           db,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		MinIdleConns: minIdleConn,
	})

	gp := goroutinePool.NewGoPool("RedisClient", poolSize, execSize)

	return &Client{RedisClient: c, goPool: gp}
}

func (c *Client) Close() {
	c.goPool.Release(0)
	c.RedisClient.Close()
}

// 异步执行命令
func (c *Client) AsyncCommond(f func(*redis.Client)) {
	c.goPool.Go(func(data interface{}) {
		f(c.RedisClient)
	}, nil)
}

// 调用异步执行，并返回数据
func (c *Client) AsyncCommondReturn(cb *selectCaseInterface.CallbackHandler, f func(*redis.Client) interface{}) {
	c.AsyncCommond(func(client *redis.Client) {
		data := f(client)
		cb.SendReturnMsgNoReturn(data)
	})
}
