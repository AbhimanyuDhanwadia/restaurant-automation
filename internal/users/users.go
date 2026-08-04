package users

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/restaurantautomation/api/internal/roles"
)

var ErrNotFound = errors.New("user not found")
var ErrInvalidUser = errors.New("invalid user")
var ErrInvalidRole = errors.New("invalid role")

type User struct {
	AuthSubject string    `json:"auth_subject"`
	Email       string    `json:"email"`
	RoleID      string    `json:"role_id"`
	RoleName    string    `json:"role_name"`
	CreatedAt   time.Time `json:"created_at"`
	LastSeenAt  time.Time `json:"last_seen_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Repository interface {
	Observe(context.Context, User) error
	List(context.Context) ([]User, error)
	UpdateRole(context.Context, string, string, time.Time) (User, error)
}

type Service struct {
	repository  Repository
	roleService *roles.Service
	adminEmails map[string]struct{}
}

func NewService(repository Repository, roleService *roles.Service, adminEmails []string) *Service {
	allowed := make(map[string]struct{}, len(adminEmails))
	for _, email := range adminEmails {
		if normalized := strings.ToLower(strings.TrimSpace(email)); normalized != "" {
			allowed[normalized] = struct{}{}
		}
	}
	return &Service{repository: repository, roleService: roleService, adminEmails: allowed}
}

func (s *Service) Observe(ctx context.Context, subject, email string) error {
	subject = strings.TrimSpace(subject)
	email = strings.ToLower(strings.TrimSpace(email))
	if subject == "" || email == "" || len(subject) > 255 || len(email) > 320 {
		return ErrInvalidUser
	}
	roleID := roles.OperatorID
	if _, ok := s.adminEmails[email]; ok {
		roleID = roles.AdministratorID
	}
	now := time.Now().UTC()
	return s.repository.Observe(ctx, User{AuthSubject: subject, Email: email, RoleID: roleID, CreatedAt: now, LastSeenAt: now, UpdatedAt: now})
}

func (s *Service) List(ctx context.Context) ([]User, error) {
	result, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}
	return s.decorateRoles(ctx, result)
}

func (s *Service) UpdateRole(ctx context.Context, subject, roleID string) (User, error) {
	subject = strings.TrimSpace(subject)
	roleID = strings.TrimSpace(roleID)
	if subject == "" || roleID == "" {
		return User{}, ErrInvalidUser
	}
	roleCatalog, err := s.rolesByID(ctx)
	if err != nil {
		return User{}, err
	}
	if _, ok := roleCatalog[roleID]; !ok {
		return User{}, ErrInvalidRole
	}
	user, err := s.repository.UpdateRole(ctx, subject, roleID, time.Now().UTC())
	if err != nil {
		return User{}, err
	}
	user.RoleName = roleCatalog[roleID].Name
	return user, nil
}

func (s *Service) decorateRoles(ctx context.Context, users []User) ([]User, error) {
	roleCatalog, err := s.rolesByID(ctx)
	if err != nil {
		return nil, err
	}
	for index := range users {
		users[index].RoleName = roleCatalog[users[index].RoleID].Name
	}
	return users, nil
}

func (s *Service) rolesByID(ctx context.Context) (map[string]roles.Role, error) {
	if s.roleService == nil {
		return nil, ErrInvalidRole
	}
	catalog, err := s.roleService.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string]roles.Role, len(catalog))
	for _, role := range catalog {
		result[role.ID] = role
	}
	return result, nil
}

type MemoryRepository struct {
	mu    sync.RWMutex
	users map[string]User
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{users: make(map[string]User)}
}

func (r *MemoryRepository) Observe(_ context.Context, user User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if existing, ok := r.users[user.AuthSubject]; ok {
		existing.Email = user.Email
		existing.LastSeenAt = user.LastSeenAt
		existing.UpdatedAt = user.UpdatedAt
		r.users[user.AuthSubject] = existing
		return nil
	}
	r.users[user.AuthSubject] = user
	return nil
}

func (r *MemoryRepository) List(_ context.Context) ([]User, error) {
	r.mu.RLock()
	result := make([]User, 0, len(r.users))
	for _, user := range r.users {
		result = append(result, user)
	}
	r.mu.RUnlock()
	sort.Slice(result, func(i, j int) bool { return result[i].Email < result[j].Email })
	return result, nil
}

func (r *MemoryRepository) UpdateRole(_ context.Context, subject, roleID string, updatedAt time.Time) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.users[subject]
	if !ok {
		return User{}, ErrNotFound
	}
	user.RoleID = roleID
	user.UpdatedAt = updatedAt
	r.users[subject] = user
	return user, nil
}

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Observe(ctx context.Context, user User) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO app_users (auth_subject,email,role_id,created_at,last_seen_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (auth_subject) DO UPDATE SET email=EXCLUDED.email,last_seen_at=EXCLUDED.last_seen_at,updated_at=EXCLUDED.updated_at`, user.AuthSubject, user.Email, user.RoleID, user.CreatedAt, user.LastSeenAt, user.UpdatedAt)
	return err
}

func (r *PostgresRepository) List(ctx context.Context) ([]User, error) {
	rows, err := r.pool.Query(ctx, `SELECT auth_subject,email,role_id,created_at,last_seen_at,updated_at FROM app_users ORDER BY LOWER(email),auth_subject`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.AuthSubject, &user.Email, &user.RoleID, &user.CreatedAt, &user.LastSeenAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, user)
	}
	return result, rows.Err()
}

func (r *PostgresRepository) UpdateRole(ctx context.Context, subject, roleID string, updatedAt time.Time) (User, error) {
	var user User
	err := r.pool.QueryRow(ctx, `UPDATE app_users SET role_id=$2,updated_at=$3 WHERE auth_subject=$1 RETURNING auth_subject,email,role_id,created_at,last_seen_at,updated_at`, subject, roleID, updatedAt).Scan(&user.AuthSubject, &user.Email, &user.RoleID, &user.CreatedAt, &user.LastSeenAt, &user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return user, err
}

var _ Repository = (*MemoryRepository)(nil)
var _ Repository = (*PostgresRepository)(nil)
