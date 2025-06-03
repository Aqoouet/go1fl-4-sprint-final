package common

import (
	"errors"
)

var (
	ErrSliceLen           = errors.New("wrong slice length")
	ErrParamLimitExceeded = errors.New("parameter exceeded allowable value")
	ErrParseDuration      = errors.New("duration parsing error")
	ErrParseInt           = errors.New("int parsing error")
	ErrEmptyString        = errors.New("empty string")
)
