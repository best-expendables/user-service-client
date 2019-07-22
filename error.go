package userclient

import (
	"errors"
	"net/http"
)

var (
	ErrServiceUnavailable = errors.New("service is unavailable")

	ErrUnauthorized = errors.New("unauthorized")

	ErrNotFound = errors.New("not found user")

	ErrMissingToken = errors.New("missing token in response")
)

var serviceUnavailableCodes = []int{
	http.StatusForbidden,
	http.StatusInternalServerError,
	http.StatusBadGateway,
	http.StatusServiceUnavailable,
	http.StatusGatewayTimeout,
}

func parseToError(statusCode int) error {
	if statusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if statusCode == http.StatusNotFound {
		return ErrNotFound
	}

	for _, c := range serviceUnavailableCodes {
		if statusCode == c {
			return ErrServiceUnavailable
		}
	}

	return errors.New(http.StatusText(statusCode))
}
