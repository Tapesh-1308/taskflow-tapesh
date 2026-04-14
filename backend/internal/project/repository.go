package project

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, p *Project) error
	GetByID(ctx context.Context, id string) (*ProjectWithUser, error)
	GetUserProjects(ctx context.Context, userID string) ([]ProjectWithUser, error)
	Update(ctx context.Context, p *ProjectWithUser) error
	Delete(ctx context.Context, id string) error
	GetProjectWithTasks(ctx context.Context, id string) (*ProjectWithTasks, error)
}

type repository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewRepository(db *pgxpool.Pool, log *slog.Logger) Repository {
	return &repository{db: db, log: log}
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

func (r *repository) GetByID(ctx context.Context, id string) (*ProjectWithUser, error) {
	query := `
		SELECT 
			p.id, p.name, p.description, p.created_at,
			u.id, u.name
		FROM projects p
		LEFT JOIN users u ON u.id = p.owner_id
		WHERE p.id = $1
	`

	var p ProjectWithUser

	var ownerID, ownerName *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.CreatedAt,
		&ownerID,
		&ownerName,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if ownerID != nil && ownerName != nil {
		p.Owner = UserInfo{
			ID:   *ownerID,
			Name: *ownerName,
		}
	}

	return &p, nil
}

func (r *repository) GetUserProjects(ctx context.Context, userID string) ([]ProjectWithUser, error) {
	query := `
		SELECT DISTINCT 
			p.id, p.name, p.description, p.created_at,
			u.id, u.name
		FROM projects p
		LEFT JOIN tasks t ON t.project_id = p.id
		LEFT JOIN users u ON u.id = p.owner_id
		WHERE p.owner_id = $1 OR t.assignee_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []ProjectWithUser

	for rows.Next() {
		var p ProjectWithUser
		var ownerID, ownerName *string

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.CreatedAt,
			&ownerID,
			&ownerName,
		)
		if err != nil {
			return nil, err
		}

		if ownerID != nil && ownerName != nil {
			p.Owner = UserInfo{
				ID:   *ownerID,
				Name: *ownerName,
			}
		}

		projects = append(projects, p)
	}

	return projects, nil
}

func (r *repository) Update(ctx context.Context, p *ProjectWithUser) error {
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
            u.name as owner_name,

            t.id, t.title, t.description, t.status, t.priority,
            t.project_id, t.assignee_id, t.due_date,
            t.created_at, t.updated_at,
            a.name as assignee_name

        FROM projects p
        LEFT JOIN users u ON u.id = p.owner_id
        LEFT JOIN tasks t ON t.project_id = p.id
        LEFT JOIN users a ON a.id = t.assignee_id
        WHERE p.id = $1
    `

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result ProjectWithTasks
	result.Tasks = []Task{}

	taskMap := make(map[string]bool)
	projectSet := false

	for rows.Next() {

		var (
			projectID   string
			projectName string
			projectDesc *string
			ownerID     string
			ownerName   string
			createdAt   time.Time

			taskID        *string
			taskTitle     *string
			taskStatus    *string
			taskPriority  *string
			taskDesc      *string
			projectID2    *string
			assigneeID    *string
			assigneeName  *string
			dueDate       *time.Time
			taskCreatedAt *time.Time
			taskUpdatedAt *time.Time
		)

		err := rows.Scan(
			&projectID,
			&projectName,
			&projectDesc,
			&ownerID,
			&createdAt,
			&ownerName,

			&taskID,
			&taskTitle,
			&taskDesc,
			&taskStatus,
			&taskPriority,
			&projectID2,
			&assigneeID,
			&dueDate,
			&taskCreatedAt,
			&taskUpdatedAt,
			&assigneeName,
		)
		if err != nil {
			return nil, err
		}

		if !projectSet {
			result.ID = projectID
			result.Name = projectName
			result.Description = projectDesc
			result.Owner = UserInfo{ID: ownerID, Name: ownerName}
			result.CreatedAt = createdAt
			projectSet = true
		}

		if taskID != nil && *taskID != "" && !taskMap[*taskID] {

			var assignee *UserInfo
			if assigneeID != nil && assigneeName != nil {
				assignee = &UserInfo{
					ID:   *assigneeID,
					Name: *assigneeName,
				}
			}

			result.Tasks = append(result.Tasks, Task{
				ID:          *taskID,
				Title:       deref(taskTitle),
				Description: taskDesc,
				Status:      deref(taskStatus),
				Priority:    deref(taskPriority),
				ProjectID:   deref(projectID2),
				Assignee:    assignee,
				DueDate:     dueDate,
				CreatedAt:   derefTime(taskCreatedAt),
				UpdatedAt:   derefTime(taskUpdatedAt),
			})

			taskMap[*taskID] = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if result.ID == "" {
		return nil, pgx.ErrNoRows
	}

	return &result, nil
}

func deref(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func derefTime(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}
