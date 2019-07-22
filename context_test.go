package userclient

import (
	"context"
	"testing"
)

func TestGetUserContext(t *testing.T) {
	ctx := ContextWithUser(context.Background(), &User{Id: "1"})
	user := GetCurrentUserFromContext(ctx)
	if user == nil {
		t.Error("User can not be nil")
	}
}

func TestGetTokenFromContext(t *testing.T) {
	expected := "TOKEN"
	ctx := ContextWithToken(context.Background(), expected)
	actual := GetTokenFromContext(ctx)
	if actual != expected {
		t.Errorf("Token not exists or not equals. Expected '%s', actual '%s'", expected, actual)
	}
}
