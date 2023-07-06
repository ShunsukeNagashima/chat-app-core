package model

type RoomHubManager struct {
	roomHubs map[string]Hub
}

func NewRoomHubManager() *RoomHubManager {
	return &RoomHubManager{
		roomHubs: make(map[string]Hub),
	}
}

func (hm *RoomHubManager) GetRoomHub(roomId string) (Hub, bool) {
	hub, exists := hm.roomHubs[roomId]
	return hub, exists
}

func (hm *RoomHubManager) CreateRoomHub(roomId string) Hub {
	hub := NewRoomHub()
	hm.roomHubs[roomId] = hub
	go hub.Run()
	return hub
}
