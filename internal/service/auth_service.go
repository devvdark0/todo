package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/devvdark0/todo/internal/auth"
	"github.com/devvdark0/todo/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
	ErrEmailInUse         = errors.New("email already in use")
)

type UserStorage interface {
	Create(user model.User) error
	GetByID(id uuid.UUID) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
}

type JWTService struct {
	secret    []byte
	tokenTTL  time.Duration
	userStore UserStorage
}

func NewJWTService(
	secret []byte, tokenTTL time.Duration, userStore UserStorage,
) *JWTService {
	return &JWTService{
		secret:    secret,
		tokenTTL:  tokenTTL,
		userStore: userStore,
	}
}

func (j *JWTService) Register(email, username, password string) error {
	user, err := j.userStore.GetByEmail(email)
	if err != nil {
		return ErrEmailInUse
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	*user = model.User{
		ID:       uuid.New(),
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	if err = j.userStore.Create(*user); err != nil {
		return err
	}

	return nil
}

func (j *JWTService) Login(email, password string) (string, error) {
	user, err := j.userStore.GetByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err = auth.VerifyPassword(user.Password, password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := j.generateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *JWTService) generateToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(j.tokenTTL)

	claims := jwt.MapClaims{
		"sub":      user.ID.String(),
		"username": user.Username,
		"email":    user.Email,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("token generation err: %w", err)
	}

	return tokenString, nil
}

func (j *JWTService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return j.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}

		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
