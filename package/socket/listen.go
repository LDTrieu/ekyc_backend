package socket

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func Listen(r *http.Request, conn *websocket.Conn, read chan<- []byte,
	write <-chan []byte) error {
	var (
		errChan  = make(chan error)
		safeSend = func(ch chan error, err error) {
			defer recover()
			// write
			ch <- err
		}
	)
	// Init read message channel
	go func() {
		defer recover()
		for {
			mtype, buff, err := conn.ReadMessage()
			if err != nil {
				safeSend(errChan, err)
				return
			}
			switch mtype {
			case websocket.TextMessage:
				println("<< Message type: TextMessage")
			case websocket.BinaryMessage:
				println("<< Other type: BinaryMessage")
			case websocket.CloseMessage:
				println("<< Other type: CloseMessage")
				continue
			case websocket.PingMessage:
				println("<< Other type: PingMessage")
				continue
			case websocket.PongMessage:
				println("<< Other type: PongMessage")
				continue
			default:
				println("<< Unknown type (", mtype, ")")
			}
			read <- buff
		}
	}()

	println(fmt.Sprintf("[LOG] %s | 101 | %13s | %15s | %10s  %#v",
		time.Now().Format("2006/01/02 - 15:04:05"),
		"WEBSOCKET",
		func() string {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				return "x.x.x.x"
			}
			return net.ParseIP(ip).String()
		}(),
		"CONNECTED",
		r.RequestURI,
	))

	// Wait and write to websocket
	for {
		select {
		case buff, ok := <-write:
			if !ok {
				return errors.New("writer channel has closed")
			}
			println(">> Write:", bytes.NewBuffer(buff).String())
			if err := conn.WriteMessage(websocket.TextMessage, buff); err != nil {
				return err
			}
		case err, ok := <-errChan:
			if !ok {
				return errors.New("writer channel has closed")
			}
			if websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
				return nil
			}
			return err
		}
	}
}
