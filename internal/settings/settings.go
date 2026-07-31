package settings

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Settings struct {
	RestaurantName     string    `json:"restaurant_name"`
	Timezone           string    `json:"timezone"`
	OperationalAlerts  bool      `json:"operational_alerts"`
	AutoAdvanceTickets bool      `json:"auto_advance_tickets"`
	UpdatedAt          time.Time `json:"updated_at"`
}

var ErrInvalidSettings = errors.New("invalid restaurant settings")

type Repository interface {
	Get(context.Context) (Settings, error)
	Save(context.Context, Settings) error
}
type Service struct{ repository Repository }

func NewService(repository Repository) *Service              { return &Service{repository: repository} }
func (s *Service) Get(ctx context.Context) (Settings, error) { return s.repository.Get(ctx) }
func (s *Service) Save(ctx context.Context, value Settings) (Settings, error) {
	value.RestaurantName = strings.TrimSpace(value.RestaurantName)
	value.Timezone = strings.TrimSpace(value.Timezone)
	if value.RestaurantName == "" || value.Timezone == "" {
		return Settings{}, ErrInvalidSettings
	}
	value.UpdatedAt = time.Now().UTC()
	return value, s.repository.Save(ctx, value)
}
func Default() Settings {
	return Settings{RestaurantName: "Restaurant Automation", Timezone: "Asia/Kolkata", OperationalAlerts: true}
}

type MemoryRepository struct {
	mu    sync.RWMutex
	value Settings
}

func NewMemoryRepository() *MemoryRepository { return &MemoryRepository{value: Default()} }
func (r *MemoryRepository) Get(_ context.Context) (Settings, error) {
	r.mu.RLock()
	value := r.value
	r.mu.RUnlock()
	return value, nil
}
func (r *MemoryRepository) Save(_ context.Context, value Settings) error {
	r.mu.Lock()
	r.value = value
	r.mu.Unlock()
	return nil
}

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}
func (r *PostgresRepository) Get(ctx context.Context) (Settings, error) {
	var value Settings
	err := r.pool.QueryRow(ctx, `SELECT restaurant_name,timezone,operational_alerts,auto_advance_tickets,updated_at FROM restaurant_settings WHERE id='current'`).Scan(&value.RestaurantName, &value.Timezone, &value.OperationalAlerts, &value.AutoAdvanceTickets, &value.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Default(), nil
	}
	return value, err
}
func (r *PostgresRepository) Save(ctx context.Context, value Settings) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO restaurant_settings (id,restaurant_name,timezone,operational_alerts,auto_advance_tickets,updated_at) VALUES ('current',$1,$2,$3,$4,$5) ON CONFLICT (id) DO UPDATE SET restaurant_name=EXCLUDED.restaurant_name,timezone=EXCLUDED.timezone,operational_alerts=EXCLUDED.operational_alerts,auto_advance_tickets=EXCLUDED.auto_advance_tickets,updated_at=EXCLUDED.updated_at`, value.RestaurantName, value.Timezone, value.OperationalAlerts, value.AutoAdvanceTickets, value.UpdatedAt)
	return err
}
