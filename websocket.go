package main

import (
	"log"
	"strconv"
	"time"

	"golang.org/x/net/websocket"
)

type websocketMsg struct {
	Type string
	Msg  interface{}
}

var connection *websocket.Conn

func wsHandler(ws *websocket.Conn) {
	log.Println("Websocket connected.")
	if connection != nil {
		connection.Close()
	}
	connection = ws
	for connection != nil && connection.LocalAddr().String() == connection.LocalAddr().String() {

	}
}

func wsWorker() {
	for {
		time.Sleep(time.Second)
		if connection == nil {
			continue
		}
		sendMetrics()
		sendPhases()
	}
}

func sendToWebsocket(msg websocketMsg) {
	if connection == nil {
		return
	}
	err := websocket.JSON.Send(connection, msg)
	if err != nil {
		connection = nil
	}
}

func sendMetrics() {
	m := make(map[string]int)
	m["Rounds"] = roundsMetrics
	m["Fights"] = fightMetrics
	sendToWebsocket(websocketMsg{"Metrics", m})
	roundsMetrics = 0
	fightMetrics = 0
}

func sendPhases() {
	m := make(map[string]phaseInformation)
	for i := 0; i < len(currentConfig.Phases); i++ {
		m[strconv.Itoa(i)] = phaseInformation{phaseTop[i], currentConfig.Phases[i].Total, currentConfig.Phases[i].Passed, currentConfig.Phases[i].GetFailed(), currentConfig.Phases[i].Bestscore}
	}
	sendToWebsocket(websocketMsg{"Phases", m})
}
