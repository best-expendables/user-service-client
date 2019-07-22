package userclient

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

var httpClient *HttpClient
var builderMock HttpRequestBuilderMock

func setup() {
	builder = HttpRequestBuilderImpl{}

	builderMock = HttpRequestBuilderMock{
		BuildLoginRequestMock: func(username string, password string) (*http.Request, error) {
			data := make(map[string]string)
			data["username"] = username
			data["password"] = password
			dataJson, _ := json.Marshal(data)

			req, _ := http.NewRequest(http.MethodPost, builderMock.RequestURL, bytes.NewBuffer(dataJson))
			return req, nil
		},
		BuildMeRequestMock: func(token string) (*http.Request, error) {
			req, _ := http.NewRequest(http.MethodGet, builderMock.RequestURL, nil)
			return req, nil
		},
		BuildPlatformsRequestMock: func(token, userID string) (*http.Request, error) {
			url := builderMock.RequestURL + "/platforms"
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			return req, nil
		},
		BuildLogoutRequestMock: func(token string) (*http.Request, error) {
			req, _ := http.NewRequest(http.MethodPost, builderMock.RequestURL, nil)
			return req, nil
		},
		BuildGetRequestMock: func(token, userID string) (*http.Request, error) {
			req, _ := http.NewRequest(http.MethodGet, builderMock.RequestURL, nil)
			return req, nil
		},
		BuildGetAllRequestMock: func(token string) (*http.Request, error) {
			req, _ := http.NewRequest(http.MethodGet, builderMock.RequestURL, nil)
			return req, nil
		},
		BuildRevokedTokensRequestMock: func(token string) (*http.Request, error) {
			req, _ := http.NewRequest(http.MethodGet, builderMock.RequestURL, nil)
			return req, nil
		},
	}

	httpClient = New(&builderMock, &http.Client{})
}

func TestMain(m *testing.M) {
	flag.Parse()

	setup()
	code := m.Run()
	os.Exit(code)
}

func TestUserHttpClient_Authenticate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			enc := json.NewEncoder(w)
			enc.Encode(map[string]string{"token": "this is a jwt token string"})
		}))

	defer ts.Close()

	builderMock.RequestURL = ts.URL
	token, err := httpClient.Authenticate("username", "password")
	if err != nil {
		t.Error("Has error when testing authenticate:", err.Error())
	}
	if token != "this is a jwt token string" {
		t.Error("Invalid token return: ", token)
	}
}

func TestUserHttpClient_AuthenticateHasReturnError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

	defer ts.Close()

	builderMock.RequestURL = ts.URL
	token, err := httpClient.Authenticate("username", "password")
	if token != "" {
		t.Error("Token should be empty if we have error")
	}

	if err != ErrNotFound {
		t.Error("Error messsage should be not found")
	}
}

func TestUserHttpClient_Me(t *testing.T) {
	returnUser := User{
		Id:        "a802918c-4471-46a1-989b-c0cf651a4b2c",
		Username:  "test",
		Email:     "test@gmail.com",
		Active:    true,
		Roles:     []string{"admin", "pm", "tester"},
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	platformsData := struct {
		Data []Platform `json:"data"`
	}{Data: []Platform{
		{Name: "OMS_VN"},
		{Name: "ALIBABA"},
	}}
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			enc := json.NewEncoder(w)
			if r.RequestURI == "/platforms" {
				w.WriteHeader(http.StatusOK)
				enc.Encode(platformsData)
				return
			}
			w.WriteHeader(http.StatusOK)
			enc.Encode(
				map[string]User{
					"data": returnUser,
				},
			)
		}))

	defer ts.Close()

	builderMock.RequestURL = ts.URL
	user, err := httpClient.Me("This is a token string")
	if err != nil {
		t.Error("Has error when testing me request:", err.Error())
	}

	assertUser(t, user, &returnUser)
}

func TestUserHttpClient_Logout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

	defer ts.Close()

	builderMock.RequestURL = ts.URL
	err := httpClient.Logout("This is a token string")
	if err != nil {
		t.Error("Has error when testing me request:", err.Error())
	}
}

func TestUserHttpClient_FindById(t *testing.T) {
	returnUser := User{
		Id:        "a802918c-4471-46a1-989b-c0cf651a4b2c",
		Username:  "test",
		Email:     "test@gmail.com",
		Active:    true,
		Roles:     []string{"admin", "pm", "tester"},
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	platformsData := struct {
		Data []Platform `json:"data"`
	}{Data: []Platform{
		{Name: "OMS_VN"},
		{Name: "ALIBABA"},
	}}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		if r.RequestURI == "/platforms" {
			w.WriteHeader(http.StatusOK)
			enc.Encode(platformsData)
			return
		}
		w.WriteHeader(http.StatusOK)
		enc.Encode(struct {
			Data User `json:"data"`
		}{returnUser})
	}))
	defer ts.Close()

	builderMock.RequestURL = ts.URL
	user, err := httpClient.FindById("This is a token string", "correct-user-id")
	if err != nil {
		t.Error("Has error when testing me request:", err.Error())
	}

	assertUser(t, user, &returnUser)
}

func TestUserHttpClient_RevokedTokens(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"data":[
				{"token": "token1", "expiredAt": "2017-08-14T18:19:03Z"},
				{"token": "token2", "expiredAt": "2017-08-18T10:20:00Z"}
			]}`))
		}))
	defer ts.Close()

	builderMock.RequestURL = ts.URL
	tokens, err := httpClient.RevokedTokens("This is a token string")
	if err != nil {
		t.Error("Has error when testing revoked tokens request:", err.Error())
	}
	if len(tokens) != 2 {
		t.Fatalf("should be 2 token, returned %d", len(tokens))
	}
	if tokens[0].Token != "token1" {
		t.Error("wrong token")
	}
	if tokens[0].ExpiredAt.Format(time.RFC3339) != "2017-08-14T18:19:03Z" {
		t.Error("wrong token expired at")
	}
	if tokens[1].Token != "token2" {
		t.Error("wrong token")
	}
	if tokens[1].ExpiredAt.Format(time.RFC3339) != "2017-08-18T10:20:00Z" {
		t.Error("wrong token expired at")
	}
}

func assertUser(t *testing.T, user, returnUser *User) {
	if user.Id != returnUser.Id {
		t.Errorf("Return user id '%s' is invalid, expected '%s'", user.Id, returnUser.Id)
	}
	if user.Username != returnUser.Username {
		t.Errorf("Return user name '%s' is invalid, expected '%s'", user.Username, returnUser.Username)
	}
	if user.Email != returnUser.Email {
		t.Errorf("Return user email '%s' is invalid, expected '%s'", user.Email, returnUser.Email)
	}
	if !reflect.DeepEqual(user.Roles, returnUser.Roles) {
		t.Errorf("Return user roles '%v' are invalid, expected '%s'", user.Roles, returnUser.Roles)
	}
	if !reflect.DeepEqual(user.PlatformNames, []string{"OMS_VN", "ALIBABA"}) {
		t.Errorf("Return user platfrom names '%v' are invalid, expected '%s'", user.PlatformNames, returnUser.PlatformNames)
	}
}
