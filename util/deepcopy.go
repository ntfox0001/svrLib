package util

import (
	"bytes"
	"encoding/gob"
)

// dstPrt必须是分配好内存的结构的指针，src必须是结构,只能拷贝public的成员变量
func DeepCopy(dstPtr, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dstPtr)
}
