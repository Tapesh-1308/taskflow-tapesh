package project

import (
	"context"
	"errors"
)

type Service interface {
	Create(ctx context.Context, name string, desc *string, userID string) (*Project, error)
	List(ctx context.Context, userID string) ([]Project, error)
	Update(ctx context.Context, id string, name string, desc *string, userID string) error
	Delete(ctx context.Context, id string, userID string) error
	GetByID(ctx context.Context, id string, userID string) (*ProjectWithTasks, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, name string, desc *string, userID string) (*Project, error) {
	p := &Project{
		Name:        name,
		Description: desc,
		OwnerID:     userID,
	}

	err := s.repo.Create(ctx, p)
	return p, err
}

func (s *service) List(ctx context.Context, userID string) ([]Project, error) {
	return s.repo.GetUserProjects(ctx, userID)
}

func (s *service) Update(ctx context.Context, id string, name string, desc *string, userID string) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if p.OwnerID != userID {
		return errors.New("forbidden")
	}

	p.Name = name
	p.Description = desc

	return s.repo.Update(ctx, p)
}

func (s *service) Delete(ctx context.Context, id string, userID string) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if p.OwnerID != userID {
		return errors.New("forbidden")
	}

	return s.repo.Delete(ctx, id)
}

func (s *service) GetByID(ctx context.Context, id string, userID string) (*ProjectWithTasks, error) {
	project, err := s.repo.GetProjectWithTasks(ctx, id)
	if err != nil {
		return nil, err
	}

	// allow if owner OR assigned to any task
	if project.OwnerID != userID {
		allowed := false

		for _, t := range project.Tasks {
			if t.AssigneeID != nil && *t.AssigneeID == userID {
				allowed = true
				break
			}
		}

		if !allowed {
			return nil, errors.New("forbidden")
		}
	}

	return project, nil
}
