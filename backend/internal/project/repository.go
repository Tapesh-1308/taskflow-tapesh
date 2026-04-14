package project

import (
	"context"
	"fmt"
	"time"

	"log/slog"

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
			p.id, 
			p.name, 
			p.description, 
			p.created_at,
			u.id,
			u.name
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

	var owner UserInfo
	if ownerID != nil && ownerName != nil {
		owner = UserInfo{
			ID:   *ownerID,
			Name: *ownerName,
		}
	}

	p.Owner = owner
	fmt.Println("ERROR::: ", p)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *repository) GetUserProjects(ctx context.Context, userID string) ([]ProjectWithUser, error) {
	query := `
		SELECT DISTINCT 
			p.id, 
			p.name, 
			p.description, 
			p.created_at,
			u.id,
			u.name
		FROM projects p
		LEFT JOIN tasks t ON t.project_id = p.id
		LEFT JOIN users u ON u.id = p.owner_id
		WHERE p.owner_id = $1 OR t.assignee_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	fmt.Println(query, userID)
	if err != nil {
		fmt.Println("ERROR,", err)
		return nil, err
	}
	defer rows.Close()

	var projects []ProjectWithUser

	for rows.Next() {
		var p ProjectWithUser
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.CreatedAt,
			&p.Owner.ID,
			&p.Owner.Name,
		)
		if err != nil {
			return nil, err
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
		var projectID, projectName string
		var projectDesc *string
		var ownerID, ownerName string
		var createdAt time.Time

		var taskID, taskTitle, taskStatus, taskPriority string
		var taskDesc *string
		var projectID2, assigneeID *string
		var assigneeName *string
		var dueDate *time.Time
		var taskCreatedAt, taskUpdatedAt time.Time

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

		// Set project fields only once
		if !projectSet {
			result.ID = projectID
			result.Name = projectName
			result.Description = projectDesc
			result.Owner = UserInfo{ID: ownerID, Name: ownerName}
			result.CreatedAt = createdAt
			projectSet = true
		}

		// Skip NULL task rows (LEFT JOIN case)
		if taskID != "" && !taskMap[taskID] {
			var assignee *UserInfo
			if assigneeID != nil && *assigneeID != "" {
				assignee = &UserInfo{ID: *assigneeID, Name: *assigneeName}
			}

			result.Tasks = append(result.Tasks, Task{
				ID:          taskID,
				Title:       taskTitle,
				Description: taskDesc,
				Status:      taskStatus,
				Priority:    taskPriority,
				ProjectID:   *projectID2,
				Assignee:    assignee,
				DueDate:     dueDate,
				CreatedAt:   taskCreatedAt,
				UpdatedAt:   taskUpdatedAt,
			})
			taskMap[taskID] = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Check if project was found
	if result.ID == "" {
		return nil, pgx.ErrNoRows
	}

	return &result, nil
}
