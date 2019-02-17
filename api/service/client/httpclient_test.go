package client

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"
	"time"
)

var HTTP_URL = "http://mockbin.com/request?foo=bar&foo=baz"

func Test_HttpClientGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), HTTP_DEFAULT_TIMEOUT*time.Second)
	defer cancel()
	httpclient := NewHTTPClient("", "", HTTP_DEFAULT_TIMEOUT, true, "")
	req, err := httpclient.NewRequest(ctx, "GET", HTTP_URL, nil, "")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	resp, err := httpclient.Do(req)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	respbyte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(string(respbyte))
	t.Log("success HttpClientGet")
}

func Test_HttpClientPost(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), HTTP_DEFAULT_TIMEOUT*time.Second)
	defer cancel()
	httpclient := NewHTTPClient("", "", HTTP_DEFAULT_TIMEOUT, true, "")
	param := url.Values{}
	param.Set("foo", "bar")
	param.Add("foo", "baz")
	req, err := httpclient.NewRequest(ctx, "POST", HTTP_URL, strings.NewReader(param.Encode()), MIMEApplicationForm)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	resp, err := httpclient.Do(req)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	respbyte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(string(respbyte))
	t.Log("success HttpClientPost")
}

func Test_HttpClientPostJson(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), HTTP_DEFAULT_TIMEOUT*time.Second)
	defer cancel()
	httpclient := NewHTTPClient("", "", HTTP_DEFAULT_TIMEOUT, true, "")
	param := bytes.NewBuffer([]byte("{\"foo\": [\"bar\", \"baz\"]}"))
	req, err := httpclient.NewRequest(ctx, "POST", HTTP_URL, param, MIMEApplicationJSON)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	resp, err := httpclient.Do(req)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	respbyte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(string(respbyte))
	t.Log("success HttpClientPostJson")
}
