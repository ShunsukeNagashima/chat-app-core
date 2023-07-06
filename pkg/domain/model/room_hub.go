package model

import "sync"

type RoomHub struct {
	Clients   map[*Client]bool
	Broadcast chan *Message
	clientMu  sync.Mutex
}

func NewRoomHub() Hub {
	return &RoomHub{
		Clients:   make(map[*Client]bool),
		Broadcast: make(chan *Message),
	}
}

func (rh *RoomHub) RegisterClient(client *Client) {
	rh.clientMu.Lock()
	defer rh.clientMu.Unlock()
	rh.Clients[client] = true
}

func (rh *RoomHub) UnregisterClient(client *Client) {
	if _, ok := rh.Clients[client]; ok {
		delete(rh.Clients, client)
		close(client.Send)
	}
}

func (rh *RoomHub) BroadcastEvent(event Event) {
	rh.Broadcast <- event.(*Message)
}

func (rh *RoomHub) Run() {
	for {
		message := <-rh.Broadcast
		for client := range rh.Clients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(rh.Clients, client)
			}
		}
	}
}
