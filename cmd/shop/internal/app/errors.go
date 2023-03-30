package app

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrSameStatus      = errors.New("order has same status")
	ErrNotFound        = errors.New("not found")
)
