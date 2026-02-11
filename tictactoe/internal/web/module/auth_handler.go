package module

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"tictactoe/internal/domain/user_service"
	"tictactoe/internal/web/mapper"
	"tictactoe/internal/web/singup_request"
)

type AuthHandler struct {
	userService *user_service.UserService
}

func NewAuthHandler(userService *user_service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {

	var req singup_request.SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapper.WriteJSON(w, http.StatusBadRequest, 
			mapper.ToErrorResponse(fmt.Errorf("invalid request body")))
		return
	}

	_, err := h.userService.Register(&req)
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

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		mapper.WriteJSON(w, http.StatusUnauthorized,
			mapper.ToErrorResponse(fmt.Errorf("authorization header required")))
		return
	}

	if !strings.HasPrefix(authHeader, "Basic ") {
		mapper.WriteJSON(w, http.StatusUnauthorized,
			mapper.ToErrorResponse(fmt.Errorf("invalid authorization format")))
		return
	}

	encoded := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		mapper.WriteJSON(w, http.StatusUnauthorized,
			mapper.ToErrorResponse(fmt.Errorf("invalid credentials encoding")))
		return
	}

	// Парсим "login:password"
	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		mapper.WriteJSON(w, http.StatusUnauthorized,
			mapper.ToErrorResponse(fmt.Errorf("invalid credentials format")))
		return
	}

	login := credentials[0]
	password := credentials[1]

	userID, err := h.userService.Authenticate(login, password)
	if err != nil {
		mapper.WriteJSON(w, http.StatusUnauthorized,
			mapper.ToErrorResponse(fmt.Errorf("invalid credentials")))
		return
	}

	response := map[string]interface{}{
		"auth user_id": userID.String(),
	}
	mapper.WriteJSON(w, http.StatusOK, response)
}