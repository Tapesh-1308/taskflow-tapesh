package user

import (
	"context"

	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAllUsers(ctx context.Context, search string) ([]User, error)
}

type repository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewRepository(db *pgxpool.Pool, log *slog.Logger) Repository {
	return &repository{db: db, log: log}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	r.log.Info("Creating user", "email", user.Email, "name", user.Name)

	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		user.Name,
		user.Email,
		user.Password,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		r.log.Error("Failed to create user", "error", err, "email", user.Email)
		return err
	}

	r.log.Info("User created successfully", "id", user.ID, "email", user.Email)
	return nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	r.log.Debug("Getting user by email", "email", email)

	query := `
		SELECT id, name, email, password, created_at
		FROM users
		WHERE email = $1
	`

	var user User
	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		r.log.Warn("User not found", "email", email, "error", err)
		return nil, err
	}

	r.log.Debug("User found", "id", user.ID, "email", user.Email)
	return &user, nil
}

func (r *repository) GetAllUsers(ctx context.Context, search string) ([]User, error) {
	r.log.Info("Getting all users", "search", search)

	query := `
		SELECT id, name
		FROM users
	`

	var users []User

	if search != "" {
		query += ` WHERE name ILIKE $1`
		rows, err := r.db.Query(ctx, query, "%"+search+"%")
		if err != nil {
			r.log.Error("Failed to query users with search", "error", err, "search", search)
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var user User
			if err := rows.Scan(&user.ID, &user.Name); err != nil {
				r.log.Error("Failed to scan user", "error", err)
				return nil, err
			}
			users = append(users, user)
		}
		if err := rows.Err(); err != nil {
			r.log.Error("Rows error", "error", err)
			return nil, err
		}
		r.log.Info("Users retrieved with search", "count", len(users), "search", search)
		return users, nil
	}

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.log.Error("Failed to query all users", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			r.log.Error("Failed to scan user", "error", err)
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		r.log.Error("Rows error", "error", err)
		return nil, err
	}

	r.log.Info("All users retrieved", "count", len(users))
	return users, nil
}
