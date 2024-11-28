package goerr

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ZapError(err error) zapcore.Field {
	if goErr := Unwrap(err); goErr != nil {
		return zap.Object("error", goErr)
	}
	return zap.Error(err)
}
