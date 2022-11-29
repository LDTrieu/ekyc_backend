package wUtil

import (
	"context"
	"fmt"
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
