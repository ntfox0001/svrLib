package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/ntfox0001/svrLib/commonError"
	"strings"

	"github.com/ntfox0001/svrLib/log"
)

type KeyValePair struct {
	Key   string
	Value string
}
type KeyValePairArray []KeyValePair

func (p KeyValePairArray) Len() int {
	return len(p)
}
func (p KeyValePairArray) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p KeyValePairArray) Less(i, j int) bool {
	return p[i].Key < p[j].Key
}

type Xml struct {
	StringMap
}

func MakeWxSign(stru interface{}) (string, error) {
	var m map[string]string
	// 将结构转为map
	if err := I2Stru(stru, &m); err != nil {
		return "", err
	} else {
		if _, ok := m["key"]; !ok {
			return "", commonError.NewStringErr("need key")
		}
		nonceStr := RandString(m["nonce_str"])
		if len(nonceStr) > 32 {
			nonceStr = string([]byte(nonceStr)[:32])
		}
		m["nonce_str"] = nonceStr

		outStr := ""
		MapSortIter_ss(m, func(k string, v string) {
			if v == "" {
				return
			}
			if k == "key" {
				return
			}
			outStr = fmt.Sprintf("%s&%s=%s", outStr, k, v)
		})

		outStr = fmt.Sprintf("%s&key=%s", outStr, m["key"])
		outStr = string([]byte(outStr)[1:])

		h := md5.New()
		h.Write([]byte(outStr))
		s := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
		//fmt.Println(s)
		m["sign"] = s
		delete(m, "key")
		x := Xml{}
		x.StringMap = m
		if xmlStr, err := xml.Marshal(x); err != nil {
			return "", err
		} else {
			return string(xmlStr), nil
		}
	}
}

func ValidateWxSign(xmlString string, key string) bool {
	x := Xml{}
	if err := xml.Unmarshal([]byte(xmlString), &x); err != nil {
		return false
	} else {
		if _, ok := x.StringMap["sign"]; !ok {
			log.Warn("sign does not exist.")
			return false
		}
		sign := x.StringMap["sign"]
		outStr := ""
		MapSortIter_ss(x.StringMap, func(k string, v string) {
			if v == "" {
				return
			}
			if k == "key" {
				return
			}
			if k == "sign" {
				return
			}
			outStr = fmt.Sprintf("%s&%s=%s", outStr, k, v)
		})

		outStr = fmt.Sprintf("%s&key=%s", outStr, key)
		outStr = string([]byte(outStr)[1:])

		fmt.Println(outStr)
		h := md5.New()
		h.Write([]byte(outStr))
		s := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))

		return s == sign
	}

}

//微信JS分享签名
func MakeWxJsShareSign(stru interface{}) (string, error) {
	var m map[string]string
	// 将结构转为map
	if err := I2Stru(stru, &m); err != nil {
		return "", err
	} else {
		// nonceStr := RandString(m["noncestr"])
		// if len(nonceStr) > 32 {
		// 	nonceStr = string([]byte(nonceStr)[:32])
		// }
		// m["noncestr"] = nonceStr

		outStr := ""
		MapSortIter_ss(m, func(k string, v string) {
			if v == "" {
				return
			}
			outStr = fmt.Sprintf("%s&%s=%s", outStr, k, v)
		})

		outStr = string([]byte(outStr)[1:])
		fmt.Printf(outStr)
		h := sha1.New()
		h.Write([]byte(outStr))
		s := hex.EncodeToString(h.Sum(nil))
		return s, nil
	}
}
