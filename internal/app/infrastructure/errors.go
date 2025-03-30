package infrastructure

import "errors"

var (
	ErrChannelNotFound = errors.New("channel not found")
	ErrImportNotFound  = errors.New("import not found")
)
