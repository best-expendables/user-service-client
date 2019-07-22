package userclient

import "sync"

type TokenHolder interface {
	GetToken(Client) (string, error)
	Invalidate()
}

// Use this one with token from auth middleware
// this holder doesn't retry to obtain new token
type StaticTokenHolder struct {
	token string
}

func NewStaticTokenHolder(token string) *StaticTokenHolder {
	return &StaticTokenHolder{token}
}

func (h *StaticTokenHolder) GetToken(Client) (string, error) {
	return h.token, nil
}

func (h *StaticTokenHolder) Invalidate() {
	h.token = ""
}

// Use this one for your consumers and other backend tasks
// this holder will automatically request a new token if the old one expired
type InMemoryTokenHolder struct {
	sync.Mutex

	Username string
	Password string

	token string
}

func NewInMemoryTokenHolder(username, password string) *InMemoryTokenHolder {
	return &InMemoryTokenHolder{
		Username: username,
		Password: password,
	}
}

func (h *InMemoryTokenHolder) GetToken(client Client) (string, error) {
	if h.token != "" {
		return h.token, nil
	}
	h.Lock()
	defer h.Unlock()
	// another goroutine could update token
	if h.token != "" {
		return h.token, nil
	}
	token, err := client.Authenticate(h.Username, h.Password)
	if token != "" {
		h.token = token
	}
	return token, err
}

func (h *InMemoryTokenHolder) Invalidate() {
	h.Lock()
	defer h.Unlock()
	h.token = ""
}
