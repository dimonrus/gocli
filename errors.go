package main

import (
	"fmt"
	"runtime/debug"
)

// Application Error
type ApplicationError struct {
	message string
	code    int
	stack   []byte
	details []ApplicationErrorDetail
}

// Interface error method
func (ae ApplicationError) Error() string {
	return ae.message
}

// Interface app error get code
func (ae ApplicationError) GetCode() int {
	return ae.code
}

// Interface app error get stack
func (ae ApplicationError) GetStack() []byte {
	return ae.stack
}

// Interface app error get details
func (ae ApplicationError) GetDetails() []ApplicationErrorDetail {
	return ae.details
}

// Interface app error add detail
func (ae ApplicationError) AddDetail(message string, code int) ApplicationErrorInterface {
	ae.details = append(ae.details, ApplicationErrorDetail{Message:message, Code:code})
	return &ae
}

// Detailed error
type ApplicationErrorDetail struct {
	Code    int
	Message string
}

// Common Error Interface
type ApplicationErrorInterface interface {
	Error() string
	GetCode() int
	GetStack() []byte
	GetDetails() []ApplicationErrorDetail
	AddDetail(message string, code int) ApplicationErrorInterface
}

// New error
func NewError(message string, code int) ApplicationErrorInterface {
	err := ApplicationError{
		message: message,
		code:    code,
		stack:   debug.Stack(),
	}

	return &err
}

// New error
func NewErrorF(message string, code int, args ... interface{}) ApplicationErrorInterface {
	err := ApplicationError{
		message: fmt.Sprintf(message, args),
		code:    code,
		stack:   debug.Stack(),
	}

	return &err
}
