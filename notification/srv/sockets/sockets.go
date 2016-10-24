package sockets

import (
	"github.com/kazoup/gabs"
	"github.com/kazoup/platform/notification/srv/helpers"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

type User struct {
	UserId string `json:"user_id"`
}

func Notify(ws *websocket.Conn) {
	var sr *helpers.SocketRef
	isOpen := true
	uId := ""

	// Listen for clients, trying to read will fail if connection is closed
	go func() {
		for isOpen {
			var u *User
			err := websocket.JSON.Receive(ws, &u)
			// If user not nil, register known connection
			if u != nil {
				uId = u.UserId
				sr = helpers.NewSocketRef(ws, uId)

				go func() {
					for {
						select {
						case c := <-sr.Codes:
							log.Println("I'M FUCKING AWESOME", c)
							jObj := gabs.New()
							jObj.SetP(c, "code")
							if err := websocket.JSON.Send(ws, jObj.String()); err != nil {
								isOpen = false
								return
							}

							time.Sleep(time.Second)
						}
					}

				}()
				log.Println("!!!!!")
				helpers.AppendConn(sr)
				helpers.NewSocketRef(ws, uId)

				log.Println("CREATE", len(helpers.GetSocketClients().Sockets))
			}

			// Connection closed
			if err != nil {
				isOpen = false

				// Clear connection from list of socket clients
				if sr != nil {
					close(sr.Codes)
				}
				//TODO: Fix delete
				//helpers.SocketClients = helpers.SocketClients.Delete(uId)
				//log.Println("DELETE", len(helpers.SocketClients))
				break
			}

		}
	}()
	for isOpen {
		jObj := gabs.New()
		jObj.SetP(true, "refresh_datasource")

		if err := websocket.JSON.Send(ws, jObj.String()); err != nil {
			isOpen = false
			return
		}
		time.Sleep(time.Second * 2)
	}
}
