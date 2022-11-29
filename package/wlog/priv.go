package wlog

import (
	"fmt"
	"runtime"
	"strings"
)

func getFileLine() (file string, line int) {
	_, file, line, _ = runtime.Caller(2)
	para := strings.Split(file, "/")
	size := len(para)
	if size > 2 {
		file = fmt.Sprintf("%v/%v", para[size-2], para[size-1])
	}
	return
}
