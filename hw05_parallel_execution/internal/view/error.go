package view

import "errors"

var (
	ErrErrorsLimitExceeded  = errors.New("errors limit exceeded")
	ErrErrorsLessUnitWorker = errors.New("errors goroutine less than unit")
)
