package options

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"errors"
)

/*
CurlOption function
*/
func CurlOption() *Option {
	return &Option{
		Method:  http.MethodGet,
		Header:  http.Header{},
		Timeout: 180 * time.Second,
	}
}

/*
SetMethod function
*/
func (ins *Option) SetMethod(method string) error {
	var temp = map[string]bool{
		http.MethodConnect: true,
		http.MethodDelete:  true,
		http.MethodGet:     true,
		http.MethodHead:    true,
		http.MethodOptions: true,
		http.MethodPatch:   true,
		http.MethodPost:    true,
		http.MethodPut:     true,
		http.MethodTrace:   true,
	}

	value, ok := temp[method]
	if ok && value {
		ins.Method = method
		return nil
	}

	return errors.New("method invalid")
}

/*
SetHeader function
*/
func (ins *Option) SetHeader(header map[string][]string) {
	for key, h := range header {
		for i, value := range h {
			if i == 0 {
				ins.Header.Set(key, value)
			} else {
				ins.Header.Add(key, value)
			}
		}
	}
}

/*
AddHeader function
*/
func (ins *Option) AddHeader(key, value string) {
	ins.Header.Add(key, value)
}

/*
SetBasicAuth function
*/
func (ins *Option) SetBasicAuth(username, password string) {
	basicAuth := func(username, password string) string {
		auth := username + ":" + password
		return base64.StdEncoding.EncodeToString([]byte(auth))
	}
	ins.Header.Set("Authorization", "Basic "+basicAuth(username, password))
}

/*
SetData function
*/
func (ins *Option) SetData(body []byte) {
	n := make([]byte, len(body))
	copy(n, body)
	ins.Body = bytes.NewReader(n)
}

/*
SetJSON function
*/
func (ins *Option) SetJSON(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	} else {
		ins.Header.Set("Content-Type", "application/json")
	}
	ins.Body = bytes.NewReader(b)
	return nil
}

func (ins *Option) SetFormData(formdata map[string]interface{}) error {
	var (
		payload     = bytes.NewBuffer(nil)
		writer      = multipart.NewWriter(payload)
		field_count = 0
	)
	defer writer.Close()
	// for key, form := range files {
	// 	w, err := writer.CreateFormFile(key, form.Filename)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	if _, err := w.Write(form.File); err != nil {
	// 		continue
	// 	}
	// }

	for key, value := range formdata {
		switch data := value.(type) {
		case string:
			println("Type: String", "; Key:", key, "; Value:", data)
			w, err := writer.CreateFormField(key)
			if err != nil {
				continue
			}
			if _, err := w.Write(bytes.NewBufferString(data).Bytes()); err != nil {
				continue
			}
			field_count++
		case *FormFile:
			println("Type: *FormFile", "; Key:", key, "; Value:", fmt.Sprintf("%+v", data))
			w, err := writer.CreateFormFile(key, data.Filename)
			if err != nil {
				continue
			}
			if _, err := w.Write(data.File); err != nil {
				continue
			}
			field_count++
		default:
			println("Type: Unknow", "; Key:", key, "; Value:", fmt.Sprintf("%+v", data))
		}
	}
	if field_count > 0 {
		ins.Header.Set("Content-Type", writer.FormDataContentType())
	}
	ins.Body = payload
	return nil
}

/*
SetTimeout function
*/
func (ins *Option) SetTimeout(d time.Duration) {
	ins.Timeout = d
}

/*
SetParam function.
*/
func (ins *Option) SetParam(param string) {
	ins.Params = param
}
