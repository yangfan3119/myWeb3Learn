package utils

import (
	"go01_4/config"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Password string `json:"password"`
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(username string, password string, userId uint) (string, error) {
	expirationTime := time.Now().Add(config.Cfg.ExpireHours * time.Hour)
	claims := &Claims{
		Username: username,
		Password: password,
		UserID:   strconv.FormatUint(uint64(userId), 10),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := []byte(config.Cfg.JwtSecret)
	return token.SignedString(jwtSecret)
}

func GetJwtClaimsUsername(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		jwtSecret := []byte(config.Cfg.JwtSecret)
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		mlog.Warn("Token无效")
		return "", err
	}
	return claims.UserID, nil
}
