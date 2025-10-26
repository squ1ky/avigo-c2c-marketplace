package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

const (
	HeaderUserID    = "X-User-ID"
	HeaderUserEmail = "X-User-Email"
	HeaderUserRole  = "X-User-Role"
)

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func JWTMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve token from cookie or Authorization header
		tokenString := extractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authentication token",
				"code":  "MISSING_TOKEN",
			})
			return
		}

		// Parse and validate token
		claims := &TokenClaims{}
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(jwtSecret), nil
			},
		)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
				"code":  "INVALID_TOKEN",
			})
			return
		}

		// Clear up the headers that the client may have passed on
		c.Request.Header.Del(HeaderUserID)
		c.Request.Header.Del(HeaderUserEmail)
		c.Request.Header.Del(HeaderUserRole)

		// Add verified user's data
		c.Request.Header.Set(HeaderUserID, claims.UserID)
		c.Request.Header.Set(HeaderUserEmail, claims.Email)
		c.Request.Header.Set(HeaderUserRole, claims.Role)

		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	if cookie, err := c.Cookie("access_token"); err == nil && cookie != "" {
		return cookie
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	if !strings.HasPrefix(authHeader, "Bearer") {
		return ""
	}

	return strings.TrimPrefix(authHeader, "Bearer ")
}
