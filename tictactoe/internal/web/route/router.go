package route

import (
	"net/http"
	"strings"
	"tictactoe/internal/web/module"
	"tictactoe/internal/web/mapper"
	"fmt"
)

func NewRouter(handler *module.GameHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
			
		switch {
		case r.URL.Path == "/health" && r.Method == http.MethodGet:
			handler.HealthCheck(w, r)
			
		case r.URL.Path == "/game" && r.Method == http.MethodPost:
			handler.CreateGame(w, r)
			
		case strings.HasPrefix(r.URL.Path, "/game/") && r.Method == http.MethodGet:
			handler.GetGame(w, r)
			
		case strings.HasPrefix(r.URL.Path, "/game/") && r.Method == http.MethodPost:
			handler.MakeMove(w, r)
			
		default:
			mapper.WriteJSON(w, http.StatusNotFound, 
				mapper.ToErrorResponse(fmt.Errorf("endpoint not found")))
		}
	})
}
