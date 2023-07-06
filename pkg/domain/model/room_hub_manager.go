package model

type RoomHubManager struct {
	RoomHubs map[string]Hub
}

func NewRoomHubManager() *RoomHubManager {
	return &RoomHubManager{
		RoomHubs: make(map[string]Hub),
	}
}

func (hm *RoomHubManager) GetRoomHub(roomId string) (Hub, bool) {
	hub, exists := hm.RoomHubs[roomId]
	return hub, exists
}

func (hm *RoomHubManager) CreateRoomHub(roomId string) Hub {
	hub := NewRoomHub()
	hm.RoomHubs[roomId] = hub
	go hub.Run()
	return hub
}
