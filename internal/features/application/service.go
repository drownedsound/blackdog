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
	ctx context.Context,
	req CreateApplicationRequest,
) (CreateApplicationResponse, error) {

	now := time.Now().UTC()
	// app := &Application{
	// 	CreditCard: CreditCard{
	// 		CardProfile: req.CardProfile,
	// 		// TODO: Retrieve matching InterestRate based on CardProfile
	// 		InterestRate: 1_250,
	// 	},
	// 	CreatedAt:         now,
	// 	UpdatedAt:         now,
	// 	MemberReferenceNo: req.MemberReferenceNo,
	// 	Status:            StatusCreated,
	// 	RequestedAmount:   req.RequestedAmount,
	// }
	app := Application{
		CreatedAt:         now,
		UpdatedAt:         now,
		OtherApplicants: []Applicant{},
		MemberReferenceNo: "ABCDE12345",
		Applicant: Applicant {
			Birthday: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC), 
			ContactNumbers: []ContactNumber {
				{ Value: "9171234567", Type: TypeMobile },
			},
			LastName: "Smith",
			FirstName: "John",
			MiddleName: "Doe",
			IsPrincipal: true,
		},
		CreditCard: CreditCard {
			ProfileId: 1,
			CurrencyId: 1,
			CreditLimit: 1_000_000,
			InterestRate: 300,
		},
		Status:            StatusCreated,
		RequestedAmount:   100_000_000,
	}

	if err := app.Validate(validator); err != nil {
		return CreateApplicationResponse{}, fmt.Errorf(
			"application.service stopped saving entity: %w", err,
		)
	}

	if err := s.repo.Save(ctx, &app); err != nil {
		return CreateApplicationResponse{}, fmt.Errorf(
			"application.service failed to save entity: %w", err,
		)
	}

	return CreateApplicationResponse{
		CreatedAt:         app.CreatedAt,
		UpdatedAt:         app.UpdatedAt,
		MemberReferenceNo: app.MemberReferenceNo,
		// CardProfile:       app.CardProfile,
		Status:            app.Status,
		Id:                app.Id,
		RequestedAmount:   app.RequestedAmount,
		// InterestRate:      app.InterestRate,
	}, nil
}

func (s *Service) GetApplicationById(
	ctx context.Context,
	req GetApplicationRequest,
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
		// CardProfile:       app.CardProfile,
		Status:            app.Status,
		Id:                app.Id,
		RequestedAmount:   app.RequestedAmount,
		// InterestRate:      app.InterestRate,
	}, nil
}
