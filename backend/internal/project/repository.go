package project

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, p *Project) error
	GetByID(ctx context.Context, id string) (*Project, error)
	GetUserProjects(ctx context.Context, userID string) ([]Project, error)
	Update(ctx context.Context, p *Project) error
	Delete(ctx context.Context, id string) error
	GetProjectWithTasks(ctx context.Context, id string) (*ProjectWithTasks, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, p *Project) error {
	query := `
		INSERT INTO projects (name, description, owner_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query,
		p.Name,
		p.Description,
		p.OwnerID,
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *repository) GetByID(ctx context.Context, id string) (*Project, error) {
	query := `
		SELECT id, name, description, owner_id, created_at
		FROM projects
		WHERE id = $1
	`

	var p Project
	err := r.db.QueryRow(ctx, query, id).
		Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *repository) GetUserProjects(ctx context.Context, userID string) ([]Project, error) {
	query := `
		SELECT DISTINCT p.id, p.name, p.description, p.owner_id, p.created_at
		FROM projects p
		LEFT JOIN tasks t ON t.project_id = p.id
		WHERE p.owner_id = $1 OR t.assignee_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project

	for rows.Next() {
		var p Project
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (r *repository) Update(ctx context.Context, p *Project) error {
	query := `
		UPDATE projects
		SET name = $1, description = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query, p.Name, p.Description, p.ID)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *repository) GetProjectWithTasks(ctx context.Context, id string) (*ProjectWithTasks, error) {
	query := `
	SELECT 
		p.id, p.name, p.description, p.owner_id, p.created_at,
		t.id, t.title, t.status, t.assignee_id
	FROM projects p
	LEFT JOIN tasks t ON t.project_id = p.id
	WHERE p.id = $1
	`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result ProjectWithTasks
	taskMap := make(map[string]bool)

	for rows.Next() {
		var p Project
		var t Task

		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt,
			&t.ID, &t.Title, &t.Status, &t.AssigneeID,
		)
		if err != nil {
			return nil, err
		}

		result.Project = p

		if t.ID != "" && !taskMap[t.ID] {
			result.Tasks = append(result.Tasks, t)
			taskMap[t.ID] = true
		}
	}

	return &result, nil
}
