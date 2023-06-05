package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shunsukenagashima/chat-api/pkg/interface/routes"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("Server failed to run with %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	router := gin.Default()

	routes.RegisterRoutes(router)

	return router.Run()
}
