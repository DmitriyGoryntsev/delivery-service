package utils

import (
	"time"
	"user-service/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	SecretKey              string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
}

func NewJWTManager(secretKey string, accessExpiry, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		SecretKey:              secretKey,
		AccessTokenExpiration:  accessExpiry,
		RefreshTokenExpiration: refreshExpiry,
	}
}

func (j *JWTManager) GenerateAccessToken(user *models.User) (string, error) {
	claims := &models.AccessTokenClaims{
		UserID:    user.ID.String(),
		Email:     user.Email,
		Role:      ifRole(user),
		IsCourier: user.IsCourier,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.AccessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWTManager) GenerateRefreshToken(user *models.User) (string, error) {
	claims := &models.RefreshTokenClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.RefreshTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWTManager) VerifyAccessToken(tokenStr string) (*models.AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil || token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*models.AccessTokenClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}

func ifRole(user *models.User) string {
	if user.IsCourier {
		return "courier"
	}
	return "user"
}
