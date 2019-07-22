package userclient

import (
	"errors"
	"testing"
)

func Test_ClientWithTokenHolder_Correct(t *testing.T) {
	user := &User{}

	c := &UserClientMock{
		MeMock: func(token string) (*User, error) {
			if token != "token" {
				return nil, errors.New("wrong token")
			}
			return user, nil
		},
		FindByIdMock: func(token, userId string) (*User, error) {
			if token != "token" {
				return nil, errors.New("wrong token")
			}
			return user, nil
		},
		FindAllMock: func(token string) ([]*User, error) {
			if token != "token" {
				return nil, errors.New("wrong token")
			}
			return []*User{user}, nil
		},
	}
	holder := &tokenHolderMock{
		GetTokenFunc: func(_ Client) (string, error) {
			return "token", nil
		},
	}
	client := NewClientWithTokenHolder(c, holder)

	u, err := client.Me()
	if err != nil {
		t.Error("Shouldn't return error")
	}
	if u != user {
		t.Error("Me should return user")
	}

	u, err = client.FindById("id")
	if err != nil {
		t.Error("Shouldn't return error")
	}
	if u != user {
		t.Error("FindById should return user")
	}

	users, err := client.FindAll()
	if err != nil {
		t.Error("Shouldn't return error")
	}
	if len(users) != 1 {
		t.Error("should return slice with 1 user")
	}
	if users[0] != user {
		t.Error("FindAll should return user")
	}
}

func Test_ClientWithTokenHolder_Invalidate(t *testing.T) {
	invalidateCalled := 0

	c := &UserClientMock{
		MeMock: func(token string) (*User, error) {
			return nil, ErrUnauthorized
		},
		FindByIdMock: func(token, userId string) (*User, error) {
			return nil, ErrUnauthorized
		},
		FindAllMock: func(token string) ([]*User, error) {
			return nil, ErrUnauthorized
		},
	}
	holder := &tokenHolderMock{
		GetTokenFunc: func(_ Client) (string, error) {
			return "token", nil
		},
		InvalidateFunc: func() {
			invalidateCalled++
		},
	}
	client := NewClientWithTokenHolder(c, holder)

	invalidateCalled = 0
	u, err := client.Me()
	if err != ErrUnauthorized {
		t.Error("Should return ErrUnauthorized error")
	}
	if u != nil {
		t.Error("Me shouldn't return user")
	}
	if invalidateCalled != 1 {
		t.Error("invalidate should be called")
	}

	invalidateCalled = 0
	u, err = client.FindById("id")
	if err != ErrUnauthorized {
		t.Error("Should return ErrUnauthorized error")
	}
	if u != nil {
		t.Error("FindById shouldn't return user")
	}
	if invalidateCalled != 1 {
		t.Error("invalidate should be called")
	}

	invalidateCalled = 0
	users, err := client.FindAll()
	if err != ErrUnauthorized {
		t.Error("Should return ErrUnauthorized error")
	}
	if users != nil {
		t.Error("should return nil")
	}
	if invalidateCalled != 1 {
		t.Error("invalidate should be called")
	}
}

func Test_ClientWithTokenHolder_Pass_Errors(t *testing.T) {
	clientError := errors.New("client error")
	c := &UserClientMock{
		MeMock: func(token string) (*User, error) {
			return nil, clientError
		},
		FindByIdMock: func(token, userId string) (*User, error) {
			return nil, clientError
		},
		FindAllMock: func(token string) ([]*User, error) {
			return nil, clientError
		},
	}
	holder := &tokenHolderMock{
		GetTokenFunc: func(_ Client) (string, error) {
			return "token", nil
		},
	}
	client := NewClientWithTokenHolder(c, holder)

	u, err := client.Me()
	if err != clientError {
		t.Error("Should return clientError error")
	}
	if u != nil {
		t.Error("Me shouldn't return user")
	}

	u, err = client.FindById("id")
	if err != clientError {
		t.Error("Should return clientError error")
	}
	if u != nil {
		t.Error("FindById shouldn't return user")
	}

	users, err := client.FindAll()
	if err != clientError {
		t.Error("Should return clientError error")
	}
	if users != nil {
		t.Error("should return nil")
	}
}

type tokenHolderMock struct {
	GetTokenFunc   func(c Client) (string, error)
	InvalidateFunc func()
}

func (t *tokenHolderMock) GetToken(c Client) (string, error) {
	return t.GetTokenFunc(c)
}

func (t *tokenHolderMock) Invalidate() {
	t.InvalidateFunc()
}
