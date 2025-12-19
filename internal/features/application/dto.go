package app

import "time"

type CategoryCode string

func (c CategoryCode) Parse() ProductCategory {
	switch c {
	case "CREDIT_CARD":
		return CategoryCard
	case "PERSONAL_LOAN":
		return CategoryLoan
	default:
		return CategoryUnknown
	}
}

type CreateApplicationRequest struct {
	MemberReferenceNo string       `json:"member_reference_no"`
	CategoryCode      CategoryCode `json:"category_code"`
	RequestedAmount   int          `json:"requested_amount"`
}

type CreateApplicationResponse struct {
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	MemberReferenceNo string    `json:"member_reference_no"`
	CategoryCode      string    `json:"category_code"`
	StatusCode        string    `json:"status_code"`
	Id                int64     `json:"id"`
	RequestedAmount   int       `json:"requested_amount"`
}

type GetApplicationRequest struct {
	Id int64 `json:"id"`
}

type GetApplicationResponse struct {
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	MemberReferenceNo string    `json:"member_reference_no"`
	CategoryCode      string    `json:"category_code"`
	StatusCode        string    `json:"status_code"`
	Id                int64     `json:"id"`
	RequestedAmount   int       `json:"requested_amount"`
}
