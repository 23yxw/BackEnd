package utils

import "errors"

var (
	ErrAlgorithmFailed  error = errors.New("failed to analyse audio")
	ErrNotPermitted     error = errors.New("operation not permitted")
	ErrRemoteCallFailed error = errors.New("remote algorithm call failed")
	ErrUnknown          error = errors.New("unknown error")
)

type ErrCode int

var ErrCodeMap map[error]ErrCode = map[error]ErrCode{
	ErrAlgorithmFailed:  10,
	ErrNotPermitted:     2,
	ErrRemoteCallFailed: 3,
}
