// Package config loads and validates all application configuration from
// environment variables and optional .env files using Viper.
//
// Configuration is loaded once at startup and injected via dependency
// injection — no global state, no direct os.Getenv() calls in business logic.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config is the top-level configuration struct for the server.
// All fields are loaded from environment variables (or a .env file).
// Viper maps ENV_VAR_NAME → struct field via mapstructure tags.
type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	Printers     PrintersConfig
	Integrations IntegrationsConfig
	Auth         AuthConfig
	Log          LogConfig
}

// PrintersConfig holds optional physical network-printer settings.
type PrintersConfig struct {
	KitchenAddress string        `mapstructure:"PRINTER_KITCHEN_ADDRESS"`
	ConnectTimeout time.Duration `mapstructure:"PRINTER_CONNECT_TIMEOUT"`
}

// IntegrationsConfig enables the generic signed webhook provider.
type IntegrationsConfig struct {
	WebhookProviderName string `mapstructure:"WEBHOOK_PROVIDER_NAME"`
	WebhookSecret       string `mapstructure:"WEBHOOK_PROVIDER_SECRET"`
	SwiggyWebhookSecret string `mapstructure:"SWIGGY_WEBHOOK_SECRET"`
}

// ServerConfig holds HTTP server tuning parameters.
type ServerConfig struct {
	// Port the HTTP server listens on. Default: 8080.
	Port int `mapstructure:"PORT"`
	// ShutdownTimeout is the maximum time allowed for in-flight requests to
	// complete during graceful shutdown. Default: 30s.
	ShutdownTimeout time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
	// CORSOrigins is a comma-separated list of allowed CORS origins.
	// Use "*" only in development.
	CORSOrigins []string `mapstructure:"CORS_ORIGINS"`
}

// DatabaseConfig holds PostgreSQL connection parameters.
type DatabaseConfig struct {
	// DSN is the full PostgreSQL connection string.
	// Example: postgres://user:pass@localhost:5432/restaurant?sslmode=disable
	DSN string `mapstructure:"DATABASE_URL"`
	// MaxConns is the maximum number of pool connections. Default: 25.
	MaxConns int32 `mapstructure:"DB_MAX_CONNS"`
	// MinConns is the minimum number of pool connections. Default: 5.
	MinConns int32 `mapstructure:"DB_MIN_CONNS"`
	// MigrationsDir is the path to versioned SQL migrations. Default: migrations.
	MigrationsDir string `mapstructure:"MIGRATIONS_DIR"`
}

// AuthConfig holds authentication parameters.
type AuthConfig struct {
	// SupabaseJWTSecret is used to verify Supabase-issued JWTs server-side.
	// Set this to the JWT secret from your Supabase project settings.
	SupabaseJWTSecret string `mapstructure:"SUPABASE_JWT_SECRET"`
	// Required enforces valid access tokens for all /api/v1 routes.
	Required bool `mapstructure:"AUTH_REQUIRED"`
	// Issuer optionally validates the iss claim on Supabase access tokens.
	Issuer string `mapstructure:"SUPABASE_JWT_ISSUER"`
	// AdminEmails is a comma-separated bootstrap allow-list for user administration.
	// Application roles are assigned in the database; this list protects the initial
	// administrative control plane before role-based request enforcement exists.
	AdminEmails []string `mapstructure:"ADMIN_EMAILS"`
}

// LogConfig controls the zerolog output format.
type LogConfig struct {
	// Level sets the minimum log level. One of: trace, debug, info, warn, error.
	// Default: info.
	Level string `mapstructure:"LOG_LEVEL"`
	// Pretty enables human-readable output in development. Set to false in prod.
	Pretty bool `mapstructure:"LOG_PRETTY"`
}

// Load reads configuration from environment variables and, if present, from a
// .env file in the working directory. Environment variables take precedence
// over .env file values. Returns an error if required values are missing or
// the config cannot be parsed.
func Load() (*Config, error) {
	v := viper.New()

	// --- Defaults ---
	v.SetDefault("PORT", 8080)
	v.SetDefault("SHUTDOWN_TIMEOUT", "30s")
	v.SetDefault("CORS_ORIGINS", "http://localhost:5173,http://localhost:5174")
	v.SetDefault("DB_MAX_CONNS", 25)
	v.SetDefault("DB_MIN_CONNS", 5)
	v.SetDefault("MIGRATIONS_DIR", "migrations")
	v.SetDefault("PRINTER_CONNECT_TIMEOUT", "3s")
	v.SetDefault("AUTH_REQUIRED", false)
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_PRETTY", false)

	// --- .env file (optional, ignored if absent) ---
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	_ = v.ReadInConfig() // intentionally ignore "not found" errors

	// --- Environment variables ---
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal failed: %w", err)
	}

	// Parse CORS_ORIGINS from comma-separated string if needed.
	if raw := v.GetString("CORS_ORIGINS"); raw != "" {
		cfg.Server.CORSOrigins = splitTrimmed(raw, ",")
	}
	if raw := v.GetString("ADMIN_EMAILS"); raw != "" {
		cfg.Auth.AdminEmails = splitTrimmed(raw, ",")
	}
	if cfg.Auth.Required && strings.TrimSpace(cfg.Auth.SupabaseJWTSecret) == "" {
		return nil, fmt.Errorf("config: SUPABASE_JWT_SECRET is required when AUTH_REQUIRED=true")
	}
	if (cfg.Integrations.WebhookProviderName == "") != (cfg.Integrations.WebhookSecret == "") {
		return nil, fmt.Errorf("config: WEBHOOK_PROVIDER_NAME and WEBHOOK_PROVIDER_SECRET must be set together")
	}

	return &cfg, nil
}

// splitTrimmed splits s by sep and trims whitespace from each element.
func splitTrimmed(s, sep string) []string {
	parts := strings.Split(s, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
