package websocket

type Config struct {
	pingMessage string
	pongMessage string
	pingTimeout time.Duration
}

func getConfig() Config {
	return config
}

func getHeartBeatTimeout() time.Duration {
	return getConfig().pingTimeout
}

var config Config
