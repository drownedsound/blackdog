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
	validator := NewValidator(time.Now().UTC())

	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		return CreateApplicationResponse{}, fmt.Errorf(
			"invalid birthday format (expected YYYY-MM-DD): %w", err,
		)
	}

	contactNumbers := make([]ContactNumber, len(req.ContactNumbers))
	for i, c := range req.ContactNumbers {
		contactNumbers[i] = ContactNumber{
			Value: c.Value,
			Type:  ContactNumberType(c.Type),
		}
	}

	app := Application{
		CreatedAt: now,
		UpdatedAt: now,
		// TODO: Map from request
		OtherApplicants:   []Applicant{},
		MemberReferenceNo: req.MemberReferenceNo,
		// TODO: Map from request
		Applicant: Applicant{
			Birthday:       birthday,
			ContactNumbers: contactNumbers,
			LastName:       req.LastName,
			FirstName:      req.FirstName,
			MiddleName:     req.MiddleName,
			IsPrincipal:    true,
		},
		CreditCard: CreditCard{
			ProfileId: req.CardProfile,
			// TODO: Get CurrencyId from cache
			CurrencyId:  1,
			CreditLimit: 0,
			// TODO: Get InterestRate from cache
			InterestRate: 300,
		},
		PersonalLoan: PersonalLoan{
			ProfileId: req.LoanProfile,
			// TODO: Get CurrencyId from cache
			CurrencyId: 1,
			LoanAmount: 0,
			// TODO: Get InterestRate from cache
			InterestRate: 300,
		},
		Status:          StatusCreated,
		RequestedAmount: req.RequestedAmount,
	}

	if err := app.Validate(validator); err != nil {
		return CreateApplicationResponse{}, fmt.Errorf(
			"application.service stopped validating entity: %w", err,
		)
	}

	if err := s.repo.Insert(ctx, &app); err != nil {
		return CreateApplicationResponse{}, fmt.Errorf(
			"application.service failed to save entity: %w", err,
		)
	}

	contactNumbersRes := make([]ContactNumberResponse, len(app.ContactNumbers))
	for i, c := range app.ContactNumbers {
		contactNumbersRes[i] = ContactNumberResponse{
			Value: c.Value,
			// Cast the domain type (byte/enum) to the DTO type (int)
			Type: int(c.Type),
		}
	}

	res := CreateApplicationResponse{
		Birthday:          app.Birthday.Format(time.DateOnly),
		CreatedAt:         app.CreatedAt,
		UpdatedAt:         app.UpdatedAt,
		MemberReferenceNo: app.MemberReferenceNo,
		LastName:          app.LastName,
		FirstName:         app.FirstName,
		MiddleName:        app.MiddleName,
		CardProfile:       app.CreditCard.ProfileId,
		LoanProfile:       app.PersonalLoan.ProfileId,
		Status:            app.Status,
		Id:                app.Id,
		RequestedAmount:   app.RequestedAmount,
		ContactNumbers:    contactNumbersRes,
	}

	if app.CreditCard.ProfileId > 0 {
		res.InterestRate = app.CreditCard.InterestRate
	}

	if app.PersonalLoan.ProfileId > 0 {
		res.InterestRate = app.PersonalLoan.InterestRate
	}

	return res, nil
}

// func (s *Service) GetApplicationById(
// 	ctx context.Context,
// 	req GetApplicationRequest,
// ) (GetApplicationResponse, error) {
// 	app, err := s.repo.GetByInternalId(ctx, req.Id)
// 	if err != nil {
// 		return GetApplicationResponse{}, fmt.Errorf(
// 			"application.service failed to get application: %w", err,
// 		)
// 	}
//
// 	return GetApplicationResponse{
// 		CreatedAt:         app.CreatedAt,
// 		UpdatedAt:         app.UpdatedAt,
// 		MemberReferenceNo: app.MemberReferenceNo,
// 		// CardProfile:       app.CardProfile,
// 		Status:          app.Status,
// 		Id:              app.Id,
// 		RequestedAmount: app.RequestedAmount,
// 		// InterestRate:      app.InterestRate,
// 	}, nil
// }
