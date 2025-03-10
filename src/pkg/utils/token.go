package utils

import (
	"ecommerce-price-tracker/internal/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type TokenConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
	Issuer             string
}

var config = TokenConfig{
	AccessTokenSecret:  "WebProject7Secret",
	RefreshTokenSecret: "Secret@WebProj",
	AccessTokenTTL:     15 * time.Minute,
	RefreshTokenTTL:    7 * 24 * time.Hour,
	Issuer:             "ecom-tracker",
}

type CustomClaims struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email,omitempty"`
	Roles     string    `json:"roles,omitempty"`
	TokenType TokenType `json:"token_type"`
	jwt.StandardClaims
}

func CreateToken(userID string, email string, role models.Role, tokenType TokenType) (string, error) {
	var ttl time.Duration
	var secret string

	switch tokenType {
	case AccessToken:
		ttl = config.AccessTokenTTL
		secret = config.AccessTokenSecret
	case RefreshToken:
		ttl = config.RefreshTokenTTL
		secret = config.RefreshTokenSecret
	default:
		return "", errors.New("invalid token type")
	}

	claims := CustomClaims{
		UserID:    userID,
		TokenType: tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    config.Issuer,
		},
	}

	claims.Email = email
	claims.Roles = string(role)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, tokenType TokenType) (*CustomClaims, error) {
	var secret string

	switch tokenType {
	case AccessToken:
		secret = config.AccessTokenSecret
	case RefreshToken:
		secret = config.RefreshTokenSecret
	default:
		return nil, errors.New("invalid token type")
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenType != tokenType {
		return nil, errors.New("invalid token type")
	}
	return claims, nil
}
