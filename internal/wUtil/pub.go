package wUtil

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Error(err error) error {
	file, line := getFileLine()
	return fmt.Errorf("file:%v line:%v %v",
		file, line, err)
}

func NewCtx(reqId, ip, hwid string) context.Context {
	ctx := SetReqId(context.Background(), reqId)
	//ctx = SetHarwareId(ctx, hwid)
	ctx = SetClientIp(ctx, ip)
	return ctx
}

func SetReqId(ctx context.Context, reqId string) context.Context {
	return context.WithValue(ctx, keyReqId, reqId)
}

func GetReqId(ctx context.Context) string {
	return ctxValByKey(ctx, keyReqId)
}

func SetClientIp(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, keyClientIp, ip)
}

func GetClientIp(ctx context.Context) string {
	return ctxValByKey(ctx, keyClientIp)
}

func GetHost(c *gin.Context) string {
	return getTLS(c) + c.Request.Host
}

func GetUrl(c *gin.Context) (string, error) {

	tls := getTLS(c)
	requestURI := c.Request.RequestURI
	host := c.Request.Host
	hostURL := tls + host + requestURI
	return hostURL, nil
}

func getTLS(c *gin.Context) string {
	if c.Request.TLS != nil {
		return "https://"
	} else {
		return "http://"
	}
}
