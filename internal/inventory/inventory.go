package inventory

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Item struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	OnHand    float64   `json:"on_hand"`
	ParLevel  float64   `json:"par_level"`
	Unit      string    `json:"unit"`
	Supplier  string    `json:"supplier"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateInput struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	OnHand   float64 `json:"on_hand"`
	ParLevel float64 `json:"par_level"`
	Unit     string  `json:"unit"`
	Supplier string  `json:"supplier"`
	Status   string  `json:"status"`
}

var ErrNotFound = errors.New("inventory item not found")
var ErrInvalidItem = errors.New("invalid inventory item")
var ErrInvalidStatus = errors.New("invalid inventory status")

type Repository interface {
	Create(context.Context, Item) error
	List(context.Context) ([]Item, error)
	UpdateStock(context.Context, string, float64, string) (Item, error)
	UpdateStatus(context.Context, string, string) (Item, error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) Create(ctx context.Context, input CreateInput) (Item, error) {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Category) == "" || strings.TrimSpace(input.Unit) == "" || input.OnHand < 0 || input.ParLevel < 0 {
		return Item{}, ErrInvalidItem
	}
	status := input.Status
	if status == "" {
		status = stockStatus(input.OnHand, input.ParLevel)
	}
	if !validStatus(status) {
		return Item{}, ErrInvalidStatus
	}
	now := time.Now().UTC()
	item := Item{ID: uuid.NewString(), Name: strings.TrimSpace(input.Name), Category: strings.TrimSpace(input.Category), OnHand: input.OnHand, ParLevel: input.ParLevel, Unit: strings.TrimSpace(input.Unit), Supplier: strings.TrimSpace(input.Supplier), Status: status, CreatedAt: now, UpdatedAt: now}
	return item, s.repository.Create(ctx, item)
}

func (s *Service) List(ctx context.Context) ([]Item, error) { return s.repository.List(ctx) }

func (s *Service) UpdateStock(ctx context.Context, id string, onHand float64) (Item, error) {
	if onHand < 0 {
		return Item{}, ErrInvalidItem
	}
	items, err := s.repository.List(ctx)
	if err != nil {
		return Item{}, err
	}
	for _, item := range items {
		if item.ID == id {
			return s.repository.UpdateStock(ctx, id, onHand, stockStatus(onHand, item.ParLevel))
		}
	}
	return Item{}, ErrNotFound
}

func (s *Service) UpdateStatus(ctx context.Context, id, status string) (Item, error) {
	if !validStatus(status) {
		return Item{}, ErrInvalidStatus
	}
	return s.repository.UpdateStatus(ctx, id, status)
}

func stockStatus(onHand, parLevel float64) string {
	if onHand < parLevel {
		return "low_stock"
	}
	return "in_stock"
}

func validStatus(status string) bool {
	return status == "in_stock" || status == "low_stock" || status == "on_order"
}

type MemoryRepository struct {
	mu    sync.RWMutex
	items map[string]Item
}

func NewMemoryRepository() *MemoryRepository { return &MemoryRepository{items: make(map[string]Item)} }

func (r *MemoryRepository) Create(_ context.Context, item Item) error {
	r.mu.Lock()
	r.items[item.ID] = item
	r.mu.Unlock()
	return nil
}

func (r *MemoryRepository) List(_ context.Context) ([]Item, error) {
	r.mu.RLock()
	result := make([]Item, 0, len(r.items))
	for _, item := range r.items {
		result = append(result, item)
	}
	r.mu.RUnlock()
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result, nil
}

func (r *MemoryRepository) UpdateStock(_ context.Context, id string, onHand float64, status string) (Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.items[id]
	if !ok {
		return Item{}, ErrNotFound
	}
	item.OnHand = onHand
	item.Status = status
	item.UpdatedAt = time.Now().UTC()
	r.items[id] = item
	return item, nil
}

func (r *MemoryRepository) UpdateStatus(_ context.Context, id, status string) (Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.items[id]
	if !ok {
		return Item{}, ErrNotFound
	}
	item.Status = status
	item.UpdatedAt = time.Now().UTC()
	r.items[id] = item
	return item, nil
}

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, item Item) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO inventory_items (id, name, category, on_hand, par_level, unit, supplier, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, item.ID, item.Name, item.Category, item.OnHand, item.ParLevel, item.Unit, item.Supplier, item.Status, item.CreatedAt, item.UpdatedAt)
	return err
}

func (r *PostgresRepository) List(ctx context.Context) ([]Item, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, category, on_hand, par_level, unit, supplier, status, created_at, updated_at FROM inventory_items ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []Item{}
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.OnHand, &item.ParLevel, &item.Unit, &item.Supplier, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *PostgresRepository) UpdateStock(ctx context.Context, id string, onHand float64, status string) (Item, error) {
	return r.update(ctx, `UPDATE inventory_items SET on_hand = $2, status = $3, updated_at = NOW() WHERE id = $1 RETURNING id, name, category, on_hand, par_level, unit, supplier, status, created_at, updated_at`, id, onHand, status)
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id, status string) (Item, error) {
	return r.update(ctx, `UPDATE inventory_items SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING id, name, category, on_hand, par_level, unit, supplier, status, created_at, updated_at`, id, status)
}

func (r *PostgresRepository) update(ctx context.Context, query string, args ...any) (Item, error) {
	var item Item
	err := r.pool.QueryRow(ctx, query, args...).Scan(&item.ID, &item.Name, &item.Category, &item.OnHand, &item.ParLevel, &item.Unit, &item.Supplier, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Item{}, ErrNotFound
	}
	return item, err
}
