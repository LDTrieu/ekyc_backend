package ws

import (
	"encoding/json"
	"errors"
)

type WebsocketType string

const (
	CmdLoginRedirect WebsocketType = "LOGIN_REDIREC"
	//CmdPortalNotify  WebsocketType = "PORTAL_NOTIFY"
)

type RedirectType string

const (
	RedirectLogin RedirectType = "LOGIN"
)

type WsRequestModel struct {
	Command WebsocketType `json:"command"` //LOGIN_REDIRECT

	/*
		Fileds for LOGIN_REDIREC
	*/
	JWT      string       `json:"jwt, omitempty"`
	Redirect RedirectType `json:"redirect,omitempty"` // MISS_WALLET; LOGIN; STORE_REGISTER; STORE_SUBMITED
}

var Station = &station{
	clients: make(map[string]chan<- []byte, 0),
}

type station struct {
	clients map[string]chan<- []byte
}

func (ins *station) AddClient(connection_id string,
	writer chan<- []byte) {
	if w, ok := ins.clients[connection_id]; ok {
		close(w)
	}
	println("[WS] add", connection_id)
	ins.clients[connection_id] = writer
}

func (ins *station) PushSender(connection_id string,
	v any) error {
	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if w, ok := ins.clients[connection_id]; ok {
		defer recover()
		println("[WS] push", connection_id, "ok")
		w <- payload
		return nil
	}
	return errors.New("sender " + connection_id + " closed or does not exist")
}
