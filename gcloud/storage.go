package gcloud

import (
	"context"
	"ekyc-app/internal/cfg"
	"ekyc-app/package/wlog"
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

func getFileAttributes(ctx context.Context, uri, bucket string) (contentType string, size int64, err error) {
	client, err := getStoreClient(ctx)
	if err != nil {
		return "", 0, err
	}
	o := client.Bucket(bucket).Object(uri)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return "", 0, err
	}
	return attrs.ContentType, attrs.Size, nil
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
func SaveEkycImageFile(ctx context.Context, student_id, filename string, file []byte) (string, error) {
	var (
		uri = fmt.Sprintf("%s/ekyc_%s", student_id, student_id)
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
func SaveFaceVideoFile(ctx context.Context, student_id, filename string, file []byte) (string, error) {
	var (
		uri = fmt.Sprintf("%s/face_video_%s", student_id, student_id)
		//uri_thumbnail = fmt.Sprintf("%s/face_thumbnail_%s", student_id, student_id)
	)
	info, err := cfg.Get(ctx)
	if err != nil {
		return "", err
	}
	// save video
	if err := saveFile(ctx, uri, info.TerminalFileBucket, file); err != nil {
		return "", err
	}

	return uri, nil
}

func SaveFaceThumnailFile(ctx context.Context, student_id, filename string, file []byte) (string, error) {
	var (
		uri_thumbnail = fmt.Sprintf("%s/face_thumbnail_%s", student_id, student_id)
	)
	info, err := cfg.Get(ctx)
	if err != nil {
		return "", err
	}
	// save video
	if err := saveFile(ctx, uri_thumbnail, info.TerminalFileBucket, file); err != nil {
		return "", err
	}

	return uri_thumbnail, nil
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
