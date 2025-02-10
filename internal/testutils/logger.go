package testutils

import (
	"bytes"
	"log/slog"
	"strings"
)

func NewLogger() (*bytes.Buffer, *slog.Logger) {
	buf := new(bytes.Buffer)
	return buf, slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{}))
}

func LogLines(buf *bytes.Buffer) []string {
	return strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n")
}
