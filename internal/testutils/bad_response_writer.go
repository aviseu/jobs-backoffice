package testutils

import (
	"errors"
	"net/http"
)

type BadResponseWriter struct {
	Code    int
	Headers http.Header
}

func NewBadResponseWriter() *BadResponseWriter {
	return &BadResponseWriter{
		Headers: http.Header{},
	}
}

func (w *BadResponseWriter) Header() http.Header {
	return w.Headers
}

func (*BadResponseWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("bad response writer")
}

func (w *BadResponseWriter) WriteHeader(status int) {
	w.Code = status
}
