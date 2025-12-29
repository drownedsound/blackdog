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
	validator := newValidator(time.Now().UTC())

	principalReq := ApplicantRequest{
		Birthday:       req.Birthday,
		LastName:       req.LastName,
		FirstName:      req.FirstName,
		MiddleName:     req.MiddleName,
		ContactNumbers: req.ContactNumbers,
	}

	principal, err := s.mapApplicant(principalReq, true)
	if err != nil {
		return CreateApplicationResponse{}, fmt.Errorf("principal mapping failed: %w", err)
	}

	otherApplicants := make([]Applicant, 0, len(req.OtherApplicants))
	for i, oaReq := range req.OtherApplicants {
		oa, err := s.mapApplicant(oaReq, false)
		if err != nil {
			return CreateApplicationResponse{}, fmt.Errorf("other applicant %d mapping failed: %w", i, err)
		}
		otherApplicants = append(otherApplicants, oa)
	}

	app := Application{
		CreatedAt:             now,
		UpdatedAt:             now,
		OtherApplicants:       otherApplicants,
		MemberReferenceNumber: req.MemberReferenceNumber,
		Applicant:             principal,
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
			Type:  int(c.Type),
		}
	}

	otherApplicantsRes := make([]ApplicantResponse, len(app.OtherApplicants))
	for i, oa := range app.OtherApplicants {
		// No error check needed here as this is a pure transformation
		// of valid data
		otherApplicantsRes[i] = s.mapApplicantToResponse(oa)
	}

	res := CreateApplicationResponse{
		Birthday:              app.Birthday.Format(time.DateOnly),
		CreatedAt:             app.CreatedAt,
		UpdatedAt:             app.UpdatedAt,
		MemberReferenceNumber: app.MemberReferenceNumber,
		LastName:              app.LastName,
		FirstName:             app.FirstName,
		MiddleName:            app.MiddleName,
		CardProfile:           app.CreditCard.ProfileId,
		CreditLimit:           app.CreditCard.CreditLimit,
		LoanProfile:           app.PersonalLoan.ProfileId,
		LoanAmount:            app.PersonalLoan.LoanAmount,
		Status:                app.Status,
		Id:                    app.Id,
		RequestedAmount:       app.RequestedAmount,
		ContactNumbers:        contactNumbersRes,
		OtherApplicants:       otherApplicantsRes,
	}

	if app.CreditCard.ProfileId > 0 {
		res.InterestRate = app.CreditCard.InterestRate
	}

	if app.PersonalLoan.ProfileId > 0 {
		res.InterestRate = app.PersonalLoan.InterestRate
	}

	return res, nil
}

// mapApplicant transforms the DTO into a Domain Entity.
// It accepts the struct by value to avoid pointer chasing in the stack.
func (s *Service) mapApplicant(req ApplicantRequest, isPrincipal bool) (Applicant, error) {
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		return Applicant{}, fmt.Errorf("invalid birthday format: %w", err)
	}

	contacts := make([]ContactNumber, len(req.ContactNumbers))
	for i, c := range req.ContactNumbers {
		contacts[i] = ContactNumber{
			Value: c.Value,
			Type:  ContactNumberType(c.Type),
		}
	}

	return Applicant{
		Birthday:       birthday,
		ContactNumbers: contacts,
		LastName:       req.LastName,
		FirstName:      req.FirstName,
		MiddleName:     req.MiddleName,
		IsPrincipal:    isPrincipal,
	}, nil
}

// mapApplicantToResponse transforms the Domain Entity to the DTO.
// It handles the specific formatting logic (e.g., DateOnly) efficiently.
func (s *Service) mapApplicantToResponse(a Applicant) ApplicantResponse {
	// Pre-allocate contact numbers to avoid resize
	contacts := make([]ContactNumberResponse, len(a.ContactNumbers))
	for i, c := range a.ContactNumbers {
		contacts[i] = ContactNumberResponse{
			Value: c.Value,
			Type:  int(c.Type),
		}
	}

	return ApplicantResponse{
		Birthday:       a.Birthday.Format(time.DateOnly), // Format as YYYY-MM-DD
		LastName:       a.LastName,
		FirstName:      a.FirstName,
		MiddleName:     a.MiddleName,
		ContactNumbers: contacts,
		IsPrincipal:    a.IsPrincipal,
	}
}
