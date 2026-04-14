package task

import (
	"context"
	"errors"

	"log/slog"
)

type Service interface {
	Create(ctx context.Context, title string, projectID string, userID string) (*Task, error)
	ListByProject(ctx context.Context, projectID string, status *string, assignee *string, userID string) ([]TaskWithUser, error)
	Update(ctx context.Context, id string, body UpdateTaskInput, userID string) error
	Delete(ctx context.Context, id string, userID string) error
}

type service struct {
	repo Repository
	log  *slog.Logger
}

func NewService(r Repository, log *slog.Logger) Service {
	return &service{repo: r, log: log}
}

func (s *service) Create(ctx context.Context, title string, projectID string, userID string) (*Task, error) {
	s.log.Info("Creating task", "title", title, "projectID", projectID, "userID", userID)

	t := &Task{
		Title:     title,
		Status:    "todo",
		ProjectID: projectID,
	}

	err := s.repo.Create(ctx, t)
	if err != nil {
		s.log.Error("Failed to create task", "error", err, "title", title, "projectID", projectID)
		return nil, err
	}

	s.log.Info("Task created successfully", "id", t.ID, "title", title)
	return t, nil
}

func (s *service) ListByProject(
	ctx context.Context,
	projectID string,
	status *string,
	assignee *string,
	userID string,
) ([]TaskWithUser, error) {
	s.log.Info("Listing tasks by project", "projectID", projectID, "status", status, "assignee", assignee, "userID", userID)

	tasks, err := s.repo.GetByProject(ctx, projectID, status, assignee)
	if err != nil {
		s.log.Error("Failed to list tasks", "error", err, "projectID", projectID)
		return nil, err
	}

	s.log.Info("Tasks listed successfully", "count", len(tasks), "projectID", projectID)
	return tasks, nil
}

func (s *service) Update(ctx context.Context, id string, body UpdateTaskInput, userID string) error {
	s.log.Info("Updating task", "id", id, "userID", userID)

	if body.Status != nil {
		valid := map[string]bool{
			"todo": true, "in_progress": true, "done": true,
		}
		if !valid[*body.Status] {
			s.log.Warn("Invalid status provided", "status", *body.Status, "id", id)
			return errors.New("invalid status")
		}
	}

	err := s.repo.Update(ctx, id, body)
	if err != nil {
		s.log.Error("Failed to update task", "error", err, "id", id)
		return err
	}

	s.log.Info("Task updated successfully", "id", id)
	return nil
}

func (s *service) Delete(ctx context.Context, id string, userID string) error {
	s.log.Info("Deleting task", "id", id, "userID", userID)

	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		s.log.Error("Failed to delete task", "error", err, "id", id)
		return err
	}

	s.log.Info("Task deleted successfully", "id", id)
	return nil
}
