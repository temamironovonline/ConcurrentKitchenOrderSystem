package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"main.go/handlers"
)

type Server struct {
	httpHandlers *handlers.HTTPHandlers
}

func NewServer(httpHandlers *handlers.HTTPHandlers) *Server {
	return &Server{
		httpHandlers: httpHandlers,
	}
}
func (s *Server) Start(ctx context.Context) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /orders", s.httpHandlers.HandleCreateOrder)
	mux.HandleFunc("GET /orders/{id}", s.httpHandlers.HandleGetIntermediateResult)
	mux.HandleFunc("PATCH /orders/{id}", s.httpHandlers.HandleCancelOrder)
	//mux Stats

	httpServer := &http.Server{
		Addr:    ":9091",
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		log.Println("server shutting down...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Println("graceful shutdown error:", err)
		}
	}()

	log.Println("server start")

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Println("connecting to HTTP server error")
	}
}
