package util

type Cache struct {
	buff     []interface{}
	signBuff []byte

	size int

	// 最低的id，指向最低的那个没用的id
	lowId int

	newFunc    func() interface{}
	returnFunc func(interface{}) interface{}
}

// cache指定的对象，对象是否可重用由调用者维护，size是cache大小，如果多了会自动添加，但是不会变小
// newfunc用于创建新对象
func NewCache(size int, newFunc func() interface{}, returnFunc func(interface{}) interface{}) *Cache {
	cache := &Cache{
		size:       size,
		buff:       make([]interface{}, size, size),
		signBuff:   make([]byte, size, size),
		lowId:      0,
		newFunc:    newFunc,
		returnFunc: returnFunc,
	}

	return cache
}
func (c *Cache) Get(i int) interface{} {
	if i < len(c.buff) {
		if c.signBuff[i] == 0 {
			return nil
		}
		return c.buff[i]
	}
	return nil
}

func (c *Cache) GetAvailableId() int {
	for i := c.lowId; i < len(c.signBuff); i++ {
		if c.signBuff[i] == 0 {
			c.signBuff[i] = 1
			if c.buff[i] == nil {
				c.buff[i] = c.newFunc()
			}
			c.lowId = i
			return i
		}
	}
	// 如果都用上了，那么创建新的
	elem := c.newFunc()
	c.buff = append(c.buff, elem)
	c.signBuff = append(c.signBuff, 1)

	return len(c.signBuff) - 1
}

func (c *Cache) Return(i int) bool {
	if i >= len(c.signBuff) {
		return false
	}

	if c.signBuff[i] == 0 {
		return false
	}

	if i < c.lowId {
		c.lowId = i
	}

	c.signBuff[i] = 0
	c.buff[i] = c.returnFunc(c.buff[i])
	return true
}
