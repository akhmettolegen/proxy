package drivers

import "errors"

var (
	ErrTaskEmpty    = errors.New("task is empty")
	ErrTaskNotFound = errors.New("task not found")
)
