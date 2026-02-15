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
	authservice "tictactoe/internal/domain/auth_service"
	"tictactoe/internal/domain/service"
	"tictactoe/internal/domain/user_service"
	_jwt "tictactoe/internal/web/jwt"
	"tictactoe/internal/web/module"
	"tictactoe/internal/web/route"
)

func NewGameStorage(lc fx.Lifecycle) (*repository.GameStorage, error) {
	log.Println("[DI] Creating GameStorage (singleton)")

	storage, err := repository.NewGameStorage()
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("[DI] Closing database connection")
			storage.DBConnection.Close(context.Background())
			return nil
		},
	})

	return storage, nil
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

func NewUserService(storage *repository.GameStorage) *user_service.UserService {
	log.Println("[DI] Creating User Service")
	return user_service.NewUserService(storage)
}

func NewJWTProvider() *_jwt.JwtProvider {
	conf := _jwt.NewJWTConfig()
	return _jwt.NewJWTProvider(conf.SecretAccessToken, conf.SecretRefreshToken, conf.AccessTTL, conf.RefreshTTL)
}

func NewAuthService(userService *user_service.UserService, jwtProvider *_jwt.JwtProvider) authservice.AuthService {
	log.Println("[DI] Creating AuthService")
	return authservice.NewAuthService(userService, jwtProvider)
}

func NewAuthHandler(authservice authservice.AuthService) *module.AuthHandler{
	log.Println("[DI] Creating AuthHandler")
	return module.NewAuthHandler(authservice)
}

func NewAuthMiddlewareHandler(jwtProvider *_jwt.JwtProvider) *module.UserAuthenticator {
	log.Println("[DI] Creating Middleware Handler")
	return module.NewUserAuthenticator(jwtProvider)
}

func NewRouter(handler *module.GameHandler, authHandler *module.AuthHandler, middleware *module.UserAuthenticator) http.Handler {
	log.Println("[DI] Creating Router")
	return route.NewRouter(handler, authHandler, middleware)
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

				log.Println("[DI]   POST  /auth/signup		- Регистрация")
				log.Println("[DI]   POST  /auth/signin		- Авторизация")
				log.Println("[DI]   POST  /refresh_access	- Обновить access token")
				log.Println("[DI]   POST  /refresh_tokens	- Обновить refresh token")
				log.Println("[DI]   POST  /game			- Создать игру")
				log.Println("[DI]   POST  /game/join/{GameID}	- Подключиться")
				log.Println("[DI]   POST  /game/{GameID}		- Сделать ход")
				log.Println("[DI]   GET   /game/{GameID}		- Получить инфо об игре")
				log.Println("[DI]   GET   /waiting		- Получить инфо об играх в ожидании")
				log.Println("[DI]   GET   /user/{UserID}		- Получить инфо о пользователе по ID")
				log.Println("[DI]   GET   /user/info_by_access_token	- Получить инфо о пользователе по токену")
				log.Println("[DI]   GET   /history		- Получить список завершенных игр для пользователя")
				log.Println("[DI]   GET   /leaderboard?limit=10	- Получить таблицу лидеров")
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
