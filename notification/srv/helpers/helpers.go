package helpers

import (
	"golang.org/x/net/websocket"
)

type SocketRef struct {
	SocketConn *websocket.Conn
	ID         string
	Codes      chan string
}

type ClientConnections []*SocketRef

var SocketClients ClientConnections

func init() {

	/*	// Keep a reference to all connections
		sockets.SocketClients = make(sockets.ClientConnections, 0) // Keep a reference to all connections
	*/

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
func (sc ClientConnections) Delete(id string) ClientConnections {
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
}

// Delete returns open ClientConnections
func AppendConn(sr *SocketRef) {
	SocketClients = append(SocketClients, sr)
}

func GetSocketClients() ClientConnections {
	return SocketClients
}
