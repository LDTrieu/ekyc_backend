package model

import "fmt"

var (
	ErrWsUnauth       = fmt.Errorf("client unauth")
	ErrDocNotFound    = fmt.Errorf("doc not found")
	ErrDocIdEmpty     = fmt.Errorf("doc id empty")
	ErrPayloadInvalid = fmt.Errorf("payload invalid")
)
