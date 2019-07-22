package userclient

import (
	"fmt"
	"net/http"
	"testing"
)

var builder HttpRequestBuilderImpl

func TestRequestBuilder_BuildLoginRequest(t *testing.T) {
	req, err := builder.BuildLoginRequest("admin", "123")

	if err != nil {
		t.Error("Has error when creating user login request!")
	}

	if req == nil {
		t.Error("Login request is nil")
	}

	if req.Method != http.MethodPost {
		t.Error("Method for login request is not POST")
	}
}

func TestRequestBuilder_BuildMeRequest(t *testing.T) {
	token := "this is a very long token"

	req, err := builder.BuildMeRequest(token)

	if err != nil {
		t.Error("Has error when creatinge me request!")
	}

	if req == nil {
		t.Error("Me request is nil")
	}

	if req.Method != http.MethodGet {
		t.Error("Method for me request is not GET")
	}

	if req.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", token) {
		t.Error("Wrong header for me request")
	}
}

func TestRequestBuilder_BuildPlatformsRequest(t *testing.T) {
	token := "this is a very long token"
	userID := "correct-id"

	req, err := builder.BuildPlatformsRequest(token, userID)

	if err != nil {
		t.Error("Has error when creating platforms request!")
	}

	if req == nil {
		t.Error("Platform request is nil")
	}

	if req.Method != http.MethodGet {
		t.Error("Method for platform request is not GET")
	}

	if req.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", token) {
		t.Error("Wrong header for platform request")
	}
}

func TestRequestBuilder_BuildLogoutRequest(t *testing.T) {
	token := "this is a very long token"

	req, err := builder.BuildLogoutRequest(token)

	if err != nil {
		t.Error("Has error when creatinge logout request!")
	}

	if req == nil {
		t.Error("Logout request is nill")
	}

	if req.Method != http.MethodPost {
		t.Error("Method for logout request is not POST")
	}

	if req.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", token) {
		t.Error("Wrong header for logout request")
	}
}

func TestRequestBuilder_BuildGetRequest(t *testing.T) {
	token := "this is a very long token"

	req, err := builder.BuildGetRequest(token, "correct-user-id")

	if err != nil {
		t.Error("Has error when creating get request!")
	}

	if req == nil {
		t.Error("Get request is nil")
	}

	if req.Method != http.MethodGet {
		t.Error("Method for get request is not GET")
	}

	if req.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", token) {
		t.Error("Wrong header for get request")
	}
}

func TestRequestBuilder_BuildGetAllRequest(t *testing.T) {
	token := "this is a very long token"

	req, err := builder.BuildGetAllRequest(token)

	if err != nil {
		t.Error("Has error when creating get all request!")
	}

	if req == nil {
		t.Error("Get all request is nil")
	}

	if req.Method != http.MethodGet {
		t.Error("Method for get all request is not GET")
	}

	if req.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", token) {
		t.Error("Wrong header for get all request")
	}
}

func TestRequestBuilder_BuildRevokedTokensRequest(t *testing.T) {
	token := "this is a very long token"

	req, err := builder.BuildRevokedTokensRequest(token)

	if err != nil {
		t.Error("Has error when creating Revoked tokens request!")
	}

	if req == nil {
		t.Error("Revoked tokens request is nil")
	}

	if req.Method != http.MethodGet {
		t.Error("Method for Revoked tokens request is not GET")
	}

	if req.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", token) {
		t.Error("Wrong header for Revoked tokens request")
	}
}
