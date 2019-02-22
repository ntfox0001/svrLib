package util

import (
	"crypto/md5"
	"encoding/hex"
	"sort"
)

// map排序
func MapSortIter_si(m map[string]interface{}, f func(k string, v interface{})) {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		f(k, m[k])
	}
}

func MapSortIter_ss(m map[string]string, f func(k string, v string)) {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		f(k, m[k])
	}
}

// 返回一个字符串map的md5
func Map2MD5(m map[string]string) string {
	b := Map2Bytes(m)
	md5 := md5.New()
	md5.Write(b)
	md5str := hex.EncodeToString(md5.Sum(nil))

	return md5str
}

// 返回字符Map的bytes
func Map2Bytes(m map[string]string) []byte {
	b := make([]byte, 0, 256)

	// 排序
	MapSortIter_ss(m, func(k string, v string) {
		b = append(b, []byte(k)...)
		b = append(b, []byte(v)...)
	})
	return b
}

// 转化一个数组map
func ArrayMap2MD5(am []map[string]string) string {
	b := make([]byte, 0, 256)
	for _, m := range am {
		mb := Map2Bytes(m)
		b = append(b, mb...)
	}

	md5 := md5.New()
	md5.Write(b)
	md5str := hex.EncodeToString(md5.Sum(nil))
	return md5str
}
