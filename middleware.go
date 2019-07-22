package userclient

import (
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const defaultMaxAttempt = 3

const defaultWaitTime = 1000

type Middleware struct {
	userServiceClient Client
	config            RetryConfig
	logger            logger
}

type RetryConfig struct {
	MaxAttempt int
	WaitTime   int
}

var DefaultRetryConfig = RetryConfig{
	MaxAttempt: defaultMaxAttempt,
	WaitTime:   defaultWaitTime,
}

type logger interface {
	Error(...interface{})
}

func NewMiddleware(userServiceClient Client, config RetryConfig) *Middleware {
	return &Middleware{
		userServiceClient: userServiceClient,
		config:            config,
		logger:            logrus.StandardLogger(),
	}
}

func (m *Middleware) SetLogger(l logger) {
	m.logger = l
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		// fallback to get token from QUERY URL if it is not available in HEADER
		if len(tokenString) == 0 {
			tokenString = r.URL.Query().Get("token")
		}

		var user *User
		var err error
		for i := 0; i < m.config.MaxAttempt; i++ {
			user, err = m.userServiceClient.Me(tokenString)
			if err != ErrServiceUnavailable {
				break
			}

			time.Sleep(time.Millisecond * time.Duration(m.config.WaitTime))
		}

		if err != nil {
			m.logger.Error(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if user == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := ContextWithUser(r.Context(), user)
		ctx = ContextWithToken(ctx, tokenString)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CheckRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetCurrentUserFromContext(r.Context())
			if user == nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			if user.HasRole(roles...) == false {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			ctx := ContextWithUser(r.Context(), user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
