package model

import "sync"

var once sync.Once
var globalHubInstance *GlobalHub

type GlobalHub struct {
	Clients   map[*Client]bool
	Broadcast chan Event
	clientMu  sync.Mutex
}

func NewGlobalHub() Hub {
	return &GlobalHub{
		Clients:   make(map[*Client]bool),
		Broadcast: make(chan Event),
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
	gh.Clients[client] = true
}

func (gh *GlobalHub) UnregisterClient(client *Client) {
	gh.clientMu.Lock()
	defer gh.clientMu.Unlock()
	if _, ok := gh.Clients[client]; ok {
		delete(gh.Clients, client)
		close(client.Send)
	}
}

func (gh *GlobalHub) BroadcastEvent(event Event) {
	gh.Broadcast <- event.(*RoomEvent)
}

func (gh *GlobalHub) Run() {
	for {
		event := <-gh.Broadcast
		for client := range gh.Clients {
			roomEvent, ok := event.(*RoomEvent)
			if ok {
				select {
				case client.Send <- roomEvent:
				default:
					close(client.Send)
					delete(gh.Clients, client)
				}
				continue
			}
		}
	}
}
