package mobile

import (
	"ekyc-app/package/wlog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Reg(router *gin.Engine) {
	router.POST("/login/submitQRLogin/:reqId", submitQRLogin)

	// Sign in
	router.POST("/mobile/auth/login/:reqId", loginTerminal)
	router.POST("/mobile/auth/signup/:reqId", signupTerminal)

	// Image
	router.POST("/mobile/auth/aten/:reqId", loginTerminal)
	
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
