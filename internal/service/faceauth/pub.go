package faceauth

import (
	"context"
	"ekyc-app/library/net"
	"ekyc-app/library/net/options"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	host                    = "http://localhost:8081"
	portal_request_endpoint = "/mobile/auth/ping-third-party"
	face_auth_endpoint      = "/mobile/auth/face-auth-session"
	timeout                 = 30 * time.Second
)

func RequestSession(ctx context.Context,
	mock_data *MockModel) (int, *FaceAuthSessionResponse, error) {
	var (
		url = func() string {
			if len(host) < 9 /* 127.0.0.1 */ {
				host = "http://127.0.0.1:8080"
			}
			if !strings.Contains(host, "http") {
				return fmt.Sprintf("http://%s%s", host, portal_request_endpoint)
			} else {
				return fmt.Sprintf("%s%s", host, portal_request_endpoint)
			}
		}()
		opt = options.CurlOption()

		// request = struct {
		// 	Name string `json:"name"`
		// 	Code int    `json:"code"`
		// }{
		// 	Name: mock_data.Name,
		// 	Code: mock_data.Code,
		// }

		result = struct {
			Code    int                     `json:"code"`    // Mã lỗi. Code 0: Thành công, Code != 0 Thất bại
			Message string                  `json:"message"` // Nội dung thông báo lỗi
			Payload FaceAuthSessionResponse `json:"payload"`
		}{}
	)
	opt.SetMethod(http.MethodPost)
	opt.SetTimeout(timeout)
	opt.AddHeader("Content-Type", "application/json")
	opt.SetJSON(mock_data)

	response, err := net.Curl(url, opt)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	if response.StatusCode != http.StatusOK &&
		response.StatusCode > http.StatusNoContent {
		return response.StatusCode, nil, fmt.Errorf("reponse bad from %s with status %d", url, response.StatusCode)
	}
	if err = json.Unmarshal(response.Body, &result); err != nil {
		return http.StatusForbidden, nil, err
	}
	return mock_data.Code, &result.Payload, nil
}

func FaceAuthSession(ctx context.Context,
	face_auth_data *options.FormFile) (int, *FaceAuthSessionResponse, error) {
	var (
		url = func() string {
			if len(host) < 9 /* 127.0.0.1 */ {
				host = "http://127.0.0.1:8080"
			}
			if !strings.Contains(host, "http") {
				return fmt.Sprintf("http://%s%s", host, face_auth_endpoint)
			} else {
				return fmt.Sprintf("%s%s", host, face_auth_endpoint)
			}
		}()
		opt  = options.CurlOption()
		form = make(map[string]interface{})

		// request = struct {
		// 	FileName string       `json:"filename"`
		// 	File     bytes.Buffer `json:"file"`
		// }{}

		result = struct {
			Code    int                     `json:"code"`    // Mã lỗi. Code 0: Thành công, Code != 0 Thất bại
			Message string                  `json:"message"` // Nội dung thông báo lỗi
			Payload FaceAuthSessionResponse `json:"payload"`
		}{}
	)

	opt.SetMethod(http.MethodPost)
	opt.SetTimeout(timeout)
	//opt.AddHeader("Content-Type", writer.FormDataContentType())
	//opt.SetJSON(mock_data)
	form["filename"] = face_auth_data
	if err := opt.SetFormData(form); err != nil {
		return http.StatusBadRequest, nil, err
	}

	response, err := net.Curl(url, opt)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if response.StatusCode != http.StatusOK &&
		response.StatusCode > http.StatusNoContent {
		return response.StatusCode, nil, fmt.Errorf("reponse bad from %s with status %d", url, response.StatusCode)
	}

	if err = json.Unmarshal(response.Body, &result); err != nil {
		return http.StatusForbidden, nil, err
	}
	return result.Code, &result.Payload, nil
}
