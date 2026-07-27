package orders

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Item struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
type Order struct {
	ID        string    `json:"id"`
	Channel   string    `json:"channel"`
	Status    string    `json:"status"`
	Notes     string    `json:"notes,omitempty"`
	Items     []Item    `json:"items"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type CreateInput struct {
	Channel string `json:"channel"`
	Notes   string `json:"notes"`
	Items   []Item `json:"items"`
}

var ErrNotFound = errors.New("order not found")
var ErrInvalidOrder = errors.New("invalid order")
var ErrInvalidStatus = errors.New("invalid order status")

type Repository interface {
	Create(context.Context, Order) error
	List(context.Context) ([]Order, error)
	UpdateStatus(context.Context, string, string) (Order, error)
}
type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }
func (s *Service) Create(ctx context.Context, input CreateInput) (Order, error) {
	if input.Channel == "" || len(input.Items) == 0 {
		return Order{}, ErrInvalidOrder
	}
	for _, item := range input.Items {
		if item.Name == "" || item.Quantity < 1 {
			return Order{}, ErrInvalidOrder
		}
	}
	now := time.Now().UTC()
	order := Order{ID: uuid.NewString(), Channel: input.Channel, Status: "received", Notes: input.Notes, Items: input.Items, CreatedAt: now, UpdatedAt: now}
	return order, s.repository.Create(ctx, order)
}
func (s *Service) List(ctx context.Context) ([]Order, error) { return s.repository.List(ctx) }
func (s *Service) UpdateStatus(ctx context.Context, id, status string) (Order, error) {
	if !validStatus(status) {
		return Order{}, ErrInvalidStatus
	}
	return s.repository.UpdateStatus(ctx, id, status)
}
func validStatus(status string) bool {
	return status == "received" || status == "preparing" || status == "ready" || status == "delivered" || status == "cancelled"
}

type MemoryRepository struct {
	mu     sync.RWMutex
	orders map[string]Order
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{orders: make(map[string]Order)}
}
func (r *MemoryRepository) Create(_ context.Context, order Order) error {
	r.mu.Lock()
	r.orders[order.ID] = order
	r.mu.Unlock()
	return nil
}
func (r *MemoryRepository) List(_ context.Context) ([]Order, error) {
	r.mu.RLock()
	result := make([]Order, 0, len(r.orders))
	for _, order := range r.orders {
		result = append(result, order)
	}
	r.mu.RUnlock()
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.After(result[j].CreatedAt) })
	return result, nil
}
func (r *MemoryRepository) UpdateStatus(_ context.Context, id, status string) (Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.orders[id]
	if !ok {
		return Order{}, ErrNotFound
	}
	order.Status = status
	order.UpdatedAt = time.Now().UTC()
	r.orders[id] = order
	return order, nil
}

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}
func (r *PostgresRepository) Create(ctx context.Context, order Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `INSERT INTO orders (id, channel, status, notes, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)`, order.ID, order.Channel, order.Status, order.Notes, order.CreatedAt, order.UpdatedAt); err != nil {
		return err
	}
	for _, item := range order.Items {
		if _, err = tx.Exec(ctx, `INSERT INTO order_items (order_id, name, quantity) VALUES ($1,$2,$3)`, order.ID, item.Name, item.Quantity); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
func (r *PostgresRepository) List(ctx context.Context) ([]Order, error) {
	rows, err := r.pool.Query(ctx, `SELECT o.id,o.channel,o.status,o.notes,o.created_at,o.updated_at,i.name,i.quantity FROM orders o JOIN order_items i ON i.order_id=o.id ORDER BY o.created_at DESC,i.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	byID := map[string]*Order{}
	ids := []string{}
	for rows.Next() {
		var id string
		var item Item
		var order Order
		if err := rows.Scan(&id, &order.Channel, &order.Status, &order.Notes, &order.CreatedAt, &order.UpdatedAt, &item.Name, &item.Quantity); err != nil {
			return nil, err
		}
		if byID[id] == nil {
			order.ID = id
			byID[id] = &order
			ids = append(ids, id)
		}
		byID[id].Items = append(byID[id].Items, item)
	}
	result := make([]Order, 0, len(ids))
	for _, id := range ids {
		result = append(result, *byID[id])
	}
	return result, rows.Err()
}
func (r *PostgresRepository) UpdateStatus(ctx context.Context, id, status string) (Order, error) {
	_, err := r.pool.Exec(ctx, `UPDATE orders SET status=$2,updated_at=NOW() WHERE id=$1`, id, status)
	if err != nil {
		return Order{}, err
	}
	orders, err := r.List(ctx)
	if err != nil {
		return Order{}, err
	}
	for _, o := range orders {
		if o.ID == id {
			return o, nil
		}
	}
	return Order{}, ErrNotFound
}
