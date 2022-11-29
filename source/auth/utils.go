package auth

import (
	"context"
	"ekyc-app/gcloud"
	"ekyc-app/source/cfg"
	"ekyc-app/source/wUtil"
	"encoding/base64"
	"time"
)

var (
	preshared_private_key []byte
	last_pull             time.Time
)

func loadPrePrivKey(ctx context.Context) (
	keyBuff []byte, err error) {
	if len(preshared_private_key) > 8 &&
		time.Since(last_pull) < 5*time.Minute {
		return preshared_private_key, nil
	} else {
		last_pull = time.Now()
	}
	info, err := cfg.Get(ctx)
	if err != nil {
		err = wUtil.Error(err)
		return
	}
	secretVal, err := gcloud.GetSecret(ctx,
		info.PrePrivKey)
	if err != nil {
		err = wUtil.Error(err)
		return
	}
	keyBuff, err = base64.StdEncoding.DecodeString(
		string(secretVal))
	if err != nil {
		err = wUtil.Error(err)
		return
	}
	return
}
