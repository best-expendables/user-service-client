package userclient

import (
	"errors"
	"testing"
)

func Test_StaticTokenHolder_GetToken(t *testing.T) {
	client := &UserClientMock{}
	token := "test-token"

	holder := NewStaticTokenHolder(token)
	recievedToken, err := holder.GetToken(client)
	if err != nil {
		t.Errorf("error '%s' returned", err)
	}
	if recievedToken != token {
		t.Errorf("token '%s' isn't equal to expected token '%s'", recievedToken, token)
	}
}

func Test_StaticTokenHolder_Invalidate(t *testing.T) {
	client := &UserClientMock{}
	token := "test-token"

	holder := NewStaticTokenHolder(token)

	holder.Invalidate()

	recievedToken, err := holder.GetToken(client)
	if err != nil {
		t.Errorf("error '%s' returned", err)
	}
	if recievedToken != "" {
		t.Errorf("token should be reset, but it is '%s'", recievedToken)
	}
}

func Test_InMemoryTokenHolder_GetToken(t *testing.T) {
	token := "correct-token"
	clientCalls := 0
	clientError := errors.New("client error")
	client := &UserClientMock{
		AuthenticateMock: func(username string, password string) (string, error) {
			if username == "incorrect" {
				return "", clientError
			}
			clientCalls++
			return token, nil
		},
	}

	holder := NewInMemoryTokenHolder("correct", "123")
	recievedToken, err := holder.GetToken(client)
	if err != nil {
		t.Errorf("error '%s' returned", err)
	}
	if recievedToken != token {
		t.Errorf("token '%s' isn't equal to expected token '%s'", recievedToken, token)
	}
	if clientCalls != 1 {
		t.Errorf("client should be called once, but it called %d", clientCalls)
	}

	// should use cache
	recievedToken, err = holder.GetToken(client)
	if err != nil {
		t.Errorf("error '%s' returned", err)
	}
	if recievedToken != token {
		t.Errorf("token '%s' isn't equal to expected token '%s'", recievedToken, token)
	}
	if clientCalls != 1 {
		t.Errorf("client should be called once, but it called %d", clientCalls)
	}

	holder = NewInMemoryTokenHolder("incorrect", "123")
	recievedToken, err = holder.GetToken(client)
	if err != clientError {
		t.Errorf("error '%s' should be returned, but it is '%s'", clientError, err)
	}
	if recievedToken != "" {
		t.Errorf("token should be empty, but it is '%s'", recievedToken)
	}
}

func Test_InMemoryTokenHolder_Invalidate(t *testing.T) {
	token := "correct-token"
	clientCalls := 0
	clientError := errors.New("client error")
	client := &UserClientMock{
		AuthenticateMock: func(username string, password string) (string, error) {
			if username == "incorrect" {
				return "", clientError
			}
			clientCalls++
			return token, nil
		},
	}

	holder := NewInMemoryTokenHolder("correct", "123")
	recievedToken, err := holder.GetToken(client)

	holder.Invalidate()

	recievedToken, err = holder.GetToken(client)
	if err != nil {
		t.Errorf("error '%s' returned", err)
	}
	if recievedToken != token {
		t.Errorf("token '%s' isn't equal to expected token '%s'", recievedToken, token)
	}
	if clientCalls != 2 {
		t.Errorf("client should be called once, but it called %d", clientCalls)
	}
}
