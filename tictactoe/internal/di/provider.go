package di

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"

	"tictactoe/internal/algorithm/minimax"
	"tictactoe/internal/datasource/repository"
	"tictactoe/internal/domain/service"
	"tictactoe/internal/web/module"
	"tictactoe/internal/web/route"
)

func NewGameStorage() *repository.GameStorage {
	log.Println("[DI] Creating GameStorage (singleton)")
	return repository.NewGameStorage()
}

func NewGameRepository(storage *repository.GameStorage) repository.GameRepository {
	log.Println("[DI] Creating GameRepository")
	return repository.NewGameRepo(storage)
}

func NewGameService(repo repository.GameRepository, algo service.MinimaxAlgorithm) service.GameService {
	log.Println("[DI] Creating GameService")
	return service.NewGameService(repo, algo)
}

func NewMinimax() service.MinimaxAlgorithm {
	log.Println("[DI] Creating MiniMax")
	return minimax.NewMinimax(2)
}

func NewGameHandler(service service.GameService) *module.GameHandler {
	log.Println("[DI] Creating GameHandler")
	return module.NewGameHandler(service)
}

func NewRouter(handler *module.GameHandler) http.Handler {
	log.Println("[DI] Creating Router")
	return route.NewRouter(handler)
}

func RegisterServer(lc fx.Lifecycle, router http.Handler) {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	
	setupGracefulShutdown(server)
	
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("[DI] Starting HTTP server on :8080")
			
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("[DI] HTTP server error: %v", err)
				}
			}()
			
			go func() {
				time.Sleep(100 * time.Millisecond)
				log.Println("[DI] Server is ready")
				log.Println("[DI] Available endpoints:")
				log.Println("[DI]   POST   /game          - Create new game")
				log.Println("[DI]   GET    /game/{id}     - Get game info")
				log.Println("[DI]   POST   /game/{id}     - Make a move")
				log.Println("[DI]   GET    /health        - Health check")
			}()
			
			return nil
		},
		
		OnStop: func(ctx context.Context) error {
			log.Println("[DI] Shutting down HTTP server...")
			
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			
			return server.Shutdown(shutdownCtx)
		},
	})
}

func setupGracefulShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-stop
		log.Println("\n[DI] Received shutdown signal")
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("[DI] Server shutdown error: %v", err)
		}
		
		log.Println("[DI] Server stopped")
		os.Exit(0)
	}()
}