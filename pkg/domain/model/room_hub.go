package model

import "sync"

type RoomHub struct {
	clients   map[*Client]bool
	broadcast chan *Message
	clientMu  sync.Mutex
}

func NewRoomHub() Hub {
	return &RoomHub{
		clients:   make(map[*Client]bool),
		broadcast: make(chan *Message),
	}
}

func (rh *RoomHub) RegisterClient(client *Client) {
	rh.clientMu.Lock()
	defer rh.clientMu.Unlock()
	rh.clients[client] = true
}

func (rh *RoomHub) UnregisterClient(client *Client) {
	if _, ok := rh.clients[client]; ok {
		delete(rh.clients, client)
		close(client.Send)
	}
}

func (rh *RoomHub) BroadcastEvent(event Event) {
	rh.broadcast <- event.(*Message)
}

func (rh *RoomHub) Run() {
	for {
		message := <-rh.broadcast
		for client := range rh.clients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(rh.clients, client)
			}
		}
	}
}
