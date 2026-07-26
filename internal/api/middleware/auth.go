package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/restaurantautomation/api/internal/config"
)

type authContextKey string

const claimsKey authContextKey = "authClaims"

type Claims struct {
	Subject string `json:"sub"`
	Email   string `json:"email"`
	Role    string `json:"role"`
}

func Auth(cfg config.AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if !cfg.Required {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			tokenString, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || strings.TrimSpace(tokenString) == "" {
				unauthorized(w)
				return
			}
			claims, err := parseClaims(tokenString, cfg)
			if err != nil {
				unauthorized(w)
				return
			}
			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseClaims(tokenString string, cfg config.AuthConfig) (Claims, error) {
	options := []jwt.ParserOption{jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()})}
	if cfg.Issuer != "" {
		options = append(options, jwt.WithIssuer(cfg.Issuer))
	}
	token, err := jwt.NewParser(options...).Parse(tokenString, func(token *jwt.Token) (any, error) { return []byte(cfg.SupabaseJWTSecret), nil })
	if err != nil || !token.Valid {
		return Claims{}, jwt.ErrTokenInvalidClaims
	}
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, jwt.ErrTokenInvalidClaims
	}
	claims := Claims{Subject: stringClaim(mapClaims, "sub"), Email: stringClaim(mapClaims, "email"), Role: stringClaim(mapClaims, "role")}
	return claims, nil
}

func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(Claims)
	return claims, ok
}
func stringClaim(claims jwt.MapClaims, key string) string {
	value, _ := claims[key].(string)
	return value
}
func unauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "authentication required"})
}
