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

// UserObserver records an already-authenticated identity in the local
// application directory. It deliberately has no access to Supabase Admin APIs.
type UserObserver interface {
	Observe(context.Context, string, string) error
}

func TrackUser(observer UserObserver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if observer == nil {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				next.ServeHTTP(w, r)
				return
			}
			if err := observer.Observe(r.Context(), claims.Subject, claims.Email); err != nil {
				http.Error(w, "user directory unavailable", http.StatusServiceUnavailable)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// PermissionChecker resolves the locally assigned permissions for an
// authenticated subject. It is defined at this boundary to avoid coupling the
// HTTP layer to a particular identity provider.
type PermissionChecker interface {
	Permissions(context.Context, string) ([]string, error)
}

func RequirePermission(checker PermissionChecker, required string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				// AUTH_REQUIRED=false is an explicit local-development mode.
				next.ServeHTTP(w, r)
				return
			}
			if checker == nil {
				http.Error(w, "authorization unavailable", http.StatusServiceUnavailable)
				return
			}
			permissions, err := checker.Permissions(r.Context(), claims.Subject)
			if err != nil {
				http.Error(w, "authorization unavailable", http.StatusServiceUnavailable)
				return
			}
			for _, permission := range permissions {
				if permission == required {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "permission required", http.StatusForbidden)
			return
		})
	}
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
