package authservice

import (
	"context"
	"tictactoe/internal/domain/user_service"
	_jwt "tictactoe/internal/web/jwt"
	"fmt"
)

type AuthServiceImpl struct {
	userService *user_service.UserService
	jwtProvider *_jwt.JwtProvider
}

func NewAuthService(userService *user_service.UserService, jwtProvider *_jwt.JwtProvider) AuthService {
    return &AuthServiceImpl{
        userService: userService,
        jwtProvider: jwtProvider,
    }
}

func (h *AuthServiceImpl) SignUp(ctx context.Context, req *_jwt.JwtRequest) (bool, error) {
	return h.userService.Register(req)
}

func (s *AuthServiceImpl) SignIn(ctx context.Context, req *_jwt.JwtRequest) (*_jwt.JwtResponse, error) {
	userID, err := s.userService.Authenticate(req.Login, req.Password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	accessToken, err := s.jwtProvider.GenerateAccessToken(userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtProvider.GenerateRefreshToken(userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &_jwt.JwtResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServiceImpl) RefreshAccessToken(ctx context.Context, refreshToken string) (*_jwt.JwtResponse, error) {
	userID, err := s.jwtProvider.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	newAccessToken, err := s.jwtProvider.GenerateAccessToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	return &_jwt.JwtResponse{
		Type:         "Bearer",
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServiceImpl) RefreshTokens(ctx context.Context, refreshToken string) (*_jwt.JwtResponse, error) {
	userID, err := s.jwtProvider.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	newAccessToken, err := s.jwtProvider.GenerateAccessToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	newRefreshToken, err := s.jwtProvider.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	return &_jwt.JwtResponse{
		Type:         "Bearer",
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}