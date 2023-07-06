package model

import "sync"

var once sync.Once
var globalHubInstance *GlobalHub

type GlobalHub struct {
	clients   map[*Client]bool
	broadcast chan Event
	clientMu  sync.Mutex
}

func NewGlobalHub() Hub {
	return &GlobalHub{
		clients:   make(map[*Client]bool),
		broadcast: make(chan Event),
	}
}

func GetGlobalHubInstance() *GlobalHub {
	once.Do(func() {
		globalHubInstance = NewGlobalHub().(*GlobalHub)
		go globalHubInstance.Run()
	})
	return globalHubInstance
}

func (gh *GlobalHub) RegisterClient(client *Client) {
	gh.clientMu.Lock()
	defer gh.clientMu.Unlock()
	gh.clients[client] = true
}

func (gh *GlobalHub) UnregisterClient(client *Client) {
	gh.clientMu.Lock()
	defer gh.clientMu.Unlock()
	if _, ok := gh.clients[client]; ok {
		delete(gh.clients, client)
		close(client.Send)
	}
}

func (gh *GlobalHub) BroadcastEvent(event Event) {
	gh.broadcast <- event.(*RoomUserDetails)
}

func (gh *GlobalHub) Run() {
	for {
		event := <-gh.broadcast
		for client := range gh.clients {
			eventData, ok := event.(*RoomUserDetails)
			if ok {
				select {
				case client.Send <- eventData:
				default:
					close(client.Send)
					delete(gh.clients, client)
				}
				continue
			}
		}
	}
}
