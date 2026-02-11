package module

import (
	"context"
	"net/http"
	"fmt"
	
	"tictactoe/internal/domain/user_service"
	"tictactoe/internal/web/mapper"
)

type UserAuthenticator struct {
	userService *user_service.UserService
}

func NewUserAuthenticator(userService *user_service.UserService) *UserAuthenticator {
	return &UserAuthenticator{userService: userService}
}

func (a *UserAuthenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, password, ok := r.BasicAuth()
		if !ok {
			mapper.WriteJSON(w, http.StatusUnauthorized,
				mapper.ToErrorResponse(fmt.Errorf("authorization required")))
			return
		}
		
		if login == "" || password == "" {
			mapper.WriteJSON(w, http.StatusUnauthorized,
				mapper.ToErrorResponse(fmt.Errorf("login and password required")))
			return
		}
		
		userID, err := a.userService.Authenticate(login, password)
		if err != nil {
			mapper.WriteJSON(w, http.StatusUnauthorized,
				mapper.ToErrorResponse(fmt.Errorf("invalid credentials")))
			return
		}
		
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}