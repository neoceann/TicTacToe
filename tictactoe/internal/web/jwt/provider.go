package _jwt

import (
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtProvider struct {
	accessToken []byte
	accessTTL time.Duration
	refreshToken []byte
	refreshTTL time.Duration
}

func NewJWTProvider(accT, refT string, accTTL, refTTL time.Duration) *JwtProvider {
	return &JwtProvider{accessToken: []byte(accT),
						accessTTL: accTTL,
						refreshToken: []byte(refT),
						refreshTTL: refTTL,}
}

func (p *JwtProvider) GenerateAccessToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims {
		Subject: userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.accessTTL)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.accessToken)	
}

func (p *JwtProvider) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims {
		Subject: userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.refreshTTL)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.refreshToken)	
}

func (p *JwtProvider) ValidateAccessToken(accToken string) (string, error) {
    token, err := jwt.ParseWithClaims(accToken, &jwt.RegisteredClaims{},
		func (token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("bad singing method: %v", token.Header["alg"])
			}
			return p.accessToken, nil
		})
	
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", fmt.Errorf("invalid token")
}

func (p *JwtProvider) ValidateRefreshToken(refToken string) (string, error) {
    token, err := jwt.ParseWithClaims(refToken, &jwt.RegisteredClaims{},
		func (token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("bad singing method: %v", token.Header["alg"])
			}
			return p.refreshToken, nil
		})
	
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", fmt.Errorf("invalid token")}

func (p *JwtProvider) GetUserIDFromToken(token string) (uuid.UUID, error) {
    userIdAsString, err := p.ValidateAccessToken(token)
    if err != nil {
        userIdAsString, err = p.ValidateRefreshToken(token)
        if err != nil {
            return uuid.Nil, fmt.Errorf("invalid token: %w", err)
        }
    }

	userID, err := uuid.Parse(userIdAsString)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed parse uuid from token")
	}

    return userID, nil
}