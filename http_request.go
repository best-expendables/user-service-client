package userclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	authenticatePath     = "/authenticate"
	mePath               = "/me"
	platformsPath        = "/users/%s/platforms"
	logoutPath           = "/logout"
	getPath              = "/users/%s"
	getAllPath           = "/users?per_page=10000"
	getRevokedTokensPath = "/revoked-tokens"
)

type HttpRequestBuilder interface {
	BuildLoginRequest(username string, password string) (*http.Request, error)
	BuildMeRequest(token string) (*http.Request, error)
	BuildPlatformsRequest(token, userID string) (*http.Request, error)
	BuildLogoutRequest(token string) (*http.Request, error)
	BuildGetRequest(token, userID string) (*http.Request, error)
	BuildGetAllRequest(token string) (*http.Request, error)
	BuildRevokedTokensRequest(token string) (*http.Request, error)
}

type HttpRequestBuilderImpl struct {
	BaseURL string
}

func (rb *HttpRequestBuilderImpl) BuildLoginRequest(username string, password string) (*http.Request, error) {
	data := make(map[string]string)
	data["username"] = username
	data["password"] = password

	dataJson, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return rb.build(http.MethodPost, authenticatePath, dataJson)
}

func (rb *HttpRequestBuilderImpl) BuildMeRequest(token string) (*http.Request, error) {
	return rb.buildWithAuth(http.MethodGet, mePath, nil, token)
}

func (rb *HttpRequestBuilderImpl) BuildPlatformsRequest(token, userID string) (*http.Request, error) {
	path := fmt.Sprintf(platformsPath, userID)
	return rb.buildWithAuth(http.MethodGet, path, nil, token)
}

func (rb *HttpRequestBuilderImpl) BuildLogoutRequest(token string) (*http.Request, error) {
	return rb.buildWithAuth(http.MethodPost, logoutPath, nil, token)
}

func (rb *HttpRequestBuilderImpl) BuildGetRequest(token, userID string) (*http.Request, error) {
	path := fmt.Sprintf(getPath, userID)
	return rb.buildWithAuth(http.MethodGet, path, nil, token)
}

func (rb *HttpRequestBuilderImpl) BuildGetAllRequest(token string) (*http.Request, error) {
	return rb.buildWithAuth(http.MethodGet, getAllPath, nil, token)
}

func (rb *HttpRequestBuilderImpl) BuildRevokedTokensRequest(token string) (*http.Request, error) {
	return rb.buildWithAuth(http.MethodGet, getRevokedTokensPath, nil, token)
}

func (rb *HttpRequestBuilderImpl) build(method string, path string, data []byte) (*http.Request, error) {
	return http.NewRequest(
		method,
		fmt.Sprintf("%s%s", rb.BaseURL, path),
		bytes.NewBuffer(data),
	)
}

func (rb *HttpRequestBuilderImpl) buildWithAuth(method string, path string, data []byte, token string) (*http.Request, error) {
	req, err := rb.build(method, path, data)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	return req, nil
}
