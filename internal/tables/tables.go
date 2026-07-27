package tables

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Table struct {
	ID          string    `json:"id"`
	Seats       int       `json:"seats"`
	Guest       string    `json:"guest"`
	Reservation string    `json:"reservation"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateInput struct {
	ID          string `json:"id"`
	Seats       int    `json:"seats"`
	Guest       string `json:"guest"`
	Reservation string `json:"reservation"`
	Status      string `json:"status"`
}

var ErrNotFound = errors.New("table not found")
var ErrInvalidTable = errors.New("invalid table")
var ErrInvalidStatus = errors.New("invalid table status")

type Repository interface {
	Create(context.Context, Table) error
	List(context.Context) ([]Table, error)
	UpdateStatus(context.Context, string, string) (Table, error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) Create(ctx context.Context, input CreateInput) (Table, error) {
	input.ID = strings.TrimSpace(input.ID)
	if input.ID == "" || input.Seats < 1 {
		return Table{}, ErrInvalidTable
	}
	status := input.Status
	if status == "" {
		status = "available"
	}
	if !validStatus(status) {
		return Table{}, ErrInvalidStatus
	}
	now := time.Now().UTC()
	table := Table{ID: input.ID, Seats: input.Seats, Guest: strings.TrimSpace(input.Guest), Reservation: strings.TrimSpace(input.Reservation), Status: status, CreatedAt: now, UpdatedAt: now}
	return table, s.repository.Create(ctx, table)
}

func (s *Service) List(ctx context.Context) ([]Table, error) { return s.repository.List(ctx) }

func (s *Service) UpdateStatus(ctx context.Context, id, status string) (Table, error) {
	if !validStatus(status) {
		return Table{}, ErrInvalidStatus
	}
	return s.repository.UpdateStatus(ctx, id, status)
}

func validStatus(status string) bool {
	return status == "available" || status == "seated" || status == "reserved" || status == "needs_check"
}

type MemoryRepository struct {
	mu     sync.RWMutex
	tables map[string]Table
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{tables: make(map[string]Table)}
}

func (r *MemoryRepository) Create(_ context.Context, table Table) error {
	r.mu.Lock()
	r.tables[table.ID] = table
	r.mu.Unlock()
	return nil
}

func (r *MemoryRepository) List(_ context.Context) ([]Table, error) {
	r.mu.RLock()
	result := make([]Table, 0, len(r.tables))
	for _, table := range r.tables {
		result = append(result, table)
	}
	r.mu.RUnlock()
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (r *MemoryRepository) UpdateStatus(_ context.Context, id, status string) (Table, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	table, ok := r.tables[id]
	if !ok {
		return Table{}, ErrNotFound
	}
	table.Status = status
	table.UpdatedAt = time.Now().UTC()
	r.tables[id] = table
	return table, nil
}

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, table Table) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO restaurant_tables (id, seats, guest, reservation, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`, table.ID, table.Seats, table.Guest, table.Reservation, table.Status, table.CreatedAt, table.UpdatedAt)
	return err
}

func (r *PostgresRepository) List(ctx context.Context) ([]Table, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, seats, guest, reservation, status, created_at, updated_at FROM restaurant_tables ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []Table{}
	for rows.Next() {
		var table Table
		if err := rows.Scan(&table.ID, &table.Seats, &table.Guest, &table.Reservation, &table.Status, &table.CreatedAt, &table.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, table)
	}
	return result, rows.Err()
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id, status string) (Table, error) {
	var table Table
	err := r.pool.QueryRow(ctx, `UPDATE restaurant_tables SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING id, seats, guest, reservation, status, created_at, updated_at`, id, status).Scan(&table.ID, &table.Seats, &table.Guest, &table.Reservation, &table.Status, &table.CreatedAt, &table.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Table{}, ErrNotFound
	}
	return table, err
}
