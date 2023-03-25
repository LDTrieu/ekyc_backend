package net

import (
	"ekyc-app/library/net/options"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

/*
Curl function.
By default, if the option is NIL then timeout is 30 seconds.
*/
func Curl(url string, option *options.Option) (*Response, error) {
	if option == nil {
		option = options.CurlOption()
	}
	if strings.HasPrefix(url, "ws") {
		return nil, errors.New("the `websocket` protocol is not supported")
	} else {
		url = UrlPrettyParse(url)
	}
	if len(option.Params) > 0 {
		url += "/" + option.Params
	}
	request, err := http.NewRequest(option.Method, url, option.Body)
	if err != nil {
		return nil, err
	}
	for key, arr := range option.Header {
		for i, value := range arr {
			if i == 0 {
				request.Header.Set(key, value)
			} else {
				request.Header.Add(key, value)
			}
		}
	}
	request.Header.Set(
		ekycTime,
		fmt.Sprintf("%d", time.Now().Unix()),
	)
	client := &http.Client{
		Timeout: option.Timeout,
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	buffer, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	return &Response{
		StatusCode: response.StatusCode,
		Header:     response.Header.Clone(),
		Body:       buffer,
	}, err
}
