package utils

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"

	db "github.com/tempfiles-Team/tempfiles-backend/database"
)

// StartServerWithGracefulShutdown function for starting server with a graceful shutdown.
func StartServerWithGracefulShutdown(a *fiber.App) {
	// Create channel for idle connections.
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := a.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnsClosed)
	}()

	// Run server.
	if err := a.Listen(os.Getenv("SERVER_URL")); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}

func StartServer(a *fiber.App) {
	if os.Getenv("BACKEND_PORT") == "" {
		os.Setenv("BACKEND_PORT", "5050")
	}

	log.Fatal(a.Listen(fmt.Sprintf(":%s", os.Getenv("BACKEND_PORT"))))
}

func ReadyComponent() {
	_ = startBatch()

	if err := CheckTmpFolder(); err != nil {
		log.Fatalf("tmp folder error: %v", err)
	}

	if err := db.NewConnection(); err != nil {
		log.Fatalf("db connection error: %v", err)
	}

	log.Println("Ready to All Server Components! 🚀")
}
