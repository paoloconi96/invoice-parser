package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = pongWait * 9 / 10
	writeWait  = 10 * time.Second
)

type Websocket struct {
	hub  *Hub
	conn *websocket.Conn
	send chan *Event
}

func HandleWebsocket(ctx *gin.Context, hub *Hub) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	connection, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("unable to upgrade: %s", err)
		return
	}

	client := &Websocket{
		hub:  hub,
		conn: connection,
		send: make(chan *Event, 16),
	}
	hub.register <- client

	go client.handlePong()
	go client.handlePing()
}

func (ws *Websocket) handlePong() {
	wsClient := ws.conn
	defer func() {
		ws.hub.unregister <- ws
		// TODO: Probably we should handle err
		wsClient.Close()
	}()

	// TODO: Probably we should handle err
	wsClient.SetReadDeadline(time.Now().Add(pongWait))
	wsClient.SetPongHandler(func(string) error {
		err := wsClient.SetReadDeadline(time.Now().Add(pongWait))
		return err
	})

	for {
		_, _, err := wsClient.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws error: %v", err)
			}

			break
		}
	}
}

func (ws *Websocket) handlePing() {
	wsClient := ws.conn

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.hub.unregister <- ws
		// TODO: Probably we should handle err
		wsClient.Close()
	}()

	for {
		select {
		case message, ok := <-ws.send:
			wsClient.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				sendCloseMessage(wsClient)
				return
			}

			// TODO: Improve this payload structure
			err := wsClient.WriteJSON(message)
			if err != nil {
				fmt.Println(err)
			}
		case <-ticker.C:
			// TODO: Probably we should handle err
			wsClient.SetWriteDeadline(time.Now().Add(writeWait))
			if err := wsClient.WriteMessage(websocket.PingMessage, nil); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("ws ping error: %s", err)
				}

				return
			}
		}
	}
}

func sendCloseMessage(ws *websocket.Conn) {
	err := ws.WriteMessage(websocket.CloseMessage, nil)
	if err != nil {
		log.Printf("ws error closing: %s", err)
	}
}
