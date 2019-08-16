package compression

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io/ioutil"
)

type Compression struct {
	// writerPool sync.Pool
	// readerPool sync.Pool
	// buffPool   sync.Pool
}

// func NewCompression() *Compression {
// 	c := Compression{
// 		writerPool: sync.Pool{
// 			New: func() interface{} {
// 				buf := new(bytes.Buffer)
// 				return gzip.NewWriter(buf)
// 			},
// 		},
// 		readerPool: sync.Pool{
// 			New: func() interface{} {
// 				r := bytes.NewReader(nil)
// 				return gzip.NewReader(r)
// 			},
// 		},
// 		buffPool: sync.Pool{
// 			New: func() interface{} {
// 				return new(bytes.Buffer)
// 			},
// 		},
// 	}
// }

func CompressGzip(src []byte) []byte {
	buf := &bytes.Buffer{}
	w, _ := gzip.NewWriterLevel(buf, gzip.BestCompression)

	leng, err := w.Write(src)
	if err != nil || leng == 0 {
		return nil
	}
	err = w.Flush()
	if err != nil {
		return nil
	}
	err = w.Close()
	if err != nil {
		return nil
	}
	b := buf.Bytes()

	return b
}
func DecompressGzip(src []byte) []byte {
	r := bytes.NewReader(src)
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil
	}
	rb, err := ioutil.ReadAll(gr)
	if err != nil {
		return nil
	}

	gr.Close()

	return rb
}

func Compresszlib(src []byte) []byte {
	buf := &bytes.Buffer{}
	w, _ := zlib.NewWriterLevel(buf, zlib.BestCompression)

	leng, err := w.Write(src)
	if err != nil || leng == 0 {
		return nil
	}
	err = w.Flush()
	if err != nil {
		return nil
	}
	err = w.Close()
	if err != nil {
		return nil
	}
	b := buf.Bytes()

	return b
}
func Decompresszlib(src []byte) []byte {
	r := bytes.NewReader(src)
	gr, err := zlib.NewReader(r)
	if err != nil {
		return nil
	}
	rb, err := ioutil.ReadAll(gr)
	if err != nil {
		return nil
	}

	gr.Close()

	return rb
}
