package authservice

import (
	"context"
	_jwt "tictactoe/internal/web/jwt"
)

type AuthService interface {
	SignUp(ctx context.Context, req *_jwt.JwtRequest) (bool, error)

    SignIn(ctx context.Context, req *_jwt.JwtRequest) (*_jwt.JwtResponse, error)
    
    RefreshAccessToken(ctx context.Context, refreshToken string) (*_jwt.JwtResponse, error)
    
    RefreshTokens(ctx context.Context, refreshToken string) (*_jwt.JwtResponse, error)
}