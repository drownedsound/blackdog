package app

import "time"

type CreateApplicationRequest struct {
	Birthday          string `json:"birthday"`
	MemberReferenceNo string `json:"member_reference_no"`
	CardProfile       int64  `json:"card_profile,omitempty"`
	LastName          string `json:"last_name"`
	FirstName         string `json:"first_name"`
	MiddleName        string `json:"middle_name,omitempty"`
	LoanProfile       int64  `json:"loan_profile,omitempty"`
	RequestedAmount   int    `json:"requested_amount"`
}

type CreateApplicationResponse struct {
	Birthday          string            `json:"birthday"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	MemberReferenceNo string            `json:"member_reference_no"`
	LastName          string            `json:"last_name"`
	FirstName         string            `json:"first_name"`
	MiddleName        string            `json:"middle_name,omitempty"`
	Status            ApplicationStatus `json:"status"`
	CardProfile       int64             `json:"card_profile,omitempty"`
	LoanProfile       int64             `json:"loan_profile,omitempty"`
	Id                int64             `json:"id"`
	RequestedAmount   int               `json:"requested_amount"`
	InterestRate      int               `json:"interest_rate"`
}

// type GetApplicationRequest struct {
// 	Id int64 `json:"id"`
// }

// type GetApplicationResponse struct {
// 	CreatedAt         time.Time         `json:"created_at"`
// 	UpdatedAt         time.Time         `json:"updated_at"`
// 	MemberReferenceNo string            `json:"member_reference_no"`
// 	Status            ApplicationStatus `json:"status"`
// 	// CardProfile       int64             `json:"card_profile"`
// 	Id              int64 `json:"id"`
// 	RequestedAmount int   `json:"requested_amount"`
// 	// InterestRate      int               `json:"interest_rate"`
// }
