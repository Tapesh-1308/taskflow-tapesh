package task

import "time"

type Task struct {
	ID         string
	Title      string
	Status     string
	ProjectID  string
	AssigneeID *string
	CreatedAt  time.Time
}
