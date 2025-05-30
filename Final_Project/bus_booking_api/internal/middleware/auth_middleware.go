package middleware

import (
	"bus_booking_api/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header is not provided")
			c.Abort()
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != authorizationTypeBearer {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unsupported authorization type: "+authType)
			c.Abort()
			return
		}

		accessToken := fields[1]
		claims, err := utils.ValidateToken(accessToken)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token: "+err.Error())
			c.Abort()
			return
		}

		c.Set(authorizationPayloadKey, claims)
		c.Next()
	}
}

func GetAuthPayload(c *gin.Context) (*utils.JWTCustomClaims, bool) {
	payload, exists := c.Get(authorizationPayloadKey)
	if !exists {
		return nil, false
	}
	claims, ok := payload.(*utils.JWTCustomClaims)
	return claims, ok
}
