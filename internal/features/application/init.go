package app

import "time"

var validator *Validator

func init()  {
	validator = NewValidator(time.Now().UTC())
}

