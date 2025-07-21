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

const BearerPrefix = "Bearer "

func OptionalAuthMiddleware(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Printf("[middleware:auth] OptionalAuth called: path=%s method=%s",
			ctx.Request.URL.Path, ctx.Request.Method,
		)

		if err := parseAndSetClaims(ctx, jwtManager); err != nil {
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

func RequireAuthMiddleware(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Printf("[middleware:auth] RequireAuth called: path=%s method=%s",
			ctx.Request.URL.Path, ctx.Request.Method,
		)

		if err := parseAndSetClaims(ctx, jwtManager); err != nil {
			log.Printf("[middleware:auth][ERROR] RequireAuth failed: %v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		log.Printf("[middleware:auth] RequireAuth: authenticated userId=%d", ctx.GetInt("userId"))
		ctx.Next()
	}
}

func parseAndSetClaims(ctx *gin.Context, jwtManager *jwt.Manager) error {
	auth := ctx.GetHeader("Authorization")
	if auth == "" {
		return fmt.Errorf("no Authorization header")
	}

	if !strings.HasPrefix(auth, BearerPrefix) {
		return fmt.Errorf("invalid Authorization format")
	}

	token := strings.TrimPrefix(auth, BearerPrefix)
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
