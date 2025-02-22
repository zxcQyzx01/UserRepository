package postgres

import (
	"UserRepository/internal/domain"
	"context"
	"database/sql"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	query := `
        INSERT INTO users (id, name, email, created_at)
        VALUES ($1, $2, $3, $4)
    `
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, time.Now())
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	var user domain.User
	query := `
        SELECT id, name, email, created_at, deleted_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.DeletedAt,
	)
	return user, err
}

func (r *UserRepository) Update(ctx context.Context, user domain.User) error {
	query := `
        UPDATE users
        SET name = $1, email = $2
        WHERE id = $3 AND deleted_at IS NULL
    `
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.ID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `
        UPDATE users
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at IS NULL
    `
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *UserRepository) List(ctx context.Context, c domain.Conditions) ([]domain.User, error) {
	query := `
        SELECT id, name, email, created_at, deleted_at
        FROM users
        WHERE deleted_at IS NULL
        ORDER BY created_at
        LIMIT $1 OFFSET $2
    `
	rows, err := r.db.QueryContext(ctx, query, c.Limit, c.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.DeletedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
