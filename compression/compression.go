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

func CompressGzip(src []byte) ([]byte, error) {
	buf := &bytes.Buffer{}
	w, _ := gzip.NewWriterLevel(buf, gzip.BestCompression)

	leng, err := w.Write(src)
	if err != nil || leng == 0 {
		return nil, err
	}
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	b := buf.Bytes()

	return b, nil
}
func DecompressGzip(src []byte) ([]byte, error) {
	r := bytes.NewReader(src)
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	rb, err := ioutil.ReadAll(gr)
	if err != nil {
		return nil, err
	}

	gr.Close()

	return rb, nil
}

func Compresszlib(src []byte) ([]byte, error) {
	buf := &bytes.Buffer{}
	w, _ := zlib.NewWriterLevel(buf, zlib.BestCompression)

	leng, err := w.Write(src)
	if err != nil || leng == 0 {
		return nil, err
	}
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	b := buf.Bytes()

	return b, nil
}
func Decompresszlib(src []byte) ([]byte, error) {
	r := bytes.NewReader(src)
	gr, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}
	rb, err := ioutil.ReadAll(gr)
	if err != nil {
		return nil, err
	}

	gr.Close()

	return rb, nil
}
