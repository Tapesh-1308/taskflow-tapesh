package task

import (
	"context"
	"errors"
)

type Service interface {
	Create(ctx context.Context, title string, projectID string, userID string) (*Task, error)
	UpdateStatus(ctx context.Context, id string, status string, userID string) error
	Assign(ctx context.Context, id string, assignTo string, userID string) error
	List(ctx context.Context, projectID string, status *string, userID string) ([]Task, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(ctx context.Context, title string, projectID string, userID string) (*Task, error) {
	t := &Task{
		Title:     title,
		Status:    "pending",
		ProjectID: projectID,
	}

	err := s.repo.Create(ctx, t)
	return t, err
}

func (s *service) UpdateStatus(ctx context.Context, id string, status string, userID string) error {
	// basic validation
	if status != "pending" && status != "in_progress" && status != "done" {
		return errors.New("invalid status")
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *service) Assign(ctx context.Context, id string, assignTo string, userID string) error {
	return s.repo.Assign(ctx, id, assignTo)
}

func (s *service) List(ctx context.Context, projectID string, status *string, userID string) ([]Task, error) {
	return s.repo.GetByProject(ctx, projectID, status)
}
