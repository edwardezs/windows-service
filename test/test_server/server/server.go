package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Server struct {
	httpServer *http.Server
	stopChan   chan os.Signal
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func New() *Server {
	http.HandleFunc("/hello", hello)
	return &Server{
		httpServer: &http.Server{Addr: ":8080"},
		stopChan:   make(chan os.Signal, 1),
	}
}

func (s *Server) Run() error {
	log.Info().Msg("Starting server")

	signal.Notify(s.stopChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Server listening failed")
		}
	}()

	<-s.stopChan

	log.Info().Msg("Shutting down server")

	if err := s.Stop(); err != nil {
		log.Error().Err(err).Msg("Failed to stop server")
		return errors.Wrap(err, "failed to stop server")
	}
	log.Info().Msg("Server stopped")
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to shutdown server")
	}

	return nil
}
