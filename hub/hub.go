package hub

import "fmt"

type Hub struct {
	clients map[*Client] bool
	Broadcast chan []byte
	register chan *Client
	unregister chan *Client
}

func NewDeploymentsHub() *Hub {
	return &Hub {
		Broadcast: make(chan []byte),
		register:  make(chan *Client),
		unregister:  make(chan *Client),
		clients : make(map[*Client] bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
		h.clients[client] = true
		fmt.Println("Client registered")
		case client := <- h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				fmt.Println("Client disconnected")
			}
		case message := <- h.Broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)

				}
			}
		}
	}
}
