package main

import (
	"os"

	"win-svc/test/test_server/server"

	"github.com/rs/zerolog/log"
)

func main() {
	srv := server.New()
	if err := srv.Run(); err != nil {
		log.Error().Err(err).Msg("An error occurred while running the server")
		os.Exit(1)
	}
}
