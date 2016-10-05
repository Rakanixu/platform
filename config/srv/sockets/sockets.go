package sockets

import (
	"github.com/kazoup/gabs"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

func PingPlatform(ws *websocket.Conn) {
	for {
		jObj := gabs.New()
		jObj.SetP(time.Now().Unix(), "timestamp")

		if err := websocket.JSON.Send(ws, jObj.String()); err != nil {
			log.Println("SOCKET ERROR", err)
			//return
		}
		time.Sleep(time.Second * 3)
	}
}
