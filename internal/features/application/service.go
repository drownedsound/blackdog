package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// FIXME: Implement slog.LogValuer to log complex object graphs
// FIXME: Group related fields using slog.GroupValue
func (s *Service) CreateApplication(
	ctx context.Context, req CreateApplicationRequest,
) (CreateApplicationResponse, error) {
	app := &Application{
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		MemberReferenceNo: req.MemberReferenceNo,
		CategoryCode:      req.CategoryCode,
		StatusCode:        StatusCreated,
		RequestedAmount:   req.RequestedAmount,
	}

	if err := app.Validate(); err != nil {
		return CreateApplicationResponse{}, fmt.Errorf(
			"application.service stopped saving entity: %w", err,
		)
	}

	if err := s.repo.Save(ctx, app); err != nil {
		return CreateApplicationResponse{}, fmt.Errorf(
			"application.service failed to save entity: %w", err,
		)
	}

	return CreateApplicationResponse{
		CreatedAt:         app.CreatedAt,
		UpdatedAt:         app.UpdatedAt,
		MemberReferenceNo: app.MemberReferenceNo,
		CategoryCode:      app.CategoryCode,
		StatusCode:        app.StatusCode.String(),
		Id:                app.Id,
		RequestedAmount:   app.RequestedAmount,
	}, nil
}

func (s *Service) GetApplicationById(
	ctx context.Context, req GetApplicationRequest,
) (GetApplicationResponse, error) {
	app, err := s.repo.GetById(ctx, req.Id)
	if err != nil {
		return GetApplicationResponse{}, fmt.Errorf(
			"application.service failed to get application: %w", err,
		)
	}

	return GetApplicationResponse{
		CreatedAt:         app.CreatedAt,
		UpdatedAt:         app.UpdatedAt,
		MemberReferenceNo: app.MemberReferenceNo,
		CategoryCode:      app.CategoryCode,
		StatusCode:        app.StatusCode.String(),
		Id:                app.Id,
		RequestedAmount:   app.RequestedAmount,
	}, nil
}
