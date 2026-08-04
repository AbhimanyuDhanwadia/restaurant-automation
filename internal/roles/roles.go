package roles

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var AllowedPermissions = map[string]struct{}{
	"analytics.view": {}, "audit.view": {}, "automation.view": {}, "backups.manage": {}, "database.view": {}, "integrations.manage": {}, "inventory.manage": {}, "kitchen.manage": {}, "operations.view": {}, "orders.manage": {}, "printers.manage": {}, "roles.manage": {}, "settings.manage": {}, "staff.manage": {}, "users.manage": {},
}

const (
	OperatorID      = "00000000-0000-0000-0000-000000000001"
	ManagerID       = "00000000-0000-0000-0000-000000000002"
	AdministratorID = "00000000-0000-0000-0000-000000000003"
)

type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions []string  `json:"permissions"`
	System      bool      `json:"system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateInput struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

var ErrNotFound = errors.New("role not found")
var ErrDuplicate = errors.New("role already exists")
var ErrInvalidRole = errors.New("invalid role")

type Repository interface {
	Create(context.Context, Role) error
	List(context.Context) ([]Role, error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) Create(ctx context.Context, input CreateInput) (Role, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" || len(name) > 80 {
		return Role{}, ErrInvalidRole
	}
	permissions, valid := normalizePermissions(input.Permissions)
	if !valid {
		return Role{}, ErrInvalidRole
	}
	now := time.Now().UTC()
	role := Role{ID: uuid.NewString(), Name: name, Description: strings.TrimSpace(input.Description), Permissions: permissions, CreatedAt: now, UpdatedAt: now}
	return role, s.repository.Create(ctx, role)
}

func (s *Service) List(ctx context.Context) ([]Role, error) { return s.repository.List(ctx) }

func normalizePermissions(values []string) ([]string, bool) {
	unique := make(map[string]struct{}, len(values))
	for _, value := range values {
		if _, ok := AllowedPermissions[value]; !ok {
			return nil, false
		}
		unique[value] = struct{}{}
	}
	if len(unique) == 0 {
		return nil, false
	}
	permissions := make([]string, 0, len(unique))
	for value := range unique {
		permissions = append(permissions, value)
	}
	sort.Strings(permissions)
	return permissions, true
}

func defaultRoles() []Role {
	now := time.Now().UTC()
	return []Role{
		{ID: OperatorID, Name: "Operator", Description: "Runs day-to-day restaurant operations.", Permissions: []string{"inventory.manage", "kitchen.manage", "operations.view", "orders.manage"}, System: true, CreatedAt: now, UpdatedAt: now},
		{ID: ManagerID, Name: "Manager", Description: "Oversees operations, staffing, and analytics.", Permissions: []string{"analytics.view", "inventory.manage", "kitchen.manage", "operations.view", "orders.manage", "settings.manage", "staff.manage"}, System: true, CreatedAt: now, UpdatedAt: now},
		{ID: AdministratorID, Name: "Administrator", Description: "Configures automation and administration services.", Permissions: []string{"analytics.view", "audit.view", "automation.view", "backups.manage", "database.view", "integrations.manage", "inventory.manage", "kitchen.manage", "operations.view", "orders.manage", "printers.manage", "roles.manage", "settings.manage", "staff.manage", "users.manage"}, System: true, CreatedAt: now, UpdatedAt: now},
	}
}

type MemoryRepository struct {
	mu    sync.RWMutex
	roles map[string]Role
}

func NewMemoryRepository() *MemoryRepository {
	repository := &MemoryRepository{roles: make(map[string]Role)}
	for _, role := range defaultRoles() {
		repository.roles[role.ID] = role
	}
	return repository
}

func (r *MemoryRepository) Create(_ context.Context, role Role) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.roles {
		if strings.EqualFold(existing.Name, role.Name) {
			return ErrDuplicate
		}
	}
	r.roles[role.ID] = role
	return nil
}

func (r *MemoryRepository) List(_ context.Context) ([]Role, error) {
	r.mu.RLock()
	result := make([]Role, 0, len(r.roles))
	for _, role := range r.roles {
		result = append(result, role)
	}
	r.mu.RUnlock()
	sort.Slice(result, func(i, j int) bool {
		if result[i].System != result[j].System {
			return result[i].System
		}
		return result[i].Name < result[j].Name
	})
	return result, nil
}

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, role Role) error {
	permissions, err := json.Marshal(role.Permissions)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `INSERT INTO app_roles (id,name,description,permissions,system,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`, role.ID, role.Name, role.Description, permissions, role.System, role.CreatedAt, role.UpdatedAt)
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" {
		return ErrDuplicate
	}
	return err
}

func (r *PostgresRepository) List(ctx context.Context) ([]Role, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,name,description,permissions,system,created_at,updated_at FROM app_roles ORDER BY system DESC,name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []Role{}
	for rows.Next() {
		var role Role
		var permissions []byte
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &permissions, &role.System, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(permissions, &role.Permissions); err != nil {
			return nil, err
		}
		result = append(result, role)
	}
	return result, rows.Err()
}

var _ Repository = (*MemoryRepository)(nil)
var _ Repository = (*PostgresRepository)(nil)
