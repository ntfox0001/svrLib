package util

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

// 字符串替换，根据map中的字符替换src中的字符
// 也可以使用go的，提供一个数组
func StringReplace(src string, replace map[string]string) string {
	patterns := make([]string, 0, len(replace))

	for k, v := range replace {
		patterns = append(patterns, k)
		patterns = append(patterns, v)
	}
	replacer := strings.NewReplacer(patterns...)

	return replacer.Replace(src)
}

// 产生随机字符串，长度不定
func RandString(seed string) string {
	b := make([]byte, 32)
	rand.Read(b)

	hashstr := fmt.Sprint(b, seed)

	h := sha256.New()
	h.Write([]byte(hashstr))
	rt := base58.CheckEncode(h.Sum(nil), 0)

	return rt
}

// 大写首字母
func UpperFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	s1 := string([]byte(s)[0:1])
	s2 := string([]byte(s)[1:])
	s3 := fmt.Sprint(strings.ToUpper(s1), s2)
	return s3
}
