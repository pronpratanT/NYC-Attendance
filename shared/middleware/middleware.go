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

func CORS() gin.HandlerFunc {
	allowedOrigins := make(map[string]struct{}, len(config.AppConfig.CORSAllowedOrigins))
	allowAllOrigins := false
	for _, origin := range config.AppConfig.CORSAllowedOrigins {
		if origin == "*" {
			allowAllOrigins = true
			break
		}
		allowedOrigins[origin] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		headers := c.Writer.Header()
		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")
		headers.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		headers.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		headers.Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		if allowAllOrigins {
			headers.Set("Access-Control-Allow-Origin", "*")
		} else {
			if _, ok := allowedOrigins[origin]; !ok {
				if c.Request.Method == http.MethodOptions {
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "origin not allowed"})
					return
				}
				c.Next()
				return
			}

			headers.Set("Access-Control-Allow-Origin", origin)
			headers.Set("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow preflight OPTIONS to pass through for CORS
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

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
