package app

import "time"

// const (
// 	StatusUnknown ApplicationStatus = iota
// 	StatusCreated
// 	StatusInProgress
// 	StatusApproved
// 	StatusDeclined
// 	StatusCancelled
// )

// type CategoryCode string

// func (c CategoryCode) Parse() ProductCategory {
// 	switch c {
// 	case "CREDIT_CARD":
// 		return CategoryCard
// 	case "PERSONAL_LOAN":
// 		return CategoryLoan
// 	default:
// 		return CategoryUnknown
// 	}
// }

type CreateApplicationRequest struct {
	MemberReferenceNo string `json:"member_reference_no"`
	// CategoryCode      CategoryCode `json:"category_code"`
	CardProfileCode int64 `json:"card_profile_code"`
	RequestedAmount int   `json:"requested_amount"`
}

type CreateApplicationResponse struct {
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	MemberReferenceNo string    `json:"member_reference_no"`
	// CategoryCode      string    `json:"category_code"`
	StatusCode      ApplicationStatus `json:"status_code"`
	CardProfileCode int64             `json:"card_profile_code"`
	Id              int64             `json:"id"`
	RequestedAmount int               `json:"requested_amount"`
	InterestRate    int               `json:"interest_rate"`
}

type GetApplicationRequest struct {
	Id int64 `json:"id"`
}

type GetApplicationResponse struct {
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	MemberReferenceNo string    `json:"member_reference_no"`
	// CategoryCode      string    `json:"category_code"`
	StatusCode      ApplicationStatus `json:"status_code"`
	CardProfileCode int64             `json:"card_profile_code"`
	Id              int64             `json:"id"`
	RequestedAmount int               `json:"requested_amount"`
	InterestRate    int               `json:"interest_rate"`
}
