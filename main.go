package main

import (
	"github.com/tanmaybinaykiya/websocket"
)

func main() {
	config = Config{
		pingMessage: "primus::ping",
		pongMessage: "primus::pong",
		pingTimeout: 5 * time.Second,
	}

	var h HTTPHandler
	err := http.ListenAndServe("localhost:4000", h)
	if err != nil {
		fmt.Println(err)
	}
}
