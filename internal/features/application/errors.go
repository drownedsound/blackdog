package app

import "errors"

var (
	ErrInvalidCreditCard = errors.New(
		"product (credit card) is not a valid value")
	ErrInvalidCreditLimit = errors.New(
		"credit limit cannot be less than zero")
	ErrInvalidPersonalLoan = errors.New(
		"product (personal loan) is not a valid value")
	ErrInvalidLoanAmount = errors.New(
		"loan amount cannot be less than zero")
	ErrInvalidCurrency = errors.New(
		"currency used is not a valid value")
	ErrInvalidInterestRate = errors.New(
		"interest rate cannot be less than or equal to zero")
	ErrMissingBirthday = errors.New(
		"birthday is required")
	ErrBirthdayInFuture = errors.New(
		"birthday cannot be in the future")
	ErrMinimumAgeNotMet = errors.New(
		"minimum age not met")
	ErrMissingLastName = errors.New(
		"last name is required")
	ErrLastNameTooLong = errors.New(
		"last name cannot be more than 30 characters")
	ErrMissingFirstName = errors.New(
		"first name is required")
	ErrFirstNameTooLong = errors.New(
		"first name cannot be more than 30 characters")
	ErrMiddleNameTooLong = errors.New(
		"middle name cannot be more than 30 characters")
	ErrMissingContactNumber = errors.New(
		"contact number is required")
	ErrInvalidContactNumber = errors.New(
		"contact number is not a valid value")
	ErrInvalidContactNumberType = errors.New(
		"contact number type is not a valid value")
	ErrCreatedAtInFuture = errors.New(
		"created at cannot be in the future")
	ErrUpdatedAtInFuture = errors.New(
		"updated at cannot be in the future")
	ErrMissingMemberReferenceNumber = errors.New(
		"member reference number is required")
	ErrInvalidRequestedAmount = errors.New(
		"requested amount cannot be less than or equal to zero")
	ErrInvalidStatus = errors.New(
		"status is not a valid value")
	ErrMissingProduct = errors.New(
		"credit card or personal loan is required")
	ErrTooManyProducts = errors.New(
		"application limited to one product only")
	ErrNotFound = errors.New(
		"database record not found")
	ErrInsertFailed = errors.New(
		"database record insertion failed")
	ErrConnectionRefused = errors.New(
		"database connection refused")
	ErrInvalidUUID = errors.New(
		"uuid cannot be less than or equal to zero")
	ErrInvalidNodeId = errors.New(
		"node id out of range")
)
