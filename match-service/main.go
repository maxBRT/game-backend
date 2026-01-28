package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	h "github.com/maxbrt/game-backend/match-service/internal/handlers"
	q "github.com/maxbrt/game-backend/match-service/internal/queue"
)

func main() {
	queueManager := q.NewQueueManager(q.NewInMemoryQueue(), q.NewInMemoryQueue(), q.NewInMemoryMatchStore())
	matcher := q.NewMatcher(queueManager, 100*time.Millisecond)
	matcher.Start()

	joinHandler := h.NewJoinHandler(queueManager)
	statusHandler := h.NewStatusHandler(queueManager, 10*time.Second, 100*time.Millisecond)

	mux := http.NewServeMux()
	mux.Handle("POST /match/join", joinHandler)
	mux.Handle("GET /match/status/{ticketID}", statusHandler)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:         ":8000",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second, // Must be > pollTimeout (10s)
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		fmt.Println("\nShutting down...")
		server.Close()
	}()

	fmt.Println("Match service started on port 8000")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
