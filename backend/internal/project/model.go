package project

import "time"

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	Tasks       []Task    `json:"tasks"`
}

type ProjectWithUser struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Owner       UserInfo  `json:"owner"`
	CreatedAt   time.Time `json:"created_at"`
	Tasks       []Task    `json:"tasks"`
}

type ProjectWithTasks struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Owner       UserInfo  `json:"owner"`
	CreatedAt   time.Time `json:"created_at"`
	Tasks       []Task    `json:"tasks"`
}

type Task struct {
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

type UserInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
