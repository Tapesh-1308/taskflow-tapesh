package task

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, t *Task) error
	GetByProject(ctx context.Context, projectID string, status *string, assignee *string) ([]TaskWithUser, error)
	Update(ctx context.Context, id string, body UpdateTaskInput) error
	Delete(ctx context.Context, id string, userID string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, t *Task) error {
	query := `
	INSERT INTO tasks (title, status, project_id, assignee_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query,
		t.Title, t.Status, t.ProjectID, t.AssigneeID,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *repository) GetByProject(
	ctx context.Context,
	projectID string,
	status *string,
	assignee *string,
) ([]TaskWithUser, error) {

	query := `
	SELECT 
		t.id, t.title, t.description, t.status, t.priority,
		t.project_id, t.due_date, t.created_at, t.updated_at,
		u.id, u.name
	FROM tasks t
	LEFT JOIN users u ON t.assignee_id = u.id
	WHERE t.project_id = $1
	`

	args := []interface{}{projectID}
	argPos := 2

	if status != nil {
		query += fmt.Sprintf(" AND t.status = $%d", argPos)
		args = append(args, *status)
		argPos++
	}

	if assignee != nil {
		query += fmt.Sprintf(" AND t.assignee_id = $%d", argPos)
		args = append(args, *assignee)
		argPos++
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []TaskWithUser

	for rows.Next() {
		var t TaskWithUser

		// nullable user fields
		var userID *string
		var userName *string

		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Priority,
			&t.ProjectID,
			&t.DueDate,
			&t.CreatedAt,
			&t.UpdatedAt,
			&userID,
			&userName,
		)
		if err != nil {
			return nil, err
		}

		// assign only if user exists
		if userID != nil {
			t.Assignee = &UserInfo{
				ID:   *userID,
				Name: safeString(userName),
			}
		} else {
			t.Assignee = nil
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *repository) Update(ctx context.Context, id string, body UpdateTaskInput) error {

	query := "UPDATE tasks SET "
	args := []interface{}{}
	i := 1

	if body.Title != nil {
		query += fmt.Sprintf("title = $%d,", i)
		args = append(args, *body.Title)
		i++
	}
	if body.Description != nil {
		query += fmt.Sprintf("description = $%d,", i)
		args = append(args, *body.Description)
		i++
	}
	if body.Status != nil {
		query += fmt.Sprintf("status = $%d,", i)
		args = append(args, *body.Status)
		i++
	}
	if body.Priority != nil {
		query += fmt.Sprintf("priority = $%d,", i)
		args = append(args, *body.Priority)
		i++
	}
	if body.AssigneeID != nil {
		query += fmt.Sprintf("assignee_id = $%d,", i)
		args = append(args, *body.AssigneeID)
		i++
	}
	if body.DueDate != nil {
		query += fmt.Sprintf("due_date = $%d,", i)
		args = append(args, *body.DueDate)
		i++
	}

	query = strings.TrimSuffix(query, ",")
	query += fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, id)

	fmt.Println(query, args)
	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *repository) Delete(ctx context.Context, id string, userID string) error {
	query := `
	DELETE FROM tasks
	WHERE id = $1
	AND (
		assignee_id = $2
		OR project_id IN (
			SELECT id FROM projects WHERE owner_id = $2
		)
	)
	`

	cmd, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("not allowed or not found")
	}

	return nil
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
