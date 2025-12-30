package module

import "net/http"

type GameHandlerInterface interface {
	MakeMove(w http.ResponseWriter, r *http.Request)
	
	CreateGame(w http.ResponseWriter, r *http.Request)
	
	GetGame(w http.ResponseWriter, r *http.Request)
}