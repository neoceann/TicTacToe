package route

import (
	"net/http"
	
	"tictactoe/internal/web/module"
)

func NewRouter(
	gameHandler *module.GameHandler, authHandler *module.AuthHandler, authMiddleware *module.UserAuthenticator) http.Handler {

	mux := http.NewServeMux()
	
	mux.HandleFunc("/auth/signup", authHandler.SignUp)
	mux.HandleFunc("/auth/signin", authHandler.SignIn)
	mux.HandleFunc("/auth/refresh_access", authHandler.RefreshAccessToken)
	mux.HandleFunc("/auth/refresh_tokens", authHandler.RefreshTokens)
	
	protected := http.NewServeMux()
	protected.HandleFunc("/history", gameHandler.GetHistoryByToken)
	protected.HandleFunc("/leaderboard", gameHandler.GetLeaderboard)
	protected.HandleFunc("/game", gameHandler.CreateGame)
	protected.HandleFunc("/game/join/", gameHandler.JoinGame)
	protected.HandleFunc("/waiting", gameHandler.GetWaitingGames)
	protected.HandleFunc("/user/", gameHandler.GetPublicUserInfoById)
	protected.HandleFunc("/user/info_by_access_token", gameHandler.GetPublicUserInfoByToken)
	protected.HandleFunc("/game/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
				gameHandler.GetGame(w, r)
		case "POST":
				gameHandler.MakeMove(w, r)
			
		}
	})
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")		

		switch r.URL.Path {
		case "/auth/signup", "/auth/signin", "/auth/refresh_access", "/auth/refresh_tokens":
			mux.ServeHTTP(w, r)
		default:
			authMiddleware.Middleware(protected).ServeHTTP(w, r)
		}
	})
}