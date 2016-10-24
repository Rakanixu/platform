package helpers

import (
	"golang.org/x/net/websocket"
)

type SocketRef struct {
	SocketConn *websocket.Conn
	ID         string
	Codes      chan string
}

type socketRefs struct {
	Sockets []*SocketRef
}

var (
	socketClients *socketRefs
)

func init() {
	socketClients = &socketRefs{
		Sockets: make([]*SocketRef, 0),
	}
}

// NewSocketRef constructor
func NewSocketRef(ws *websocket.Conn, id string) *SocketRef {
	return &SocketRef{
		SocketConn: ws,
		ID:         id,
		Codes:      make(chan string, 10),
	}
}

// Delete returns open ClientConnections
/*func (sc ClientConnections) Delete(id string) []*SocketRef {
	index := 0
	found := false

	for k, v := range sc {
		if id == v.ID {
			index = k
			found = true
			break
		}
	}

	if found {
		sc = append(sc[:index], sc[index+1:]...)
	}

	return sc
}*/

// Delete returns open ClientConnections
func AppendConn(sr *SocketRef) {
	socketClients.Sockets = append(socketClients.Sockets, sr)
}

func GetSocketClients() *socketRefs {
	return socketClients
}
