package userclient

import (
	"net/http"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		StatusCode int
		Expect     error
	}{
		{401, ErrUnauthorized},
		{404, ErrNotFound},
		{403, ErrServiceUnavailable},
		{500, ErrServiceUnavailable},
	}

	for _, c := range cases {
		result := parseToError(c.StatusCode)
		if result != c.Expect {
			t.Errorf("Expect %v - Got %v", c.Expect, result)
		}
	}
}

func TestParseWithUnknownStatus(t *testing.T) {
	cases := []struct {
		StatusCode int
		Expect     string
	}{
		{510, http.StatusText(510)},
	}

	for _, c := range cases {
		result := parseToError(c.StatusCode)
		if result.Error() != c.Expect {
			t.Errorf("Expect %v - Got %v", c.Expect, result)
		}
	}
}
