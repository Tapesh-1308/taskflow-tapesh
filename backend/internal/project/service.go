package project

import (
	"context"
	"errors"

	"log/slog"
)

type Service interface {
	Create(ctx context.Context, name string, desc *string, userID string) (*Project, error)
	List(ctx context.Context, userID string) ([]ProjectWithUser, error)
	Update(ctx context.Context, id string, name string, desc *string, userID string) error
	Delete(ctx context.Context, id string, userID string) error
	GetByID(ctx context.Context, id string, userID string) (*ProjectWithTasks, error)
}

type service struct {
	repo Repository
	log  *slog.Logger
}

func NewService(repo Repository, log *slog.Logger) Service {
	return &service{repo: repo, log: log}
}

func (s *service) Create(ctx context.Context, name string, desc *string, userID string) (*Project, error) {
	s.log.Info("Creating project", "name", name, "userID", userID)

	p := &Project{
		Name:        name,
		Description: desc,
		OwnerID:     userID,
	}

	err := s.repo.Create(ctx, p)
	if err != nil {
		s.log.Error("Failed to create project", "error", err, "name", name, "userID", userID)
		return nil, err
	}

	s.log.Info("Project created successfully", "id", p.ID, "name", name)
	return p, nil
}

func (s *service) List(ctx context.Context, userID string) ([]ProjectWithUser, error) {
	s.log.Info("Listing projects for user", "userID", userID)

	projects, err := s.repo.GetUserProjects(ctx, userID)
	if err != nil {
		s.log.Error("Failed to list projects", "error", err, "userID", userID)
		return nil, err
	}

	s.log.Info("Projects listed successfully", "count", len(projects), "userID", userID)
	return projects, nil
}

func (s *service) Update(ctx context.Context, id string, name string, desc *string, userID string) error {
	s.log.Info("Updating project", "id", id, "name", name, "userID", userID)

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("Failed to get project for update", "error", err, "id", id)
		return err
	}

	if p.Owner.ID != userID {
		s.log.Warn("Forbidden: user not owner", "userID", userID, "projectID", id, "ownerID", p.Owner.ID)
		return errors.New("forbidden")
	}

	p.Name = name
	p.Description = desc

	err = s.repo.Update(ctx, p)
	if err != nil {
		s.log.Error("Failed to update project", "error", err, "id", id)
		return err
	}

	s.log.Info("Project updated successfully", "id", id)
	return nil
}

func (s *service) Delete(ctx context.Context, id string, userID string) error {
	s.log.Info("Deleting project", "id", id, "userID", userID)

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("Failed to get project for deletion", "error", err, "id", id)
		return err
	}

	if p.Owner.ID != userID {
		s.log.Warn("Forbidden: user not owner", "userID", userID, "projectID", id, "ownerID", p.Owner.ID)
		return errors.New("forbidden")
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.log.Error("Failed to delete project", "error", err, "id", id)
		return err
	}

	s.log.Info("Project deleted successfully", "id", id)
	return nil
}

func (s *service) GetByID(ctx context.Context, id string, userID string) (*ProjectWithTasks, error) {
	s.log.Info("Getting project by ID", "id", id, "userID", userID)

	project, err := s.repo.GetProjectWithTasks(ctx, id)
	if err != nil {
		s.log.Error("Failed to get project", "error", err, "id", id)
		return nil, err
	}

	// allow if owner OR assigned to any task
	if project.Owner.ID != userID {
		allowed := false

		for _, t := range project.Tasks {
			if t.Assignee != nil && t.Assignee.ID == userID {
				allowed = true
				break
			}
		}

		if !allowed {
			s.log.Warn("Forbidden: user not authorized", "userID", userID, "projectID", id, "ownerID", project.Owner.ID)
			return nil, errors.New("forbidden")
		}
	}

	s.log.Info("Project retrieved successfully", "id", id, "taskCount", len(project.Tasks))
	return project, nil
}
