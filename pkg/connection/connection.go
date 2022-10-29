package connection

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

type Connection struct {
	connection *websocket.Conn
}

// TODO: Implement auth

func Connect(address string) Connection {
	url := url.URL{Scheme: "ws", Host: address, Path: "/"}
	fmt.Printf("Connecting to %s", url.String())

	connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		fmt.Errorf("Error dialing %s: %v", url.String(), err)
	}

	return Connection{connection: connection}
}

func (c Connection) Notify(payload string) error {
	return c.connection.WriteMessage(websocket.TextMessage, []byte(payload))
}

func (c Connection) Read() ([]byte, error) {
	_, message, err := c.connection.ReadMessage()

	return message, err
}

func (c Connection) Close() error {
	c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return c.connection.Close()
}
