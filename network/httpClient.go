package network

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ntfox0001/svrLib/log"
)

const (
	ContentTypeText = "text/plain"
	ContentTypeJson = "application/json"
	ContentTypeFrom = "application/x-www-form-urlencoded"
	ContentTypeFile = "multipart/form-data"
)

func SyncHttpGet(url string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 10}
	resp, err := client.Get(url)

	if err != nil {
		log.Error("syncHttpGet", "url", url, "error", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil

}

func SyncHttpPost(url string, content string, contentType string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 10}
	resp, err := client.Post(url, contentType, strings.NewReader(content))

	if err != nil {
		log.Error("SyncHttpPost", "url", url, "error", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func SyncHttpPostByHeader(url string, content string, contentType string, header map[string]string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 10}

	req, err := http.NewRequest("POST", url, strings.NewReader(content))
	if err != nil {
		log.Error("SyncHttpPostByHeader", "url", url, "error", err.Error())
		return "", err
	}

	req.Header.Set("Content-Type", contentType)

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("SyncHttpPostByHeader", "url", url, "error", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
