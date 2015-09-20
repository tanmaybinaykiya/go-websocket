package websocket

import (
	"fmt"
	"github.com/tanmaybinaykiya/websocket"
	"net/http"
	"time"
)

type HTTPHandler struct {
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// func (conn Connection) dataListener(string) error

func getCurrentTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func getPongMessage() string {
	return fmt.Sprintf(getConfig().pongMessage+":%d", getCurrentTime())
}

func (h HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Upgrade") == "websocket" {
		wconn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("err upgrading " + err.Error())
			return
		}
		fmt.Println("Connection succesfully upgraded")
		conn := Connection{
			wsConn:           wconn,
			heartBeatChannel: make(chan string),
			closeChannel:     make(chan error),
		}
		conn.dataListener = func(message string) error {
			err := conn.sendMessage("Reply: " + message)
			return err
		}
		if conn.openListener != nil {
			conn.openListener("Connection succesfully upgraded")
		}
		conn.heartbeatHandler()
		conn.websocketListener()
	} else {
		w.Write([]byte("Hi! This is Tanmay's server. You should not be meddling around here"))
	}
}
