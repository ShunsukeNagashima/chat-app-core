package model

type HubManager struct {
	Hubs map[string]*Hub
}

func NewHubManager() *HubManager {
	return &HubManager{
		Hubs: make(map[string]*Hub),
	}
}

func (hm *HubManager) GetHub(roomID string) *Hub {
	hub, exists := hm.Hubs[roomID]
	if !exists {
		hub = NewHub()
		hm.Hubs[roomID] = hub
		go hub.Run()
	}
	return hub
}
