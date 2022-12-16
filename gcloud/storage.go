package gcloud

import (
	"context"
	"ekyc-app/package/wlog"
	"ekyc-app/source/cfg"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
)

var (
	storageClient *storage.Client
)

func getStoreClient(ctx context.Context) (_ *storage.Client, err error) {
	if storageClient != nil {
		return storageClient, nil
	}
	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		return
	}
	return storageClient, nil
}
func SaveFaceImage(ctx context.Context, accId string, file []byte) (
	uri string, err error) {
	uri = accId + "/faceimage"
	info, err := cfg.Get(ctx)
	if err != nil {
		return
	}

	if err := saveFile(ctx, uri, info.TerminalFileBucket, file); err != nil {
		log.Println("err", err)
		return "", err
	}
	return
}

func SaveFaceImageFile(ctx context.Context, account_id, filename string, file []byte) (string, error) {
	var (
		uri = fmt.Sprintf("%s/validate_%s", account_id, filename)
	)
	info, err := cfg.Get(ctx)
	if err != nil {
		return "", err
	}
	if err := saveFile(ctx, uri, info.TerminalFileBucket, file); err != nil {
		return "", err
	}
	return uri, nil
}
func saveFile(ctx context.Context, uri, bucket string, file []byte) (
	err error) {
	client, err := getStoreClient(ctx)
	if err != nil {
		err = wlog.Error(ctx, err)
		return
	}
	o := client.Bucket(bucket).Object(uri)
	wc := o.NewWriter(ctx)
	if _, err = wc.Write(file); err != nil {
		return
	}
	err = wc.Close()
	return
}

// func loadFile(ctx context.Context, uri, bucket string) (
// 	file []byte, err error) {
// 	client, err := getStoreClient(ctx)
// 	if err != nil {
// 		err = wlog.Error(ctx, err)
// 		return
// 	}
// 	o := client.Bucket(bucket).Object(uri)
// 	rc, err := o.NewReader(ctx)
// 	if err != nil {
// 		return
// 	}
// 	defer rc.Close()
// 	file, err = io.ReadAll(rc)
// 	return
// }
