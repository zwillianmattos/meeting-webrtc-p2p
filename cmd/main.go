package main

import (
	"os"
	"strconv"
	"github.com/zwillianmattos/webrtc/packages/logger"
	"github.com/zwillianmattos/webrtc/packages/signaler"
	"github.com/zwillianmattos/webrtc/packages/turn"
	"github.com/zwillianmattos/webrtc/packages/websocket"
	"gopkg.in/ini.v1"
)

func main() {
	cfg, err := ini.Load("configs/config.ini")
	if err != nil {
		logger.Errorf("Fail to read file: %v", err)
		os.Exit(1)
	}

	publicIP := cfg.Section("turn").Key("public_ip").String()
	stunPort, err := cfg.Section("turn").Key("port").Int()
	if err != nil {
		stunPort = 3478
	}
	realm := cfg.Section("turn").Key("realm").String()

	turnConfig := turn.DefaultConfig()
	turnConfig.PublicIP = publicIP
	turnConfig.Port = stunPort
	turnConfig.Realm = realm
	turn := turn.NewTurnServer(turnConfig)

	signaler := signaler.NewSignaler(turn)
	wsServer := websocket.NewWebSocketServer(signaler.HandleNewWebSocket, signaler.HandleTurnServerCredentials)

	sslCert := cfg.Section("general").Key("cert").String()
	sslKey := cfg.Section("general").Key("key").String()
	bindAddress := cfg.Section("general").Key("bind").String()

	port := os.Getenv("PORT")
	port_int, err := strconv.Atoi(port)
    if port == "" {
        logger.Errorf("$PORT must be set")
    }

	htmlRoot := cfg.Section("general").Key("html_root").String()

	config := websocket.DefaultConfig()
	config.Host = bindAddress
	config.Port =  port_int
	config.CertFile = sslCert
	config.KeyFile = sslKey
	config.HTMLRoot = htmlRoot

	wsServer.Bind(config)
}
