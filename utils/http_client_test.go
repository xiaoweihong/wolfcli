package utils

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	req := Get("http://httpbin.org/get")
	resp, err := req.Response()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
	t.Log(req.String())
}

func TestBeegoHTTPRequest_DoRequest(t *testing.T) {
	req := Get("https://goolnk.com/33BD2j")
	retryAmount := 1
	req.Retries(1)
	req.RetryDelay(1400 * time.Millisecond)
	retryDelay := 1400 * time.Millisecond

	req.setting.CheckRedirect = func(redirectReq *http.Request, redirectVia []*http.Request) error {
		return errors.New("Redirect triggered")
	}

	startTime := time.Now().UnixNano() / int64(time.Millisecond)

	_, err := req.Response()
	if err == nil {
		t.Fatal("Response should have yielded an error")
	}

	endTime := time.Now().UnixNano() / int64(time.Millisecond)
	elapsedTime := endTime - startTime
	delayedTime := int64(retryAmount) * retryDelay.Milliseconds()

	if elapsedTime < delayedTime {
		t.Errorf("Not enough retries. Took %dms. Delay was meant to take %dms", elapsedTime, delayedTime)
	}
}

func TestSimplePost(t *testing.T) {
	v := "smallfish"
	req := Post("http://httpbin.org/post")
	req.Body("xiaowei")

	str, err := req.String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)

	n := strings.Index(str, v)
	if n == -1 {
		t.Fatal(v + " not found in post")
	}
}