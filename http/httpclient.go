package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client http client
type Client struct {
	*http.Client
}

// NewHTTPClient 创建一个http client
func NewHTTPClient(timeout time.Duration) *Client {
	c := Client{
		&http.Client{
			Timeout: timeout,
		},
	}
	return &c
}

// Do 执行http请求
func (c *Client) Do(method, url string, req interface{}, rsp interface{}, opts ...Option) error {
	qopts := options{
		serialization: JSON,
	}
	for _, o := range opts {
		o(&qopts)
	}

	// marshaler
	m, ok := Marshalers[qopts.serialization]
	if !ok {
		return errors.New("invalid serialization")
	}
	buf, err := m.Marshal(req)
	if err != nil {
		return err
	}

	// http request
	httpReq, err := http.NewRequest(method, url, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	// headers
	switch qopts.serialization {
	case JSON:
		httpReq.Header.Add("Content-Type", "application/json")
	case FORM:
		httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	default:
		return errors.New("invalid serialization type")
	}

	for k, v := range qopts.headers {
		httpReq.Header.Add(k, v)
	}

	// send/recv
	httpRsp, err := c.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRsp.Body.Close()

	dat, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return err
	}

	// check error
	if httpRsp.StatusCode != http.StatusOK {
		return fmt.Errorf("http statusCode:%d, status:%s, details:%s",
			httpRsp.StatusCode, httpRsp.Status, string(dat))
	}

	// unmarshal
	err = json.Unmarshal(dat, rsp)
	if err != nil {
		return err
	}
	return nil
}
