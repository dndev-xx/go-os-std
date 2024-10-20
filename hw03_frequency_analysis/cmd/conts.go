package cmd

import (
	"errors"
	"regexp"
)

var (
	ErrorEmpty = errors.New("slice is empty")
	regex      = regexp.MustCompile(REGEX)
)

const (
	NEGATIVE = -1
	REGEX    = `\S+`
	MAX      = 10
	ZERO     = 0
	DIV      = 2
)
