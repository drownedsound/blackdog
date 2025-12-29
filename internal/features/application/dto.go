package app

import "time"

type ContactNumberRequest struct {
	Value string `json:"value"`
	Type  int    `json:"type"`
}

type ContactNumberResponse struct {
	Value string `json:"value"`
	Type  int    `json:"type"`
}

type ApplicantRequest struct {
	Birthday       string                 `json:"birthday"`
	LastName       string                 `json:"last_name"`
	FirstName      string                 `json:"first_name"`
	MiddleName     string                 `json:"middle_name,omitempty"`
	ContactNumbers []ContactNumberRequest `json:"contact_numbers"`
}

type ApplicantResponse struct {
	Birthday       string                  `json:"birthday"`
	LastName       string                  `json:"last_name"`
	FirstName      string                  `json:"first_name"`
	MiddleName     string                  `json:"middle_name,omitempty"`
	ContactNumbers []ContactNumberResponse `json:"contact_numbers"`
	IsPrincipal    bool                    `json:"is_principal"`
}

// TODO: Rename to ApplicationRequest
type CreateApplicationRequest struct {
	Birthday string `json:"birthday"`
	// FIXME: Use Number instead of No
	MemberReferenceNumber string                 `json:"member_reference_no"`
	CardProfile       int64                  `json:"card_profile,omitempty"`
	LastName          string                 `json:"last_name"`
	FirstName         string                 `json:"first_name"`
	MiddleName        string                 `json:"middle_name,omitempty"`
	LoanProfile       int64                  `json:"loan_profile,omitempty"`
	RequestedAmount   int                    `json:"requested_amount"`
	ContactNumbers    []ContactNumberRequest `json:"contact_numbers"`
	OtherApplicants   []ApplicantRequest     `json:"other_applicants,omitempty"`
}

// TODO: Rename to ApplicationResponse
type CreateApplicationResponse struct {
	Birthday          string                  `json:"birthday"`
	CreatedAt         time.Time               `json:"created_at"`
	UpdatedAt         time.Time               `json:"updated_at"`
	MemberReferenceNumber string                  `json:"member_reference_no"`
	LastName          string                  `json:"last_name"`
	FirstName         string                  `json:"first_name"`
	MiddleName        string                  `json:"middle_name,omitempty"`
	Status            ApplicationStatus       `json:"status"`
	CardProfile       int64                   `json:"card_profile,omitempty"`
	LoanProfile       int64                   `json:"loan_profile,omitempty"`
	Id                int64                   `json:"id"`
	RequestedAmount   int                     `json:"requested_amount"`
	InterestRate      int                     `json:"interest_rate"`
	ContactNumbers    []ContactNumberResponse `json:"contact_numbers"`
	OtherApplicants   []ApplicantResponse     `json:"other_applicants,omitempty"`
}

// type GetApplicationRequest struct {
// 	Id int64 `json:"id"`
// }

// type GetApplicationResponse struct {
// 	CreatedAt         time.Time         `json:"created_at"`
// 	UpdatedAt         time.Time         `json:"updated_at"`
// 	MemberReferenceNumber string            `json:"member_reference_no"`
// 	Status            ApplicationStatus `json:"status"`
// 	// CardProfile       int64             `json:"card_profile"`
// 	Id              int64 `json:"id"`
// 	RequestedAmount int   `json:"requested_amount"`
// 	// InterestRate      int               `json:"interest_rate"`
// }
