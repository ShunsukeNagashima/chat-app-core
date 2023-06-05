package main

import (
	"log"

	scripts "github.com/shunsukenagashima/chat-api/scripts/setup"
)

func main() {
	if err := scripts.SetupUsers(); err != nil {
		log.Panicf("Failed to set up users: %v", err)
	}

	if err := scripts.SetupLikes(); err != nil {
		log.Panicf("Failed to set up likes: %v", err)
	}

	if err := scripts.SetupReadby(); err != nil {
		log.Panicf("Failed to set up readby: %v", err)
	}

	if err := scripts.SetUpRooms(); err != nil {
		log.Panicf("Failed to set up rooms: %v", err)
	}

	if err := scripts.SetupMessages(); err != nil {
		log.Panicf("Failed to set up messages: %v", err)
	}
}
