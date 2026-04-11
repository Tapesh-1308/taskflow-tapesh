package project

import "time"

type Project struct {
	ID          string
	Name        string
	Description *string
	OwnerID     string
	CreatedAt   time.Time
}

type ProjectWithTasks struct {
	Project
	Tasks []Task
}

type Task struct {
	ID         string
	Title      string
	Status     string
	AssigneeID *string
}
