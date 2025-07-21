package auth

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/1URose/marketplace/internal/common/jwt"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	BearerPrefix string
	jwtManager   *jwt.Manager
}

func NewMiddleware(bearerPrefix string, jwtManager *jwt.Manager) *Middleware {
	log.Println("[middleware:auth] NewMiddleware initialized")
	return &Middleware{
		BearerPrefix: bearerPrefix,
		jwtManager:   jwtManager,
	}
}

func (m *Middleware) Optional() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Printf("[middleware:auth] OptionalAuth called: path=%s method=%s",
			ctx.Request.URL.Path, ctx.Request.Method,
		)

		if err := m.parseAndSetClaims(ctx, m.jwtManager); err != nil {
			log.Printf("[middleware:auth][INFO] OptionalAuth: unauthenticated user (%v)", err)
			ctx.Set("isAuthenticated", false)
		} else {
			uid := ctx.GetInt("userId")
			log.Printf("[middleware:auth] OptionalAuth: authenticated userId=%d", uid)
			ctx.Set("isAuthenticated", true)
		}

		ctx.Next()
	}
}

func (m *Middleware) Require() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Printf("[middleware:auth] RequireAuth called: path=%s method=%s",
			ctx.Request.URL.Path, ctx.Request.Method,
		)

		if err := m.parseAndSetClaims(ctx, m.jwtManager); err != nil {
			log.Printf("[middleware:auth][ERROR] RequireAuth failed: %v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		log.Printf("[middleware:auth] RequireAuth: authenticated userId=%d", ctx.GetInt("userId"))
		ctx.Next()
	}
}

func (m *Middleware) parseAndSetClaims(ctx *gin.Context, jwtManager *jwt.Manager) error {
	raw := ctx.GetHeader("Authorization")
	if raw == "" {
		return fmt.Errorf("no Authorization header")
	}

	var token string
	if strings.HasPrefix(raw, m.BearerPrefix) {
		token = strings.TrimPrefix(raw, m.BearerPrefix)
	} else {
		token = raw
		log.Printf("[middleware:auth] parseAndSetClaims: no Bearer prefix, using raw token")
	}

	claims, err := jwtManager.ValidateAccessToken(token)
	if err != nil {
		return err
	}

	if uid, err := strconv.Atoi(claims.Subject); err == nil {
		ctx.Set("userId", uid)
	}
	ctx.Set("userEmail", claims.Email)
	return nil
}
