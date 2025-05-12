package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ankit8125/crud-golang-practice/internal/config"
)

func main(){
	// fmt.Println("Entry point of students-api")
	
	// load config
	cfg := config.MustLoad()
	// database setup
	
	// setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("welcome to students api"))
	})
	
	// setup server
	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}
	slog.Info("Server started", slog.String("address", cfg.Addr)) 
	fmt.Printf("Server started %s", cfg.Addr)  

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}	
	} ()

	<-done
	// Graceful shutdown of server
	slog.Info("Shutting down the server")
	 
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shut down server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
}