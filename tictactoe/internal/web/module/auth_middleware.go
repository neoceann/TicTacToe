package module

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	_jwt "tictactoe/internal/web/jwt"
	"tictactoe/internal/web/mapper"
)

type UserAuthenticator struct {
	jwtProvider *_jwt.JwtProvider
}

func NewUserAuthenticator(jwtProvider *_jwt.JwtProvider) *UserAuthenticator {
	return &UserAuthenticator{
		jwtProvider: jwtProvider,
	}
}

func (a *UserAuthenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			mapper.WriteJSON(w, http.StatusUnauthorized,
				mapper.ToErrorResponse(fmt.Errorf("authorization header required")))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			mapper.WriteJSON(w, http.StatusUnauthorized,
				mapper.ToErrorResponse(fmt.Errorf("invalid authorization format, use Bearer <token>")))
			return
		}

		accessToken := parts[1]

		userID, err := a.jwtProvider.ValidateAccessToken(accessToken)
		if err != nil {
			mapper.WriteJSON(w, http.StatusUnauthorized,
				mapper.ToErrorResponse(fmt.Errorf("invalid or expired token")))
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}