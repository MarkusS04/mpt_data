package auth

import (
	"errors"
	"fmt"
	"mpt_data/helper"
	dbModel "mpt_data/models/dbmodel"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	ExpiresIn = time.Hour * 2
)

func GetUserIDFromToken(tokenString string) (uint, error) {
	token, err := getTokenFromString(strings.TrimPrefix(tokenString, "Bearer "))
	if err != nil {
		return 0, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("error reading token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id not found in token claims")
	}

	userID := uint(userIDFloat)
	return userID, nil
}

func generateJWT(user dbModel.User) (string, error) {

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(ExpiresIn).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(helper.Config.API.JWTKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := getTokenFromString(tokenString)

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid Token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("error reading token claims")
	}

	return claims, nil
}

func getTokenFromString(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(helper.Config.API.JWTKey), nil
		})
}
