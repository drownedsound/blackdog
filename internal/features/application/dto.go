package app

import "time"

type CreateApplicationRequest struct {
	MemberReferenceNo string `json:"member_reference_no"`
	CardProfile       int64  `json:"card_profile"`
	RequestedAmount   int    `json:"requested_amount"`
}

type CreateApplicationResponse struct {
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	MemberReferenceNo string            `json:"member_reference_no"`
	Status            ApplicationStatus `json:"status"`
	CardProfile       int64             `json:"card_profile"`
	Id                int64             `json:"id"`
	RequestedAmount   int               `json:"requested_amount"`
	InterestRate      int               `json:"interest_rate"`
}

type GetApplicationRequest struct {
	Id int64 `json:"id"`
}

type GetApplicationResponse struct {
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	MemberReferenceNo string            `json:"member_reference_no"`
	Status            ApplicationStatus `json:"status"`
	CardProfile       int64             `json:"card_profile"`
	Id                int64             `json:"id"`
	RequestedAmount   int               `json:"requested_amount"`
	InterestRate      int               `json:"interest_rate"`
}
