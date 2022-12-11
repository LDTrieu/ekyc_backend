package mobile

import (
	"ekyc-app/package/wlog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Reg(router *gin.Engine) {
	router.POST("/login/submitQRLogin/:reqId", submitQRLogin)

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
