package queue

import "errors"

var ErrRetriable = errors.New("retriable error")
var ErrFatal = errors.New("fatal error")
