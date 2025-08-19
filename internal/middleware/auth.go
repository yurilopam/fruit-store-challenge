package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthConfig struct{ Secret []byte }

type ContextKeys string

const (
	CtxUserID ContextKeys = "user_id"
	CtxRole   ContextKeys = "role"
)

func AuthRequired(cfg AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || len(h) < 8 || h[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tokenStr := h[7:]
		t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return cfg.Secret, nil
		})
		if err != nil || !t.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}
		c.Set(string(CtxUserID), claims["user_id"])
		c.Set(string(CtxRole), claims["role"])
		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get(string(CtxRole))
		if !ok || v.(string) != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
