// Package helpers provides api errors helpers for the CloudAvenue Terraform Provider.
package helpers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	apiclient "github.com/orange-cloudavenue/cloudavenue-sdk-go"
)

// APIError is an error returned by the CloudAvenue API.
type APIError struct {
	lastError  error
	statusCode int
	apiclient.ApiError
}

// Error returns the error message.
func (e *APIError) Error() string {
	return e.lastError.Error()
}

// Unwrap returns the underlying error.
func (e *APIError) Unwrap() error {
	return e.lastError
}

// GetSummary returns a summary of the error.
func (e *APIError) GetSummary() string {
	return fmt.Sprintf("%s (HTTP Code => %d)", e.Reason, e.GetStatusCode())
}

// GetDetail returns a detailed description of the error.
func (e *APIError) GetDetail() string {
	return e.Message
}

// GetStatusCode returns the HTTP status code.
func (e *APIError) GetStatusCode() int {
	return e.statusCode
}

// GetTerraformDiagnostic returns a Terraform Diagnostic for the error.
func (e *APIError) GetTerraformDiagnostic() diag.Diagnostic {
	var summary, detail string
	if e.Reason != "" {
		summary = e.GetSummary()
	} else {
		summary = fmt.Sprintf("HTTP response error: %d", e.GetStatusCode())
	}
	if e.Message != "" {
		detail = e.GetDetail()
	} else {
		detail = e.lastError.Error()
	}

	if e.IsNotFound() {
		return diag.NewWarningDiagnostic(summary, detail)
	}
	return diag.NewErrorDiagnostic(summary, detail)
}

// IsNotFound returns true if the error is a 404 Not Found error.
func (e *APIError) IsNotFound() bool {
	return e.GetStatusCode() == http.StatusNotFound
}

// CheckAPIError checks the HTTP response for errors and returns an APIError
// if the response code is >= 400.
// If the response code is < 400, nil is returned.
func CheckAPIError(err error, httpR *http.Response) *APIError {
	if err == nil {
		return nil
	}

	if httpR != nil && httpR.StatusCode >= http.StatusBadRequest {
		var apiErr *apiclient.GenericSwaggerError
		if errors.As(err, &apiErr) {
			x := &APIError{
				lastError:  err,
				statusCode: httpR.StatusCode,
			}

			if apiErr.Model() != nil {
				if m, ok := apiErr.Model().(apiclient.ApiError); ok {
					x.ApiError = m
				}
			}

			return x
		}
	} else {
		return nil
	}

	return nil
}
