package middleware

import (
	"context"
	"encoding/json"
	"hr-program/shared/auth"
	"hr-program/shared/config"
	db "hr-program/shared/connection"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		claims, err := auth.ParseAccessToken(config.AppConfig.JWTSecret, parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		key := "auth:session:" + claims.JTI
		raw, err := db.ConnectRedis().Get(context.Background(), key).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "session expired or revoked",
			})
			return
		}

		var session auth.SessionData
		if err := json.Unmarshal([]byte(raw), &session); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid session",
			})
			return
		}

		if session.Revoked {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "session revoked",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("employee_id", claims.EmployeeID)
		c.Set("jti", claims.JTI)
		c.Next()
	}
}
