package model

type Hub interface {
	RegisterClient(*Client)
	UnregisterClient(*Client)
	BroadcastEvent(Event)
	Run()
}
