package net

import (
	"bytes"
	"errors"
	"io"
	"log"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

func NewMultipartForm(c *gin.Context) (*MultipartForm, error) {
	if err := c.Request.ParseMultipartForm(0); err != nil {
		return nil, err
	}
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}
	return &MultipartForm{form}, nil
}

type MultipartForm struct {
	*multipart.Form
}

func (ins *MultipartForm) GetForm(key string) (string, error) {
	for _, value := range ins.Value[key] {
		log.Println("value", value)
		return value, nil
	}
	log.Println("no content name")
	return "", errors.New("no content")
}

func (ins *MultipartForm) GetFile(key string) (
	filename string, payload []byte, err error) {
	for _, file := range ins.File[key] {
		var (
			temp = bytes.NewBuffer(nil)
		)
		content, err := file.Open()
		if err != nil {
			continue
		}
		if _, err := io.Copy(temp, content); err != nil {
			continue
		}
		return file.Filename, temp.Bytes(), nil
	}
	log.Println("no content file")
	return "", nil, errors.New("no content")
}
