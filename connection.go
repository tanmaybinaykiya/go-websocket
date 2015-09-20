package websocket

func (conn Connection) websocketListener() {
	for {
		messageType, message, err := conn.wsConn.ReadMessage()
		if err != nil {
			fmt.Println("tanmay" + err.Error())
			break
		}
		switch {
		case messageType == websocket.TextMessage && strings.HasPrefix(string(message), getConfig().pingMessage):
			fmt.Println("\t heartBeat Message: %s", string(message))
			conn.heartBeatChannel <- string(message)
		case messageType == websocket.TextMessage && strings.HasPrefix(string(message), getConfig().pongMessage):
			fmt.Println("\t pong Message: %s", string(message))
			if conn.pongListener != nil {
				conn.pongListener(string(message))
			}
		case messageType == websocket.TextMessage:
			fmt.Println("\t regular Message: %s", string(message))
			// if dataListener != nil {
			conn.dataListener(string(message))
			// } else {
			// 	conn.messageChannel <- string(message)
			// }
		case messageType == websocket.BinaryMessage:
			fmt.Println("Binary Frame, not handling for now")
		}
	}
	return
}

func (conn Connection) sendMessage(message string) error {
	err := conn.wsConn.WriteMessage(websocket.TextMessage, []byte(message))
	return err
}

func (conn Connection) addCloseListener(eventHandler func(string) error) {
	conn.closeListener = eventHandler
}
func (conn Connection) addOpenListener(eventHandler func(string) error) {
	conn.openListener = eventHandler
}
func (conn Connection) addPingDataListener(eventHandler func(string) error) {
	conn.pingDataListener = eventHandler
}
func (conn Connection) addPongDataListener(eventHandler func(string) error) {
	conn.pongDataListener = eventHandler
}
func (conn Connection) addDataListener(eventHandler func(string) error) {
	conn.dataListener = eventHandler
}
func (conn Connection) addPingListener(eventHandler func(string) error) {
	conn.pingListener = eventHandler
}
func (conn Connection) addPongListener(eventHandler func(string) error) {
	conn.pongListener = eventHandler
}

func (conn Connection) On(eventType string, eventHandler func(string) error) {
	switch eventType {
	case "close":
		conn.addCloseListener(eventHandler)
	case "open":
		conn.addOpenListener(eventHandler)
	case "pingData":
		conn.addPingDataListener(eventHandler)
	case "pongData":
		conn.addPongDataListener(eventHandler)
	case "data":
		conn.addDataListener(eventHandler)
	case "ping":
		conn.addPingListener(eventHandler)
	case "pong":
		conn.addPongListener(eventHandler)
	}
}

func (conn Connection) cleanUp(timer *time.Timer, err error) {
	conn.wsConn.Close()
	timer.Stop()
	close(conn.heartBeatChannel)
	// close(conn.messageChannel)
	if conn.closeListener != nil {
		conn.closeListener(err.Error())
	}
}

func (conn Connection) heartbeatHandler() {
	timer := time.NewTimer(getHeartBeatTimeout())
	go func() {
		for {
			select {
			case <-conn.heartBeatChannel:
				if conn.pingListener != nil {
					val := timer.Reset(getHeartBeatTimeout())
					if !val {
						conn.cleanUp(timer, errors.New("Heartbeat not sent in time"))
						return
					}
					conn.sendMessage(getPongMessage())
				}
			case <-timer.C:
				fmt.Println("Closing connection as heartbeat not sent in time")
				conn.cleanUp(timer, errors.New("Heartbeat not sent in time"))
				return
			}
		}
	}()
}

type Connection struct {
	wsConn           *websocket.Conn
	heartBeatChannel chan string
	// messageChannel   chan string
	closeChannel     chan error
	closeListener    func(string) error
	openListener     func(string) error
	pingDataListener func(string) error
	pongDataListener func(string) error
	pingListener     func(string) error
	pongListener     func(string) error
	dataListener     func(string) error
}
