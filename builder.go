package goerr

import (
	"fmt"
)

// Builder keeps a set of key-value pairs and can create a new error and wrap error with the key-value pairs.
type Builder struct {
	values   values
	metadata map[string]string
}

// NewBuilder creates a new Builder
func NewBuilder() *Builder {
	return &Builder{
		values:   make(values),
		metadata: make(map[string]string),
	}
}

// With copies the current Builder and adds a new key-value pair.
func (x *Builder) With(key string, value any) *Builder {
	newVS := &Builder{values: x.values.clone(), metadata: x.cloneMetadata()}
	newVS.values[key] = value
	return newVS
}

// WithMetadata copies the current Builder and adds a response-visible key-value pair.
func (x *Builder) WithMetadata(key, value string) *Builder {
	newVS := &Builder{values: x.values.clone(), metadata: x.cloneMetadata()}
	newVS.metadata[key] = value
	return newVS
}

func (x *Builder) cloneMetadata() map[string]string {
	m := make(map[string]string, len(x.metadata))
	for k, v := range x.metadata {
		m[k] = v
	}
	return m
}

// New creates a new error with message
func (x *Builder) New(format string, args ...any) *Error {
	err := newError()
	err.msg = fmt.Sprintf(format, args...)
	err.values = x.values.clone()
	for k, v := range x.metadata {
		err.metadata[k] = v
	}

	return err
}

// Wrap creates a new Error with caused error and add message.
func (x *Builder) Wrap(cause error, msg ...any) *Error {
	err := newError()
	err.msg = toWrapMessage(msg)
	err.cause = cause
	err.values = x.values.clone()
	for k, v := range x.metadata {
		err.metadata[k] = v
	}
	return err
}
