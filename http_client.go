package userclient

import (
	"encoding/json"
	"net/http"
	"time"
)

type HttpClient struct {
	requestBuilder HttpRequestBuilder
	requestClient  *http.Client
}

const DEFAULT_TIME_OUT = 10

func NewDefault(baseURL string) *HttpClient {
	return &HttpClient{
		requestBuilder: &HttpRequestBuilderImpl{BaseURL: baseURL},
		requestClient:  &http.Client{Timeout: time.Second * DEFAULT_TIME_OUT},
	}
}

func New(rb HttpRequestBuilder, hc *http.Client) *HttpClient {
	return &HttpClient{
		requestBuilder: rb,
		requestClient:  hc,
	}
}

func (c *HttpClient) Authenticate(username string, password string) (string, error) {
	req, err := c.requestBuilder.BuildLoginRequest(username, password)
	if err != nil {
		return "", err
	}

	resp, err := c.makeRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if val, ok := data["token"]; ok {
		return val.(string), nil
	}
	return "", ErrMissingToken
}

func (c *HttpClient) Me(token string) (*User, error) {
	var response struct {
		Data User
	}

	req, err := c.requestBuilder.BuildMeRequest(token)
	if err != nil {
		return nil, err
	}

	resp, err := c.makeRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	c.addPlatformsToUser(token, &response.Data)

	return &response.Data, nil
}

func (c *HttpClient) Logout(token string) error {
	req, err := c.requestBuilder.BuildLogoutRequest(token)
	if err != nil {
		return err
	}

	resp, err := c.makeRequest(req)
	if err == nil {
		resp.Body.Close()
	}
	return err
}

func (c *HttpClient) FindById(token, userId string) (*User, error) {
	var response struct {
		Data User
	}

	req, err := c.requestBuilder.BuildGetRequest(token, userId)
	if err != nil {
		return nil, err
	}

	resp, err := c.makeRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	c.addPlatformsToUser(token, &response.Data)

	return &response.Data, nil
}

func (c *HttpClient) FindAll(token string) ([]*User, error) {
	req, err := c.requestBuilder.BuildGetAllRequest(token)
	if err != nil {
		return nil, err
	}

	resp, err := c.makeRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	usersResponse := struct {
		Data []*User `json:"data"`
	}{Data: make([]*User, 0)}
	if err := json.NewDecoder(resp.Body).Decode(&usersResponse); err != nil {
		return nil, err
	}

	return usersResponse.Data, nil
}

func (c *HttpClient) RevokedTokens(token string) ([]RevokedToken, error) {
	req, err := c.requestBuilder.BuildRevokedTokensRequest(token)
	if err != nil {
		return nil, err
	}
	resp, err := c.makeRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tokensResponse := struct {
		Data []RevokedToken
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&tokensResponse); err != nil {
		return nil, err
	}

	return tokensResponse.Data, nil
}

func (c *HttpClient) makeRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.requestClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, parseToError(resp.StatusCode)
	}

	return resp, nil
}

func (c *HttpClient) addPlatformsToUser(token string, user *User) error {
	req, err := c.requestBuilder.BuildPlatformsRequest(token, user.Id)
	if err != nil {
		return err
	}
	resp, err := c.makeRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	platforms := struct {
		Data []Platform `json:"data"`
	}{Data: make([]Platform, 0)}
	if err := json.NewDecoder(resp.Body).Decode(&platforms); err != nil {
		return err
	}

	user.PlatformNames = make([]string, len(platforms.Data), len(platforms.Data))
	for i, p := range platforms.Data {
		user.PlatformNames[i] = p.Name
	}

	return nil
}
