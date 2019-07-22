package userclient

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var authCases = []struct {
	ClientMock   UserClientMock
	ExpectCode   int
	ExpectOutput string
}{
	{
		UserClientMock{
			MeMock: func(token string) (*User, error) {
				return nil, errors.New("not found user")
			},
		},
		401,
		"not found user",
	},
	{
		UserClientMock{
			MeMock: func(token string) (*User, error) {
				return nil, nil
			},
		},
		401,
		"Unauthorized",
	},
	{
		UserClientMock{
			MeMock: func(token string) (*User, error) {
				return nil, ErrServiceUnavailable
			},
		},
		401,
		"service is unavailable",
	},
	{
		UserClientMock{
			MeMock: func(token string) (*User, error) {
				return &User{}, nil
			},
		},
		200,
		"ok",
	},
}

var aclCases = []struct {
	ClientMock   UserClientMock
	ExpectCode   int
	ExpectOutput string
}{
	{
		UserClientMock{
			MeMock: func(token string) (*User, error) {
				return &User{Roles: []string{"admin"}}, nil
			},
		},
		403,
		"Forbidden",
	},
	{
		UserClientMock{
			MeMock: func(token string) (*User, error) {
				return &User{Roles: []string{"pm", "admin", "tester"}}, nil
			},
		},
		200,
		"ok",
	},
}

var req *http.Request
var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("ok"))
	req = r
})

func TestAuth(t *testing.T) {
	var makeRequest = func(userCl Client, token string) (int, string) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://localhost/test", nil)

		r.Header.Set("Authorization", token)

		middleware := NewMiddleware(userCl, DefaultRetryConfig)
		middleware.Auth(testHandler).ServeHTTP(w, r)

		actual, _ := ioutil.ReadAll(w.Body)
		return w.Result().StatusCode, strings.TrimSpace(string(actual))
	}

	for _, c := range authCases {
		respStatusCode, respText := makeRequest(&c.ClientMock, "this is token string")
		if respStatusCode != c.ExpectCode {
			t.Errorf("Wrong response code. Expect %v - Got %v", c.ExpectCode, respStatusCode)
		}

		if respText != c.ExpectOutput {
			t.Errorf("Wrong response message. Expect %v - Got %v", c.ExpectOutput, respText)
		}
	}
}

func TestAcl(t *testing.T) {
	var makeRequest = func(userCl Client, token string, roles ...string) (int, string) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://localhost/test", nil)

		r.Header.Set("Authorization", token)

		middleware := NewMiddleware(userCl, DefaultRetryConfig)
		middleware.Auth(CheckRoles(roles...)(testHandler)).ServeHTTP(w, r)

		actual, _ := ioutil.ReadAll(w.Body)
		return w.Result().StatusCode, strings.TrimSpace(string(actual))
	}

	for _, c := range aclCases {
		respStatusCode, respText := makeRequest(&c.ClientMock, "this is token string", "tester", "pm")
		if respStatusCode != c.ExpectCode {
			t.Errorf("Wrong response code. Expect %v - Got %v", c.ExpectCode, respStatusCode)
		}

		if respText != c.ExpectOutput {
			t.Errorf("Wrong response message. Expect %v - Got %v", c.ExpectOutput, respText)
		}
	}
}

func TestAclWithoutAuth(t *testing.T) {
	var makeRequest = func(token string, roles ...string) (int, string) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://localhost/test", nil)

		r.Header.Set("Authorization", token)

		CheckRoles(roles...)(testHandler).ServeHTTP(w, r)

		actual, _ := ioutil.ReadAll(w.Body)
		return w.Result().StatusCode, strings.TrimSpace(string(actual))
	}

	respStatusCode, respText := makeRequest("this is token string", "tester", "pm")
	if respStatusCode != http.StatusUnauthorized {
		t.Errorf("Wrong response code. Expect %v - Got %v", http.StatusUnauthorized, respStatusCode)
	}
	if respText != http.StatusText(http.StatusUnauthorized) {
		t.Errorf("Wrong response message. Expect %v - Got %v", http.StatusText(http.StatusUnauthorized), respText)
	}
}

func TestAuthWithTokenInURL(t *testing.T) {
	var makeRequest = func(userCl Client, token string) (int, string) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://localhost/test?token=THIS_IS_A_TOKEN", nil)

		middleware := NewMiddleware(userCl, DefaultRetryConfig)
		middleware.Auth(testHandler).ServeHTTP(w, r)

		actual, _ := ioutil.ReadAll(w.Body)
		return w.Result().StatusCode, strings.TrimSpace(string(actual))
	}

	for _, c := range authCases {
		respStatusCode, respText := makeRequest(&c.ClientMock, "this is token string")
		if respStatusCode != c.ExpectCode {
			t.Errorf("Wrong response code. Expect %v - Got %v", c.ExpectOutput, respStatusCode)
		}

		if respText != c.ExpectOutput {
			t.Errorf("Wrong response message. Expect %v - Got %v", c.ExpectOutput, respText)
		}
	}
}

func TestGetCurrentUserFromContext(t *testing.T) {
	var makeRequest = func(userCl Client, token string) *http.Request {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://localhost/test", nil)

		r.Header.Set("Authorization", token)

		middleware := NewMiddleware(userCl, DefaultRetryConfig)
		middleware.Auth(testHandler).ServeHTTP(w, r)

		return req
	}

	clientMock := UserClientMock{
		MeMock: func(token string) (*User, error) {
			return &User{
				Username: "Test",
				Email:    "test@lazada.com",
			}, nil
		},
	}

	r := makeRequest(&clientMock, "this is token string")
	user := GetCurrentUserFromContext(r.Context())

	if user == nil {
		t.Error("Can not get current user from context")
	}

	if user.Username != "Test" {
		t.Errorf("Wrong user name. Expect %v - Got %v", "Test", user.Username)
	}

	if user.Email != "test@lazada.com" {
		t.Errorf("Wrong user name. Expect %v - Got %v", "test@lazada.com", user.Email)
	}
}
