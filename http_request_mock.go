package userclient

import (
	"net/http"
)

type HttpRequestBuilderMock struct {
	RequestURL                    string
	BuildLoginRequestMock         func(username string, password string) (*http.Request, error)
	BuildMeRequestMock            func(token string) (*http.Request, error)
	BuildPlatformsRequestMock     func(token, userID string) (*http.Request, error)
	BuildLogoutRequestMock        func(token string) (*http.Request, error)
	BuildGetRequestMock           func(token, userID string) (*http.Request, error)
	BuildGetAllRequestMock        func(token string) (*http.Request, error)
	BuildRevokedTokensRequestMock func(token string) (*http.Request, error)
}

func (rb *HttpRequestBuilderMock) BuildLoginRequest(username string, password string) (*http.Request, error) {
	return rb.BuildLoginRequestMock(username, password)
}

func (rb *HttpRequestBuilderMock) BuildMeRequest(token string) (*http.Request, error) {
	return rb.BuildMeRequestMock(token)
}

func (rb *HttpRequestBuilderMock) BuildPlatformsRequest(token, userID string) (*http.Request, error) {
	return rb.BuildPlatformsRequestMock(token, userID)
}

func (rb *HttpRequestBuilderMock) BuildLogoutRequest(token string) (*http.Request, error) {
	return rb.BuildLogoutRequestMock(token)
}

func (rb *HttpRequestBuilderMock) BuildGetRequest(token, userID string) (*http.Request, error) {
	return rb.BuildGetRequestMock(token, userID)
}
func (rb *HttpRequestBuilderMock) BuildGetAllRequest(token string) (*http.Request, error) {
	return rb.BuildGetAllRequestMock(token)
}

func (rb *HttpRequestBuilderMock) BuildRevokedTokensRequest(token string) (*http.Request, error) {
	return rb.BuildRevokedTokensRequestMock(token)
}

type CacheMock struct {
	GetFn    func(key string, obj interface{}) error
	SetFn    func(key string, obj interface{}) error
	DeleteFn func(key string) error
}

func (c *CacheMock) Get(key string, obj interface{}) error {
	return c.GetFn(key, obj)
}
func (c *CacheMock) Set(key string, obj interface{}) error {
	return c.SetFn(key, obj)
}
func (c *CacheMock) Delete(key string) error {
	return c.DeleteFn(key)
}
