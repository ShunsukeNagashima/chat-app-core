package model

type HubManager struct {
	Hubs map[string]*Hub
}

func NewHubManager() *HubManager {
	return &HubManager{
		Hubs: make(map[string]*Hub),
	}
}

func (hm *HubManager) GetHub(roomId string) *Hub {
	hub, exists := hm.Hubs[roomId]
	if !exists {
		hub = NewHub()
		hm.Hubs[roomId] = hub
		go hub.Run()
	}
	return hub
}
