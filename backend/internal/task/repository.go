package task

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, t *Task) error
	UpdateStatus(ctx context.Context, id string, status string) error
	Assign(ctx context.Context, id string, userID string) error
	GetByProject(ctx context.Context, projectID string, status *string) ([]Task, error)
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

func (r *repository) UpdateStatus(ctx context.Context, id string, status string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE tasks SET status = $1 WHERE id = $2`,
		status, id,
	)
	return err
}

func (r *repository) Assign(ctx context.Context, id string, userID string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE tasks SET assignee_id = $1 WHERE id = $2`,
		userID, id,
	)
	return err
}

func (r *repository) GetByProject(ctx context.Context, projectID string, status *string) ([]Task, error) {
	query := `
	SELECT id, title, status, project_id, assignee_id, created_at
	FROM tasks
	WHERE project_id = $1
	`

	args := []interface{}{projectID}

	if status != nil {
		query += " AND status = $2"
		args = append(args, *status)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Status, &t.ProjectID, &t.AssigneeID, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}
