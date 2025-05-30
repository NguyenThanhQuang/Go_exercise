package utils

import (
	"bus_booking_api/internal/config"
	"bus_booking_api/internal/model"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTCustomClaims struct {
	UserID    string        `json:"userId"`
	Email     string        `json:"email"`
	Role      model.UserRole `json:"role"`
	CompanyID string        `json:"companyId,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(user *model.User) (string, error) {
	claims := &JWTCustomClaims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AppConfig.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "bus_booking_api",
			Subject:   user.ID.Hex(),
		},
	}

	if user.Role == model.RoleCompanyAdmin && !user.CompanyID.IsZero() {
		claims.CompanyID = user.CompanyID.Hex()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.AppConfig.JWTSecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateToken(tokenString string) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.JWTSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GetUserIDFromToken(claims *JWTCustomClaims) (primitive.ObjectID, error) {
	if claims == nil {
		return primitive.NilObjectID, errors.New("claims are nil")
	}
	return primitive.ObjectIDFromHex(claims.UserID)
}
