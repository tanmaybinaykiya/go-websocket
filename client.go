package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"net/url"
	"time"
)

func sendHeartbeats(wsConn *websocket.Conn) {
	ticker := time.NewTicker(3 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				wsConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("primus::ping:%d", time.Now().UnixNano()/int64(time.Millisecond))))
				// scount++
				// fmt.Printf("s: %d\n", scount)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func addReadListener(recd chan []byte, quit chan bool) {
	go func() {
		for {
			select {
			case <-recd:
				// rcount++
				fmt.Println("received: " + string(<-recd))
				// fmt.Printf("r: %d\n", rcount)
			}
		}
	}()
}

func startClient() {
	fmt.Println("start...")
	// rcount = 0
	// scount = 0
	url, err := url.Parse("http://localhost:4000/")
	if err != nil {
		fmt.Println("url parse failed")
		fmt.Println(err.Error())
		return
	}

	conn, err := net.Dial("tcp", url.Host)
	if err != nil {
		fmt.Println("dial failed")
		fmt.Println(err.Error())
		return
	}
	if conn == nil {
		fmt.Println("Conn nil")
	}

	requestHeader := http.Header{
		"Origin":                   {"http://localhost:4000"},
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
	}

	wsConn, response, err := websocket.NewClient(conn, url, requestHeader, 1024, 1024)
	if err != nil {
		fmt.Println("newClient failed")
		fmt.Println(response)
		fmt.Println(err.Error())
	}
	// fmt.Println("Creating Channels")
	recd := make(chan []byte, 1)
	quit := make(chan bool, 2)

	wsConn.WriteMessage(websocket.TextMessage, []byte("Hey"))
	// scount++
	// fmt.Printf("s: %d\n", scount)

	addReadListener(recd, quit)
	sendHeartbeats(wsConn)

	for {
		_, received, err := wsConn.ReadMessage()
		if err != nil {
			fmt.Println("read failed")
			fmt.Println(err.Error())
			break
		}
		recd <- received
	}
}

func main() {

	for i := 0; i < 10; i++ {
		go startClient()
	}
	select {}
}
