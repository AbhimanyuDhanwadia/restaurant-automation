package backups

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

type Status string

const (
	StatusRecorded Status = "recorded"
	StatusVerified Status = "verified"
	StatusFailed   Status = "failed"
)

type Record struct {
	ID          string    `json:"id"`
	Target      string    `json:"target"`
	SizeBytes   int64     `json:"size_bytes"`
	Status      Status    `json:"status"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type RecordInput struct {
	Target    string `json:"target"`
	SizeBytes int64  `json:"size_bytes"`
}

var ErrInvalidRecord = errors.New("invalid backup record")
var ErrNotFound = errors.New("backup record not found")

type Repository interface {
	Create(context.Context, Record) error
	List(context.Context) ([]Record, error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) Record(ctx context.Context, input RecordInput) (Record, error) {
	target := strings.TrimSpace(input.Target)
	if target == "" || len(target) > 500 || input.SizeBytes < 0 {
		return Record{}, ErrInvalidRecord
	}
	now := time.Now().UTC()
	record := Record{ID: uuid.NewString(), Target: target, SizeBytes: input.SizeBytes, Status: StatusRecorded, CompletedAt: now, CreatedAt: now}
	return record, s.repository.Create(ctx, record)
}

func (s *Service) List(ctx context.Context) ([]Record, error) { return s.repository.List(ctx) }

type MemoryRepository struct {
	mu      sync.RWMutex
	records map[string]Record
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{records: make(map[string]Record)}
}

func (r *MemoryRepository) Create(_ context.Context, record Record) error {
	r.mu.Lock()
	r.records[record.ID] = record
	r.mu.Unlock()
	return nil
}

func (r *MemoryRepository) List(_ context.Context) ([]Record, error) {
	r.mu.RLock()
	result := make([]Record, 0, len(r.records))
	for _, record := range r.records {
		result = append(result, record)
	}
	r.mu.RUnlock()
	sort.Slice(result, func(i, j int) bool { return result[i].CompletedAt.After(result[j].CompletedAt) })
	return result, nil
}

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, record Record) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO backup_records (id,target,size_bytes,status,completed_at,created_at) VALUES ($1,$2,$3,$4,$5,$6)`, record.ID, record.Target, record.SizeBytes, record.Status, record.CompletedAt, record.CreatedAt)
	return err
}

func (r *PostgresRepository) List(ctx context.Context) ([]Record, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,target,size_bytes,status,completed_at,created_at FROM backup_records ORDER BY completed_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []Record{}
	for rows.Next() {
		var record Record
		if err := rows.Scan(&record.ID, &record.Target, &record.SizeBytes, &record.Status, &record.CompletedAt, &record.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return result, rows.Err()
}

var _ Repository = (*MemoryRepository)(nil)
var _ Repository = (*PostgresRepository)(nil)
