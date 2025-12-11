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

func (s *Service) CreateApplication(
	ctx context.Context, req CreateApplicationRequest,
) (*CreateApplicationResponse, error) {
	app := &Application{
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		MemberReferenceNo: req.MemberReferenceNo,
		CategoryCode:      req.CategoryCode,
		StatusCode:        "CREATED",
		RequestedAmount:   req.RequestedAmount,
	}

	// FIXME: Implement slog.LogValuer to log complex object graphs
	// FIXME: Group related fields to create structured subgraphs using
	//        slog.GroupValue
	s.logger.Debug(
		"Persisting transaction to the database",
		slog.Any("txn", app),
	)
	if err := s.repo.Save(ctx, app); err != nil {
		return nil, fmt.Errorf("failed to save application: %w", err)
	}

	return &CreateApplicationResponse{
		CreatedAt:         app.CreatedAt,
		UpdatedAt:         app.UpdatedAt,
		MemberReferenceNo: app.MemberReferenceNo,
		CategoryCode:      app.CategoryCode,
		StatusCode:        app.StatusCode,
		Id:                app.Id,
		RequestedAmount:   app.RequestedAmount,
	}, nil
}

func (s *Service) GetApplicationById(
	ctx context.Context, req GetApplicationRequest,
) (*GetApplicationResponse, error) {
	app, err := s.repo.GetById(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	return &GetApplicationResponse{
		CreatedAt:         app.CreatedAt,
		UpdatedAt:         app.UpdatedAt,
		MemberReferenceNo: app.MemberReferenceNo,
		CategoryCode:      app.CategoryCode,
		StatusCode:        app.StatusCode,
		Id:                app.Id,
		RequestedAmount:   app.RequestedAmount,
	}, nil
}
