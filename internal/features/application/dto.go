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

type CreateApplicationRequest struct {
	Birthday              string                 `json:"birthday"`
	MemberReferenceNumber string                 `json:"member_reference_no"`
	CardProfile           int64                  `json:"card_profile,omitempty"`
	LastName              string                 `json:"last_name"`
	FirstName             string                 `json:"first_name"`
	MiddleName            string                 `json:"middle_name,omitempty"`
	LoanProfile           int64                  `json:"loan_profile,omitempty"`
	RequestedAmount       int32                  `json:"requested_amount"`
	ContactNumbers        []ContactNumberRequest `json:"contact_numbers"`
	OtherApplicants       []ApplicantRequest     `json:"other_applicants,omitempty"`
}

type CreateApplicationResponse struct {
	CreatedAt             time.Time               `json:"created_at"`
	UpdatedAt             time.Time               `json:"updated_at"`
	MemberReferenceNumber string                  `json:"member_reference_no"`
	LastName              string                  `json:"last_name"`
	FirstName             string                  `json:"first_name"`
	MiddleName            string                  `json:"middle_name,omitempty"`
	Birthday              string                  `json:"birthday"`
	ContactNumbers        []ContactNumberResponse `json:"contact_numbers"`
	OtherApplicants       []ApplicantResponse     `json:"other_applicants,omitempty"`
	CardProfile           int64                   `json:"card_profile,omitempty"`
	LoanProfile           int64                   `json:"loan_profile,omitempty"`
	Id                    int64                   `json:"id"`
	CreditLimit           int32                   `json:"credit_limit,omitempty"`
	LoanAmount            int32                   `json:"loan_amount,omitempty"`
	RequestedAmount       int32                   `json:"requested_amount"`
	InterestRate          int16                   `json:"interest_rate"`
	Status                ApplicationStatus       `json:"status"`
}

type GetApplicationRequest struct {
	Id int64 `json:"id"`
}

type GetApplicationResponse struct {
	CreatedAt             time.Time               `json:"created_at"`
	UpdatedAt             time.Time               `json:"updated_at"`
	MemberReferenceNumber string                  `json:"member_reference_no"`
	LastName              string                  `json:"last_name"`
	FirstName             string                  `json:"first_name"`
	MiddleName            string                  `json:"middle_name,omitempty"`
	Birthday              string                  `json:"birthday"`
	ContactNumbers        []ContactNumberResponse `json:"contact_numbers"`
	OtherApplicants       []ApplicantResponse     `json:"other_applicants,omitempty"`
	CardProfile           int64                   `json:"card_profile,omitempty"`
	LoanProfile           int64                   `json:"loan_profile,omitempty"`
	Id                    int64                   `json:"id"`
	CreditLimit           int32                   `json:"credit_limit,omitempty"`
	LoanAmount            int32                   `json:"loan_amount,omitempty"`
	RequestedAmount       int32                   `json:"requested_amount"`
	InterestRate          int16                   `json:"interest_rate"`
	Status                ApplicationStatus       `json:"status"`
}
