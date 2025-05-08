package xmatters

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	// ErrNoContent is a generic 204 Error output used to return appropriate output to the user after a failed DELETE request.
	ErrNoContent = XMattersError{
		Code:    StatusNoContent,
		Message: "A resource was not found in response to a DELETE request.",
		Reason:  "No Content",
	}
	// ErrInavlidCredentials is a generic 401 Error output used to return appropriate output to the user after a failed request due to invalid credentials.
	ErrInavlidCredentials = XMattersError{
		Code:    StatusUnauthorized,
		Message: "Invalid Credentials",
		Reason:  "Unauthorized",
	}
	// ErrNoHostname is a generic Error output used to return appropriate output to the user after a failed request due to missing hostname.
	ErrNoHostname = XMattersError{
		Code:    0,
		Message: "Missing Hostname",
		Reason:  "Bad Request",
	}
	// General error message content
	errUnmarshalError     = "error unmarshalling the JSON response"
	errUnmarshalErrorBody = "error unmarshalling the JSON response error body"
)

// XMattersError is a custom error type with helpful fields.
type XMattersError struct {
	Code    int    `json:"code,omitempty"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Subcode string `json:"subcode,omitempty"`
}

// Error implements the error interface for xMattersError.
func (e XMattersError) Error() string {
	return fmt.Sprintf("xMatters API Error: %d - %s. %s\nSubcode: %s", e.Code, e.Reason, e.Message, e.Subcode)
}

// getFunctionName retrieves the name of the function that called `newUnmarshalError`.
// It uses runtime.Caller to get the program counter and function name.
func getFunctionName() string {
	pc, _, _, ok := runtime.Caller(2) // 2 means two levels up from this function (the function that called `newUnmarshalError`)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	// Extract function name without package
	fullName := fn.Name()
	return filepath.Base(fullName[strings.LastIndex(fullName, ".")+1:])
}

// newUnmarshalError creates a new XMattersError with a generic unmarshal error message.
// It uses the getFunctionName function to include the name of the function that called it.
// This is useful for debugging and understanding where the error occurred.
// The error code is set to 0, indicating a generic error.
func newUnmarshalError() error {
	return XMattersError{
		Code:    0,
		Message: errUnmarshalError,
		Reason:  fmt.Sprintf("Internal Server Error in %s", getFunctionName()),
	}
}

// NewXMattersError is a constructor function to create a new xMattersError instance
func newXMattersError(body []byte) error {
	var xmerr XMattersError
	err := json.Unmarshal(body, &xmerr)
	if err != nil {
		return fmt.Errorf("%s in xMatters Error Construction: %w \n%s", errUnmarshalErrorBody, err, string(body))
	}
	return xmerr
}
