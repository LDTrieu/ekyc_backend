package qrcode

import "github.com/skip2/go-qrcode"

func GenQRCode(content string) (qrByte []byte) {
	qrCodeObject, err := GenQRObject(content)
	if err != nil {
		return
	}
	qrByte, err = qrCodeObject.PNG(256)
	if err != nil {
		return nil
	}
	return qrByte
}

func GenQRObject(content string) (
	qrCodeObject *qrcode.QRCode, err error) {
	qrCodeObject, err = qrcode.New(content, qrcode.Low)
	if err != nil {
		return nil, err
	}
	return qrCodeObject, nil
}
