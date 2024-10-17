package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid input string")

const (
	EMPTY = ""
	ZERO  = 0
	UNIT  = 1
	TWO   = 2
)

func Unpack(input string) (string, error) {
	if input == EMPTY {
		return EMPTY, nil
	}
	prev := rune(input[ZERO])
	if unicode.IsDigit(prev) {
		return EMPTY, ErrInvalidString
	}
	var isNumber int
	var builder strings.Builder
	isShielding := false
	for _, val := range input {
		if isShielding {
			if !(val == '\\' || unicode.IsDigit(val)) {
				return EMPTY, ErrInvalidString
			}
			prev = val
			builder.WriteRune(val)
			isShielding = false
			continue
		}
		if val == '\\' {
			isShielding = true
			continue
		}
		err := processCharacter(val, &builder, &prev, &isNumber)
		if err != nil {
			return EMPTY, err
		}
	}

	if isShielding {
		return EMPTY, ErrInvalidString
	}
	return builder.String(), nil
}

func processCharacter(val rune, builder *strings.Builder, prev *rune, isNumber *int) error {
	if unicode.IsDigit(val) {
		return processDigit(val, builder, prev, isNumber)
	}
	processNonDigit(val, builder, prev, isNumber)
	return nil
}

func processDigit(val rune, builder *strings.Builder, prev *rune, isNumber *int) error {
	*isNumber++
	if *isNumber >= TWO {
		return ErrInvalidString
	}
	cnt, err := strconv.Atoi(string(val))
	if err != nil {
		return ErrInvalidString
	}
	if cnt > ZERO {
		builder.WriteString(strings.Repeat(string(*prev), cnt-UNIT))
	} else {
		removeLetter(builder, *prev)
	}
	return nil
}

func processNonDigit(val rune, builder *strings.Builder, prev *rune, isNumber *int) {
	builder.WriteRune(val)
	*prev = val
	*isNumber = ZERO
}

func removeLetter(sb *strings.Builder, char rune) {
	var newBuilder strings.Builder
	removed := false
	for _, val := range sb.String() {
		if val == char && !removed {
			removed = true
			continue
		}
		newBuilder.WriteRune(val)
	}
	sb.Reset()
	sb.WriteString(newBuilder.String())
}
