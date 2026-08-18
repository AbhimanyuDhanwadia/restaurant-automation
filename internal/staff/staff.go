package staff

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

type Member struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Station   string    `json:"station"`
	Status    string    `json:"status"`
	Handoff   string    `json:"handoff"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Owner     string    `json:"owner"`
	DueLabel  string    `json:"due_label"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Handoff struct {
	Note      string    `json:"note"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaskSummary struct {
	Completed int `json:"completed"`
	Total     int `json:"total"`
}

type CreateMemberInput struct {
	Name    string `json:"name"`
	Role    string `json:"role"`
	Station string `json:"station"`
	Status  string `json:"status"`
	Handoff string `json:"handoff"`
}
type CreateTaskInput struct {
	Title    string `json:"title"`
	Owner    string `json:"owner"`
	DueLabel string `json:"due_label"`
}

var ErrNotFound = errors.New("staff resource not found")
var ErrInvalidMember = errors.New("invalid staff member")
var ErrInvalidTask = errors.New("invalid shift task")
var ErrInvalidStatus = errors.New("invalid staff status")

type Repository interface {
	CreateMember(context.Context, Member) error
	ListMembers(context.Context) ([]Member, error)
	UpdateMemberStatus(context.Context, string, string) (Member, error)
	UpdateMemberHandoff(context.Context, string, string) (Member, error)
	CreateTask(context.Context, Task) error
	ListTasks(context.Context) ([]Task, error)
	TaskSummary(context.Context) (TaskSummary, error)
	CompleteTask(context.Context, string) (Task, error)
	GetHandoff(context.Context) (Handoff, error)
	SaveHandoff(context.Context, Handoff) error
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }
func (s *Service) CreateMember(ctx context.Context, input CreateMemberInput) (Member, error) {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Role) == "" || strings.TrimSpace(input.Station) == "" {
		return Member{}, ErrInvalidMember
	}
	status := input.Status
	if status == "" {
		status = "off_shift"
	}
	if !validStatus(status) {
		return Member{}, ErrInvalidStatus
	}
	now := time.Now().UTC()
	member := Member{ID: uuid.NewString(), Name: strings.TrimSpace(input.Name), Role: strings.TrimSpace(input.Role), Station: strings.TrimSpace(input.Station), Status: status, Handoff: strings.TrimSpace(input.Handoff), CreatedAt: now, UpdatedAt: now}
	return member, s.repository.CreateMember(ctx, member)
}
func (s *Service) ListMembers(ctx context.Context) ([]Member, error) {
	return s.repository.ListMembers(ctx)
}
func (s *Service) UpdateMemberStatus(ctx context.Context, id, status string) (Member, error) {
	if !validStatus(status) {
		return Member{}, ErrInvalidStatus
	}
	return s.repository.UpdateMemberStatus(ctx, id, status)
}
func (s *Service) UpdateMemberHandoff(ctx context.Context, id, handoff string) (Member, error) {
	return s.repository.UpdateMemberHandoff(ctx, id, strings.TrimSpace(handoff))
}
func (s *Service) CreateTask(ctx context.Context, input CreateTaskInput) (Task, error) {
	if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Owner) == "" {
		return Task{}, ErrInvalidTask
	}
	now := time.Now().UTC()
	task := Task{ID: uuid.NewString(), Title: strings.TrimSpace(input.Title), Owner: strings.TrimSpace(input.Owner), DueLabel: strings.TrimSpace(input.DueLabel), CreatedAt: now, UpdatedAt: now}
	return task, s.repository.CreateTask(ctx, task)
}
func (s *Service) ListTasks(ctx context.Context) ([]Task, error) { return s.repository.ListTasks(ctx) }
func (s *Service) TaskSummary(ctx context.Context) (TaskSummary, error) {
	return s.repository.TaskSummary(ctx)
}
func (s *Service) CompleteTask(ctx context.Context, id string) (Task, error) {
	return s.repository.CompleteTask(ctx, id)
}
func (s *Service) GetHandoff(ctx context.Context) (Handoff, error) {
	return s.repository.GetHandoff(ctx)
}
func (s *Service) SaveHandoff(ctx context.Context, note string) (Handoff, error) {
	handoff := Handoff{Note: strings.TrimSpace(note), UpdatedAt: time.Now().UTC()}
	return handoff, s.repository.SaveHandoff(ctx, handoff)
}
func validStatus(status string) bool {
	return status == "on_shift" || status == "on_break" || status == "off_shift"
}

type MemoryRepository struct {
	mu      sync.RWMutex
	members map[string]Member
	tasks   map[string]Task
	handoff Handoff
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{members: map[string]Member{}, tasks: map[string]Task{}}
}
func (r *MemoryRepository) CreateMember(_ context.Context, member Member) error {
	r.mu.Lock()
	r.members[member.ID] = member
	r.mu.Unlock()
	return nil
}
func (r *MemoryRepository) ListMembers(_ context.Context) ([]Member, error) {
	r.mu.RLock()
	result := make([]Member, 0, len(r.members))
	for _, member := range r.members {
		result = append(result, member)
	}
	r.mu.RUnlock()
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result, nil
}
func (r *MemoryRepository) UpdateMemberStatus(_ context.Context, id, status string) (Member, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	member, ok := r.members[id]
	if !ok {
		return Member{}, ErrNotFound
	}
	member.Status = status
	member.UpdatedAt = time.Now().UTC()
	r.members[id] = member
	return member, nil
}
func (r *MemoryRepository) UpdateMemberHandoff(_ context.Context, id, handoff string) (Member, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	member, ok := r.members[id]
	if !ok {
		return Member{}, ErrNotFound
	}
	member.Handoff = handoff
	member.UpdatedAt = time.Now().UTC()
	r.members[id] = member
	return member, nil
}
func (r *MemoryRepository) CreateTask(_ context.Context, task Task) error {
	r.mu.Lock()
	r.tasks[task.ID] = task
	r.mu.Unlock()
	return nil
}
func (r *MemoryRepository) ListTasks(_ context.Context) ([]Task, error) {
	r.mu.RLock()
	result := make([]Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		if !task.Completed {
			result = append(result, task)
		}
	}
	r.mu.RUnlock()
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.Before(result[j].CreatedAt) })
	return result, nil
}
func (r *MemoryRepository) TaskSummary(_ context.Context) (TaskSummary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	summary := TaskSummary{Total: len(r.tasks)}
	for _, task := range r.tasks {
		if task.Completed {
			summary.Completed++
		}
	}
	return summary, nil
}
func (r *MemoryRepository) CompleteTask(_ context.Context, id string) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	task, ok := r.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	task.Completed = true
	task.UpdatedAt = time.Now().UTC()
	r.tasks[id] = task
	return task, nil
}
func (r *MemoryRepository) GetHandoff(_ context.Context) (Handoff, error) {
	r.mu.RLock()
	handoff := r.handoff
	r.mu.RUnlock()
	return handoff, nil
}
func (r *MemoryRepository) SaveHandoff(_ context.Context, handoff Handoff) error {
	r.mu.Lock()
	r.handoff = handoff
	r.mu.Unlock()
	return nil
}

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}
func (r *PostgresRepository) CreateMember(ctx context.Context, member Member) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO staff_members (id, name, role, station, status, handoff, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, member.ID, member.Name, member.Role, member.Station, member.Status, member.Handoff, member.CreatedAt, member.UpdatedAt)
	return err
}
func (r *PostgresRepository) ListMembers(ctx context.Context) ([]Member, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,name,role,station,status,handoff,created_at,updated_at FROM staff_members ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []Member{}
	for rows.Next() {
		var member Member
		if err := rows.Scan(&member.ID, &member.Name, &member.Role, &member.Station, &member.Status, &member.Handoff, &member.CreatedAt, &member.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, member)
	}
	return result, rows.Err()
}
func (r *PostgresRepository) UpdateMemberStatus(ctx context.Context, id, status string) (Member, error) {
	return r.member(ctx, `UPDATE staff_members SET status=$2,updated_at=NOW() WHERE id=$1 RETURNING id,name,role,station,status,handoff,created_at,updated_at`, id, status)
}
func (r *PostgresRepository) UpdateMemberHandoff(ctx context.Context, id, handoff string) (Member, error) {
	return r.member(ctx, `UPDATE staff_members SET handoff=$2,updated_at=NOW() WHERE id=$1 RETURNING id,name,role,station,status,handoff,created_at,updated_at`, id, handoff)
}
func (r *PostgresRepository) member(ctx context.Context, query string, args ...any) (Member, error) {
	var member Member
	err := r.pool.QueryRow(ctx, query, args...).Scan(&member.ID, &member.Name, &member.Role, &member.Station, &member.Status, &member.Handoff, &member.CreatedAt, &member.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Member{}, ErrNotFound
	}
	return member, err
}
func (r *PostgresRepository) CreateTask(ctx context.Context, task Task) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO shift_tasks (id,title,owner,due_label,completed,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`, task.ID, task.Title, task.Owner, task.DueLabel, task.Completed, task.CreatedAt, task.UpdatedAt)
	return err
}
func (r *PostgresRepository) ListTasks(ctx context.Context) ([]Task, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,title,owner,due_label,completed,created_at,updated_at FROM shift_tasks WHERE completed=FALSE ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []Task{}
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Owner, &task.DueLabel, &task.Completed, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, task)
	}
	return result, rows.Err()
}
func (r *PostgresRepository) TaskSummary(ctx context.Context) (TaskSummary, error) {
	var summary TaskSummary
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FILTER (WHERE completed), COUNT(*) FROM shift_tasks`).Scan(&summary.Completed, &summary.Total)
	return summary, err
}
func (r *PostgresRepository) CompleteTask(ctx context.Context, id string) (Task, error) {
	var task Task
	err := r.pool.QueryRow(ctx, `UPDATE shift_tasks SET completed=TRUE,updated_at=NOW() WHERE id=$1 RETURNING id,title,owner,due_label,completed,created_at,updated_at`, id).Scan(&task.ID, &task.Title, &task.Owner, &task.DueLabel, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, ErrNotFound
	}
	return task, err
}
func (r *PostgresRepository) GetHandoff(ctx context.Context) (Handoff, error) {
	var handoff Handoff
	err := r.pool.QueryRow(ctx, `SELECT note,updated_at FROM shift_handoffs WHERE id='current'`).Scan(&handoff.Note, &handoff.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Handoff{}, nil
	}
	return handoff, err
}
func (r *PostgresRepository) SaveHandoff(ctx context.Context, handoff Handoff) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO shift_handoffs (id,note,updated_at) VALUES ('current',$1,$2) ON CONFLICT (id) DO UPDATE SET note=EXCLUDED.note,updated_at=EXCLUDED.updated_at`, handoff.Note, handoff.UpdatedAt)
	return err
}
