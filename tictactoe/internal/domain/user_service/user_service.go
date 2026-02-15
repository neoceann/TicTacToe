package user_service

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"

    "github.com/google/uuid"
    "tictactoe/internal/datasource/model"
    "tictactoe/internal/datasource/repository"
    "tictactoe/internal/web/jwt"
)

type UserService struct {
    storage *repository.GameStorage
}

func NewUserService(storage *repository.GameStorage) *UserService {
    return &UserService{storage: storage}
}

func (s *UserService) hashPassword(password string) string {
    hash := sha256.Sum256([]byte(password))
    return hex.EncodeToString(hash[:])
}

func (s *UserService) Register(req *_jwt.JwtRequest) (bool, error) {
    _, err := s.storage.FindUserByLogin(req.Login)
    if err == nil {
        return false, fmt.Errorf("user already exists")
    }

    user := &model.User{
        ID:           uuid.New(),
        Login:        req.Login,
        PasswordHash: s.hashPassword(req.Password),
    }

    err = s.storage.SaveUser(user)
    return err == nil, err
}

func (s *UserService) Authenticate(login, password string) (uuid.UUID, error) {
    user, err := s.storage.FindUserByLogin(login)
    if err != nil {
        return uuid.Nil, fmt.Errorf("user not found")
    }

    if s.hashPassword(password) != user.PasswordHash {
        return uuid.Nil, fmt.Errorf("invalid password")
    }

    return user.ID, nil
}