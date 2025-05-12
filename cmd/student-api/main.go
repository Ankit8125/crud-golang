package main
// How to run? => go run cmd/student-api/main.go -config config/local.yaml
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
	"github.com/ankit8125/crud-golang-practice/internal/http/handlers/student"
	"github.com/ankit8125/crud-golang-practice/internal/storage/sqlite"
)

func main(){
	// fmt.Println("Entry point of students-api")
	
	// load config
	cfg := config.MustLoad()
	
	// database setup
	storage, err := sqlite.New(cfg) // cfg is a pointer
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0 "))

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))

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

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shut down server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
}