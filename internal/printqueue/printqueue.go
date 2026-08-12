package printqueue

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/restaurantautomation/api/internal/printers"
)

type Status string

const (
	StatusQueued   Status = "queued"
	StatusPrinting Status = "printing"
	StatusPrinted  Status = "printed"
	StatusFailed   Status = "failed"
)

type Job struct {
	ID          string          `json:"id"`
	OrderID     string          `json:"order_id"`
	Destination string          `json:"destination"`
	Lines       []printers.Line `json:"lines"`
	Reprint     bool            `json:"reprint"`
	Status      Status          `json:"status"`
	Attempts    int             `json:"attempts"`
	LastError   string          `json:"last_error,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

var ErrNotFound = errors.New("print job not found")
var ErrInvalidJob = errors.New("invalid print job")
var ErrJobNotFailed = errors.New("print job is not failed")
var ErrPrinterManagerUnavailable = errors.New("printer manager unavailable")

type Repository interface {
	Create(context.Context, Job) error
	Get(context.Context, string) (Job, error)
	List(context.Context) ([]Job, error)
	LatestForOrder(context.Context, string) (Job, error)
	ClaimQueued(context.Context, string) (bool, error)
	Update(context.Context, string, Status, int, string) error
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (service *Service) Queue(ctx context.Context, ticket printers.Ticket) (Job, error) {
	if strings.TrimSpace(ticket.OrderID) == "" || strings.TrimSpace(ticket.Destination) == "" || len(ticket.Lines) == 0 {
		return Job{}, ErrInvalidJob
	}
	now := time.Now().UTC()
	job := Job{ID: uuid.NewString(), OrderID: strings.TrimSpace(ticket.OrderID), Destination: strings.TrimSpace(ticket.Destination), Lines: ticket.Lines, Reprint: ticket.Reprint, Status: StatusQueued, CreatedAt: now, UpdatedAt: now}
	return job, service.repository.Create(ctx, job)
}

func (service *Service) Reprint(ctx context.Context, orderID string) (Job, error) {
	job, err := service.repository.LatestForOrder(ctx, strings.TrimSpace(orderID))
	if err != nil {
		return Job{}, err
	}
	return service.Queue(ctx, printers.Ticket{OrderID: job.OrderID, Destination: job.Destination, Lines: job.Lines, Reprint: true})
}

func (service *Service) RetryFailed(ctx context.Context, id string) (Job, error) {
	job, err := service.repository.Get(ctx, strings.TrimSpace(id))
	if err != nil {
		return Job{}, err
	}
	if job.Status != StatusFailed {
		return Job{}, ErrJobNotFailed
	}
	return service.Queue(ctx, printers.Ticket{OrderID: job.OrderID, Destination: job.Destination, Lines: job.Lines, Reprint: true})
}

func (service *Service) List(ctx context.Context) ([]Job, error) { return service.repository.List(ctx) }

// RecoverQueued dispatches tickets that were durably queued before a process
// restart. Jobs already marked printing are intentionally not replayed because
// the previous process may have already sent them to the physical printer.
func (service *Service) RecoverQueued(ctx context.Context, manager *printers.Manager) (int, error) {
	if manager == nil {
		return 0, ErrPrinterManagerUnavailable
	}
	jobs, err := service.List(ctx)
	if err != nil {
		return 0, err
	}
	recovered := 0
	for _, job := range jobs {
		if job.Status != StatusQueued {
			continue
		}
		claimed, err := service.repository.ClaimQueued(ctx, job.ID)
		if err != nil {
			return recovered, err
		}
		if !claimed {
			continue
		}
		ticket := printers.Ticket{OrderID: job.OrderID, Destination: job.Destination, Lines: job.Lines, Reprint: job.Reprint, PrintJobID: job.ID}
		if err := manager.Print(ctx, ticket); err != nil {
			service.Failed(ctx, ticket, job.Attempts+1, err)
			continue
		}
		recovered++
	}
	return recovered, nil
}

func (service *Service) Printing(ctx context.Context, ticket printers.Ticket, attempt int) {
	service.update(ctx, ticket.PrintJobID, StatusPrinting, attempt, "")
}
func (service *Service) Printed(ctx context.Context, ticket printers.Ticket, attempt int) {
	service.update(ctx, ticket.PrintJobID, StatusPrinted, attempt, "")
}
func (service *Service) Failed(ctx context.Context, ticket printers.Ticket, attempt int, err error) {
	message := "print failed"
	if err != nil {
		message = err.Error()
	}
	service.update(ctx, ticket.PrintJobID, StatusFailed, attempt, message)
}
func (service *Service) update(ctx context.Context, id string, status Status, attempt int, message string) {
	if id != "" {
		_ = service.repository.Update(ctx, id, status, attempt, message)
	}
}

type MemoryRepository struct {
	mu   sync.RWMutex
	jobs map[string]Job
}

func NewMemoryRepository() *MemoryRepository { return &MemoryRepository{jobs: make(map[string]Job)} }
func (repository *MemoryRepository) Create(_ context.Context, job Job) error {
	repository.mu.Lock()
	repository.jobs[job.ID] = job
	repository.mu.Unlock()
	return nil
}
func (repository *MemoryRepository) List(_ context.Context) ([]Job, error) {
	repository.mu.RLock()
	result := make([]Job, 0, len(repository.jobs))
	for _, job := range repository.jobs {
		result = append(result, job)
	}
	repository.mu.RUnlock()
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.After(result[j].CreatedAt) })
	return result, nil
}
func (repository *MemoryRepository) Get(_ context.Context, id string) (Job, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	job, ok := repository.jobs[id]
	if !ok {
		return Job{}, ErrNotFound
	}
	return job, nil
}
func (repository *MemoryRepository) LatestForOrder(_ context.Context, orderID string) (Job, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	var latest Job
	for _, job := range repository.jobs {
		if job.OrderID == orderID && (latest.ID == "" || job.CreatedAt.After(latest.CreatedAt)) {
			latest = job
		}
	}
	if latest.ID == "" {
		return Job{}, ErrNotFound
	}
	return latest, nil
}
func (repository *MemoryRepository) Update(_ context.Context, id string, status Status, attempts int, lastError string) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	job, ok := repository.jobs[id]
	if !ok {
		return ErrNotFound
	}
	job.Status, job.Attempts, job.LastError, job.UpdatedAt = status, attempts, lastError, time.Now().UTC()
	repository.jobs[id] = job
	return nil
}
func (repository *MemoryRepository) ClaimQueued(_ context.Context, id string) (bool, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	job, ok := repository.jobs[id]
	if !ok {
		return false, ErrNotFound
	}
	if job.Status != StatusQueued {
		return false, nil
	}
	job.Status, job.Attempts, job.LastError, job.UpdatedAt = StatusPrinting, job.Attempts+1, "", time.Now().UTC()
	repository.jobs[id] = job
	return true, nil
}

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}
func (repository *PostgresRepository) Create(ctx context.Context, job Job) error {
	lines, err := json.Marshal(job.Lines)
	if err != nil {
		return err
	}
	_, err = repository.pool.Exec(ctx, `INSERT INTO print_jobs (id, order_id, destination, lines, reprint, status, attempts, last_error, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, job.ID, job.OrderID, job.Destination, lines, job.Reprint, job.Status, job.Attempts, job.LastError, job.CreatedAt, job.UpdatedAt)
	return err
}
func (repository *PostgresRepository) List(ctx context.Context) ([]Job, error) {
	rows, err := repository.pool.Query(ctx, `SELECT id, order_id, destination, lines, reprint, status, attempts, last_error, created_at, updated_at FROM print_jobs ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJobs(rows)
}
func (repository *PostgresRepository) Get(ctx context.Context, id string) (Job, error) {
	rows, err := repository.pool.Query(ctx, `SELECT id, order_id, destination, lines, reprint, status, attempts, last_error, created_at, updated_at FROM print_jobs WHERE id=$1`, id)
	if err != nil {
		return Job{}, err
	}
	defer rows.Close()
	jobs, err := scanJobs(rows)
	if err != nil {
		return Job{}, err
	}
	if len(jobs) == 0 {
		return Job{}, ErrNotFound
	}
	return jobs[0], nil
}
func (repository *PostgresRepository) LatestForOrder(ctx context.Context, orderID string) (Job, error) {
	rows, err := repository.pool.Query(ctx, `SELECT id, order_id, destination, lines, reprint, status, attempts, last_error, created_at, updated_at FROM print_jobs WHERE order_id=$1 ORDER BY created_at DESC LIMIT 1`, orderID)
	if err != nil {
		return Job{}, err
	}
	defer rows.Close()
	jobs, err := scanJobs(rows)
	if err != nil {
		return Job{}, err
	}
	if len(jobs) == 0 {
		return Job{}, ErrNotFound
	}
	return jobs[0], nil
}
func (repository *PostgresRepository) Update(ctx context.Context, id string, status Status, attempts int, lastError string) error {
	command, err := repository.pool.Exec(ctx, `UPDATE print_jobs SET status=$2, attempts=$3, last_error=$4, updated_at=NOW() WHERE id=$1`, id, status, attempts, lastError)
	if err != nil {
		return err
	}
	if command.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
func (repository *PostgresRepository) ClaimQueued(ctx context.Context, id string) (bool, error) {
	command, err := repository.pool.Exec(ctx, `UPDATE print_jobs SET status=$2, attempts=attempts+1, last_error='', updated_at=NOW() WHERE id=$1 AND status=$3`, id, StatusPrinting, StatusQueued)
	if err != nil {
		return false, err
	}
	return command.RowsAffected() > 0, nil
}
func scanJobs(rows pgx.Rows) ([]Job, error) {
	result := make([]Job, 0)
	for rows.Next() {
		var job Job
		var lines []byte
		if err := rows.Scan(&job.ID, &job.OrderID, &job.Destination, &lines, &job.Reprint, &job.Status, &job.Attempts, &job.LastError, &job.CreatedAt, &job.UpdatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(lines, &job.Lines); err != nil {
			return nil, err
		}
		result = append(result, job)
	}
	return result, rows.Err()
}
