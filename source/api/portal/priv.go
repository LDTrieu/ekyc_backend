package portal

import (
	"context"
	"ekyc-app/package/socket"
	"ekyc-app/source/auth"
	"ekyc-app/source/fsdb"
	"ekyc-app/source/model"
	"ekyc-app/source/ws"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func guestRendQRLogin(ctx context.Context,
	request *rendQRLoginRequest) (
	rendQRLoginResponse, error) {
	var (
		login_session_id = uuid.NewString()
	)
	// make new JWT authen with account_id equal to empty
	_, jwt_login, err := auth.GenerateJWTLoginSession(
		ctx, login_session_id, "")
	if err != nil {
		return rendQRLoginResponse{
			Code:    model.StatusInternalServerError,
			Message: err.Error()}, err
	}
	//save to cache
	doc_id, err := fsdb.LoginSessionDBC.
		Add(ctx, login_session_id,
			jwt_login.Token, jwt_login.ExpiresAt)
	if err != nil {
		return rendQRLoginResponse{Code: model.StatusServiceUnavailable,
			Message: err.Error()}, err
	}
	log.Println(doc_id)
	// end
	return rendQRLoginResponse{Payload: rend_qr_code_login_resp{
		//Path:      fmt.Sprintf("/login/qr/download/%s/%s?action=%s", doc_id, uuid.NewString(), qrActionLoginWebPortal),
		JWT:       jwt_login.Token,
		IssuedAt:  jwt_login.IssuedAt,
		ExpiresIn: jwt_login.ExpiresAt - jwt_login.IssuedAt,
	}}, nil

}

/* WS */
func websocketConnection(c *gin.Context) {
	var (
		r   = c.Request
		w   = c.Writer
		ctx = c.Request.Context()

		get_token = func(r *http.Request) (string, error) {
			h := strings.TrimSpace(r.Header.Get("Sec-Websocket-Protocol"))
			if h == "" {
				return "", errors.New("subprotocols is invalid")
			}
			protocols := strings.Split(h, ",")
			for i := range protocols {
				protocols[i] = strings.TrimSpace(protocols[i])
			}
			if len(protocols) > 0 {
				return protocols[0], nil
			}
			return "", errors.New("subprotocols is invalid")
		}

		get_session_id = func(token string) string {
			// Try search on LoginSessionDBC
			if _, session_id, ok, err := fsdb.LoginSessionDBC.GetSessionId(ctx, token); err != nil {
				println("Websocket-Connecting:", err.Error(), "\nTry again search on Profile.")
			} else if ok && len(session_id) > 0 {
				return session_id
			}
			// Retry search on ProfileToken
			if _, session_id, ok, err := fsdb.PersonProfile.GetSessionID(ctx, token); err != nil {
				println("Websocket-Connecting:", err.Error())
				return ""
			} else if !ok {
				println("Websocket-Connecting: token does not exist")
				return ""
			} else {
				return session_id
			}
		}

		connection_id = "default"
		read          = make(chan []byte, 100)
		write         = make(chan []byte, 100)

		websocket_receiver = func(ctx context.Context, read <-chan []byte) {
			var failure_count int
			for {
				select {
				case <-ctx.Done():
					if err := ctx.Err(); err != nil {
						println("Close connection by client:", err.Error())
					} else {
						println("Close connection by client.")
					}
					return
				case _, ok := <-read:
					if !ok {
						if failure_count > 3 {
							println("can not read payload data becase READ channel is closed.")
						}
						failure_count++
						time.Sleep(300 * time.Millisecond)
						continue
					}
				}
			}
		}
	)

	token, err := get_token(r)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	println("Websocket-Connecting: Protocol [", token, "]")
	connection_id = get_session_id(token)
	if len(connection_id) == 0 {
		println("Websocket-Connecting: disconnected")
		c.AbortWithError(http.StatusUnauthorized, errors.New("token does not exist"))
		return
	}
	println("Websocket-Connecting: ConnectionID=", connection_id)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  0,
		WriteBufferSize: 0,
		Subprotocols:    []string{token},
		CheckOrigin: func(r *http.Request) bool {
			// allow all connections by default
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
		return
	}
	defer conn.Close()

	go websocket_receiver(ctx, read)

	// Add connection to Station
	ws.Station.AddClient(connection_id, write)

	// create read-only-chanel (task) & write-only-chanel (task) and put in the function Listen
	if err := socket.Listen(r, conn, read, write); err != nil {
		println("Error:", err.Error())
	}
}
