package websocket

import "github.com/paoloconi96/invoice-parser/internal/invparser"

type Hub struct {
	websockets map[*Websocket]interface{}
	Broadcast  chan *invparser.InvoiceId
	register   chan *Websocket
	unregister chan *Websocket
}

type Event struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan *invparser.InvoiceId),
		register:   make(chan *Websocket),
		unregister: make(chan *Websocket),
		websockets: make(map[*Websocket]interface{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case Websocket := <-h.register:
			h.websockets[Websocket] = nil
		case Websocket := <-h.unregister:
			if _, ok := h.websockets[Websocket]; ok {
				delete(h.websockets, Websocket)
				close(Websocket.send)
			}
		case invoiceId := <-h.Broadcast:
			message := &Event{
				Type:  "processed",
				Value: string(*invoiceId),
			}
			for Websocket := range h.websockets {
				select {
				case Websocket.send <- message:
				default:
					close(Websocket.send)
					delete(h.websockets, Websocket)
				}
			}
		}
	}
}
