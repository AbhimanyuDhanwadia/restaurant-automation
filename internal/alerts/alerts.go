package alerts

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Alert struct {
	ID        string    `json:"id"`
	Label     string    `json:"label"`
	Detail    string    `json:"detail"`
	Severity  string    `json:"severity"`
	CreatedAt time.Time `json:"created_at"`
}
type CreateInput struct {
	Label    string `json:"label"`
	Detail   string `json:"detail"`
	Severity string `json:"severity"`
}

var ErrNotFound = errors.New("alert not found")
var ErrInvalidAlert = errors.New("invalid alert")
var ErrInvalidSeverity = errors.New("invalid alert severity")

type Repository interface {
	Create(context.Context, Alert) error
	ListActive(context.Context) ([]Alert, error)
	Acknowledge(context.Context, string) error
}
type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }
func (s *Service) Create(ctx context.Context, input CreateInput) (Alert, error) {
	if strings.TrimSpace(input.Label) == "" || strings.TrimSpace(input.Detail) == "" {
		return Alert{}, ErrInvalidAlert
	}
	if !validSeverity(input.Severity) {
		return Alert{}, ErrInvalidSeverity
	}
	alert := Alert{ID: uuid.NewString(), Label: strings.TrimSpace(input.Label), Detail: strings.TrimSpace(input.Detail), Severity: input.Severity, CreatedAt: time.Now().UTC()}
	return alert, s.repository.Create(ctx, alert)
}
func (s *Service) ListActive(ctx context.Context) ([]Alert, error) {
	return s.repository.ListActive(ctx)
}
func (s *Service) Acknowledge(ctx context.Context, id string) error {
	return s.repository.Acknowledge(ctx, id)
}
func validSeverity(value string) bool { return value == "high" || value == "medium" || value == "low" }

type MemoryRepository struct {
	mu     sync.RWMutex
	alerts map[string]Alert
}

func NewMemoryRepository() *MemoryRepository { return &MemoryRepository{alerts: map[string]Alert{}} }
func (r *MemoryRepository) Create(_ context.Context, alert Alert) error {
	r.mu.Lock()
	r.alerts[alert.ID] = alert
	r.mu.Unlock()
	return nil
}
func (r *MemoryRepository) ListActive(_ context.Context) ([]Alert, error) {
	r.mu.RLock()
	result := make([]Alert, 0, len(r.alerts))
	for _, alert := range r.alerts {
		result = append(result, alert)
	}
	r.mu.RUnlock()
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.After(result[j].CreatedAt) })
	return result, nil
}
func (r *MemoryRepository) Acknowledge(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.alerts[id]; !ok {
		return ErrNotFound
	}
	delete(r.alerts, id)
	return nil
}

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}
func (r *PostgresRepository) Create(ctx context.Context, alert Alert) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO operational_alerts (id,label,detail,severity,created_at) VALUES ($1,$2,$3,$4,$5)`, alert.ID, alert.Label, alert.Detail, alert.Severity, alert.CreatedAt)
	return err
}
func (r *PostgresRepository) ListActive(ctx context.Context) ([]Alert, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,label,detail,severity,created_at FROM operational_alerts WHERE acknowledged_at IS NULL ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []Alert{}
	for rows.Next() {
		var alert Alert
		if err := rows.Scan(&alert.ID, &alert.Label, &alert.Detail, &alert.Severity, &alert.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, alert)
	}
	return result, rows.Err()
}
func (r *PostgresRepository) Acknowledge(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `UPDATE operational_alerts SET acknowledged_at=NOW() WHERE id=$1 AND acknowledged_at IS NULL`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
