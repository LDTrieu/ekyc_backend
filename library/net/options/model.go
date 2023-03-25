package options

import (
	"io"
	"net/http"
	"time"
)

/*
Option object
*/
type Option struct {
	Method  string // http.MethodX
	Params  string
	Body    io.Reader
	Header  http.Header
	Timeout time.Duration
}

type FormFile struct {
	Filename string
	File     []byte
}
