package main

import (
	"log"

	cleanupScripts "github.com/shunsukenagashima/chat-api/scripts/cleanup"
	setupScripts "github.com/shunsukenagashima/chat-api/scripts/setup"
)

func main() {
	if err := cleanupScripts.CleanUpDynamodb(); err != nil {
		log.Panicf("Failed to clean up dynamodb: %v", err)
	}

	if err := cleanupScripts.CleanUpElasticsearch(); err != nil {
		log.Panicf("Failed to clean up elasticsearch: %v", err)
	}

	users, err := setupScripts.SetupUsers()
	if err != nil {
		log.Panicf("Failed to set up users: %v", err)
	}

	roomIDs, err := setupScripts.SetupRooms()
	if err != nil {
		log.Panicf("Failed to set up rooms: %v", err)
	}

	if err := setupScripts.SetupRoomUsers(roomIDs, users); err != nil {
		log.Panicf("Failed to set up room users: %v", err)
	}

	if err := setupScripts.SetupLikes(); err != nil {
		log.Panicf("Failed to set up likes: %v", err)
	}

	if err := setupScripts.SetupReadby(); err != nil {
		log.Panicf("Failed to set up readby: %v", err)
	}

	if err := setupScripts.SetupMessages(users, roomIDs[0]); err != nil {
		log.Panicf("Failed to set up messages: %v", err)
	}

	if err := setupScripts.SetUpElasticsearch(users); err != nil {
		log.Panicf("Failed to set up elasticsearch: %v", err)
	}
}
