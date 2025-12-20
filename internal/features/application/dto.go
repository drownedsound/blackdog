package app

import "time"

type CreateApplicationRequest struct {
	MemberReferenceNo string `json:"member_reference_no"`
	CardProfileCode   int64  `json:"card_profile_code"`
	RequestedAmount   int    `json:"requested_amount"`
}

type CreateApplicationResponse struct {
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	MemberReferenceNo string            `json:"member_reference_no"`
	StatusCode        ApplicationStatus `json:"status_code"`
	CardProfileCode   int64             `json:"card_profile_code"`
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
	StatusCode        ApplicationStatus `json:"status_code"`
	CardProfileCode   int64             `json:"card_profile_code"`
	Id                int64             `json:"id"`
	RequestedAmount   int               `json:"requested_amount"`
	InterestRate      int               `json:"interest_rate"`
}
