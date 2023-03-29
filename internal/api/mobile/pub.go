package mobile

import (
	"ekyc-app/internal/model"
	"ekyc-app/library/net"
	"ekyc-app/package/wlog"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Reg(router *gin.Engine) {
	router.POST("/login/submitQRLogin/:reqId", submitQRLogin)

	// Sign in
	router.POST("/mobile/auth/login/:reqId", loginTerminal)
	router.POST("/mobile/auth/signup/:reqId", signupTerminal)

	// Image
	router.POST("/mobile/auth/face-auth-session/:reqId", faceAuthSession)
	router.GET("/mobile/auth/ping-third-party/:reqId", pingThirdParty)
	// Push notify

}

// Device send a JSON to Server
func submitQRLogin(c *gin.Context) {
	var (
		request = submitQRLoginRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
		}
	)
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := guestSubmitQRLogin(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}
	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)
}

// /* */
func loginTerminal(c *gin.Context) {
	var (
		request = loginTerminalRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
		}
	)
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := __loginTerminal(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)
}

/* */
func signupTerminal(c *gin.Context) {
	// validate token
	status, _, from, err := validateBearer(c.Request.Context(), c.Request)
	if err != nil {
		c.AbortWithError(status, err)
		return
	}
	var (
		request = signupTerminalRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			Permit: from,
		}
	)

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := __sigupTerminal(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)
}

/* */
func pingThirdParty(c *gin.Context) {
	// validate token
	// status, _, auth_data, err := validateBearer(c.Request.Context(), c.Request)
	// if err != nil {
	// 	c.AbortWithError(status, err)
	// 	return
	// }
	//log.Println("auth_data.AccountID: ", auth_data.AccountID)
	var (
		request = pingThirdPartyRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			// Payload: face_image_req{
			// 	TerminalId: auth_data.AccountID,
			// 	StudentId:  c.Param("studentId"),
			// 	FileName:   fmt.Sprintf("%s_%s.bin", c.Param("studentId"), primitive.NewObjectID().Hex()),
			// },
		}
	)

	resp, err := __pingThirdParty(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}
	// Trace client and result
	//resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)
}

/* */
func faceAuthSession(c *gin.Context) {
	// validate token
	status, _, auth_data, err := validateBearer(c.Request.Context(), c.Request)
	if err != nil {
		c.AbortWithError(status, err)
		return
	}

	var (
		request = faceAuthSessionRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			Payload: face_image_req{
				TerminalId: auth_data.AccountID,
				// StudentId:  "studentId",
				FileName: fmt.Sprintf("%s_%s.bin", "StudentId", primitive.NewObjectID().Hex()),
			},
		}
	)

	// validate image
	// Read multipart form files
	multipart_form, err := net.NewMultipartForm(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	filename, err := multipart_form.GetForm("filename")
	if err != nil {
		wlog.Error(c, err)
	}

	file_name, file, err := multipart_form.GetFile("filename")
	if err != nil {
		wlog.Error(c, err)
	}
	if len(file) < 1 {
		c.JSON(http.StatusBadRequest, faceAuthSessionResponse{
			Code:       model.StatusFileNotFound,
			Message:    "FILE_NOT_FOUND",
			traceField: request.traceField,
		},
		)
		return
	}
	request.Payload.FileName = func() string {
		if len(filename) > 0 {
			return filename
		}
		return file_name
	}()
	request.Payload.File = file
	// call API to Django
	resp, err := __faceAuthSession(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
		c.AbortWithError(http.StatusServiceUnavailable, err)
		return
	}

	// Trace client and result
	//resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)
}
