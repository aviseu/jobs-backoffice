package testutils

import (
	"io"
	"net/http"
	"strings"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestLog struct {
	Header   http.Header
	Method   string
	URL      string
	Body     string
	Response string
}

type RequestLogger struct {
	Client HTTPClient
	Logs   []*RequestLog
}

func NewRequestLogger(c HTTPClient) *RequestLogger {
	return &RequestLogger{
		Client: c,
		Logs:   make([]*RequestLog, 0),
	}
}

func (rl *RequestLogger) Do(req *http.Request) (*http.Response, error) {
	reqBuf := new(strings.Builder)
	if req.Body != nil {
		_, err := io.Copy(reqBuf, req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(strings.NewReader(reqBuf.String()))
	}

	l := &RequestLog{
		Header: req.Header,
		Method: req.Method,
		URL:    req.URL.String(),
		Body:   reqBuf.String(),
	}

	resp, err := rl.Client.Do(req)
	if err != nil {
		return nil, err
	}

	respBuf := new(strings.Builder)
	if resp.Body != nil {
		_, err = io.Copy(respBuf, resp.Body)
		if err != nil {
			return nil, err
		}
		resp.Body = io.NopCloser(strings.NewReader(respBuf.String()))

		l.Response = respBuf.String()
	}

	rl.Logs = append(rl.Logs, l)

	return resp, err
}
