package testutils

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type HTTPClientMock struct {
	mock.Mock
}

func NewHTTPClientMock() *HTTPClientMock {
	return &HTTPClientMock{}
}

func (c *HTTPClientMock) Do(req *http.Request) (*http.Response, error) {
	args := c.Called(req)
	if args.Get(1) == nil {
		return args.Get(0).(*http.Response), nil
	}

	return nil, args.Get(1).(error)
}
