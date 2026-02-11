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
	
	protected := http.NewServeMux()
	protected.HandleFunc("/game", gameHandler.CreateGame)
	protected.HandleFunc("/game/join/", gameHandler.JoinGame)
	protected.HandleFunc("/waiting", gameHandler.GetWaitingGames)
	protected.HandleFunc("/user/", gameHandler.GetPublicUserInfoById)
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

		if r.URL.Path == "/auth/signup" || r.URL.Path == "/auth/signin" {
			mux.ServeHTTP(w, r)
		} else {
			authMiddleware.Middleware(protected).ServeHTTP(w, r)
		}
	})
}