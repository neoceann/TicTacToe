package module

import (
	"encoding/json"
	"fmt"
	"net/http"

	authservice "tictactoe/internal/domain/auth_service"
	_jwt "tictactoe/internal/web/jwt"
	"tictactoe/internal/web/mapper"
)

type AuthHandler struct {
	authService authservice.AuthService
}

func NewAuthHandler(authService authservice.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (s *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {

	var req _jwt.JwtRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapper.WriteJSON(w, http.StatusBadRequest, 
			mapper.ToErrorResponse(fmt.Errorf("invalid request body")))
		return
	}

	_, err := s.authService.SignUp(r.Context(), &req)
	if err != nil {
		mapper.WriteJSON(w, http.StatusBadRequest, 
			mapper.ToErrorResponse(err))
		return
	}

	response := map[string]interface{}{
		"status": "User registered successfully",
	}
	mapper.WriteJSON(w, http.StatusCreated, response)
}

func (s *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request){
    var req _jwt.JwtRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        mapper.WriteJSON(w, http.StatusBadRequest,
            mapper.ToErrorResponse(fmt.Errorf("invalid request body")))
        return
    }
    
    resp, err := s.authService.SignIn(r.Context(), &req)
    if err != nil {
        mapper.WriteJSON(w, http.StatusUnauthorized,
            mapper.ToErrorResponse(fmt.Errorf("singin error: %w", err)))
        return
    }
    
    mapper.WriteJSON(w, http.StatusOK, resp)
}

func (s *AuthHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
    var req _jwt.RefreshJwtRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        mapper.WriteJSON(w, http.StatusBadRequest,
            mapper.ToErrorResponse(fmt.Errorf("invalid request body")))
        return
    }
    
    resp, err := s.authService.RefreshAccessToken(r.Context(), req.RefreshToken)
    if err != nil {
        mapper.WriteJSON(w, http.StatusUnauthorized,
            mapper.ToErrorResponse(fmt.Errorf("refresh accsess token error: %w", err)))
        return
    }
    
    mapper.WriteJSON(w, http.StatusOK, resp)
}

func (s *AuthHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
    var req _jwt.RefreshJwtRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        mapper.WriteJSON(w, http.StatusBadRequest,
            mapper.ToErrorResponse(fmt.Errorf("invalid request body")))
        return
    }
    
    resp, err := s.authService.RefreshTokens(r.Context(), req.RefreshToken)
    if err != nil {
        mapper.WriteJSON(w, http.StatusUnauthorized,
            mapper.ToErrorResponse(fmt.Errorf("refresh tokens error: %w", err)))
        return
    }
    
    mapper.WriteJSON(w, http.StatusOK, resp)
}