package task

import "time"

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	ProjectID   string     `json:"project_id"`
	AssigneeID  *string    `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskWithUser struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	ProjectID   string     `json:"project_id"`
	Assignee    *UserInfo  `json:"assignee"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type UpdateTaskInput struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	Priority    *string    `json:"priority"`
	AssigneeID  *string    `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
}

type UserInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
