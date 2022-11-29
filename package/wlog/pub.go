package wlog

import (
	"context"
	"ekyc-app/source/wUtil"
	"encoding/json"
	"fmt"
	"log"
)

func Info(ctx context.Context, a ...interface{}) {
	txt := fmt.Sprint(a...)
	file, line := getFileLine()
	obj := map[string]interface{}{
		"file":  file,
		"line":  line,
		"info":  txt,
		"reqId": wUtil.GetReqId(ctx),
		//	"hardwareId": wUtil.GetHarwareId(ctx),
		"ip":    wUtil.GetClientIp(ctx),
		"level": "info",
	}

	jBuff, _ := json.Marshal(obj)
	log.Println(string(jBuff))
}
