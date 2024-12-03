package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{Status: StatusOk}
}

func Error(msg string) Response {
	return Response{Status: StatusError, Error: msg}
}

func ValidationError(errs validator.ValidationErrors) Response {
	errMsgs := make([]string, len(errs))

	for idx, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs[idx] = fmt.Sprintf("field %s is a required field", err.Field())
		case "url":
			errMsgs[idx] = fmt.Sprintf("field %s is not a valid URL", err.Field())
		default:
			errMsgs[idx] = fmt.Sprintf("field %s is not valid", err.Field())
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
