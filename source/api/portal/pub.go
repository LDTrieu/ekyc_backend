package portal

import (
	"ekyc-app/library/qrcode"
	"ekyc-app/package/wlog"
	"ekyc-app/source/fsdb"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Reg(router *gin.Engine) {
	router.GET("/ws/connect", websocketConnection)
	router.GET("/ping", ping)
	router.GET("/login/qr/download/:qrId/:reqId", downloadQR)
	router.GET("/login/qr/rend/:reqId", rendQRLogin)
	router.POST("/login/auth/:reqId", loginBasic)
	router.POST("/signup/auth/:reqId", signupBasic)

	router.GET("portal/filterListUser/:reqId", filterListUser)
	router.GET("portal/filterListStudent/:reqId", filterListStudent)
	router.POST("portal/createStudentProfile/:reqId", createStudentProfile)
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "OKE",
	})
}

func rendQRLogin(c *gin.Context) {
	var (
		request = rendQRLoginRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
		}
	)
	response, err := guestRendQRLogin(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}
	// Trace client and result
	response.traceField = request.traceField
	c.JSON(http.StatusOK, response)
}

func downloadQR(c *gin.Context) {
	var (
		qrdata []byte
	)
	switch c.Query("action") {
	case qrActionLoginWebPortal:
		token, ok, err := fsdb.LoginSessionDBC.
			GetToken(c, c.Param("qrId"))
		if err != nil || !ok {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		qrdata = qrcode.GenQRCode(token)

	default:
		c.AbortWithError(http.StatusNotFound,
			errors.New("action undefine"))
		return
	}
	//WRITE PHOTO
	c.Writer.Header().Set("Content-Type",
		http.DetectContentType(qrdata))
	if _, err := c.Writer.Write(qrdata); err != nil {
		c.AbortWithError(http.StatusServiceUnavailable, err)
		return
	}
}

/* */
func loginBasic(c *gin.Context) {
	var (
		request = loginBasicRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
		}
	)
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := __loginBasic(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)

}

/* */
func signupBasic(c *gin.Context) {
	var (
		request = signupBasicRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
		}
	)
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := __signupBasic(c.Request.Context(),
		&request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)

}

/* */
func filterListUser(c *gin.Context) {
	// validate token
	status, _, from, err := validateBearer(c.Request.Context(), c.Request)
	if err != nil {
		c.AbortWithError(status, err)
		return
	}
	var (
		request = filterListUserRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			Permit: from,
		}
	)
	// if err := c.BindJSON(&request.Payload); err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, err)
	// 	return
	// }
	resp, err := __filterListUser(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)

}

/* */
func filterListStudent(c *gin.Context) {
	// validate token
	status, _, from, err := validateBearer(c.Request.Context(), c.Request)
	if err != nil {
		c.AbortWithError(status, err)
		return
	}
	var (
		request = filterListStudentRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			Permit: from,
		}
	)
	// if err := c.BindJSON(&request.Payload); err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, err)
	// 	return
	// }
	resp, err := __filterListStudent(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)

}

/* */
func createStudentProfile(c *gin.Context) {
	// validate token
	status, _, from, err := validateBearer(c.Request.Context(), c.Request)
	if err != nil {
		c.AbortWithError(status, err)
		return
	}
	var (
		request = createStudentProfileRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			Permit: from,
		}
	)
	if err := c.BindJSON(&request.Payload); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := __createStudentProfile(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)
}
