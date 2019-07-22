package userclient

import (
	"errors"
	"testing"
)

func TestUserCacheClient_Me(t *testing.T) {
	calls := 0
	uncachedClient := &UserClientMock{
		MeMock: func(token string) (*User, error) {
			calls++
			return &User{}, nil
		},
	}
	cache := newCachedMock()

	client := NewCacheClient(uncachedClient, cache)
	_, err := client.Me("token")
	if err != nil {
		t.Error("Cached client returned error")
	}
	if calls != 1 {
		t.Error("Underlying client hasn't been called")
	}
	_, err = client.Me("token")
	if err != nil {
		t.Error("Cached client returned error")
	}
	if calls > 1 {
		t.Error("Underlying client has been called")
	}
}

func TestUserCacheClient_FindById(t *testing.T) {
	calls := 0
	uncachedClient := &UserClientMock{
		FindByIdMock: func(token, id string) (*User, error) {
			calls++
			return &User{}, nil
		},
	}
	cache := newCachedMock()

	client := NewCacheClient(uncachedClient, cache)
	_, err := client.FindById("token", "id")
	if err != nil {
		t.Error("Cached client returned error")
	}
	if calls != 1 {
		t.Error("Underlying client hasn't been called")
	}
	_, err = client.FindById("token", "id")
	if err != nil {
		t.Error("Cached client returned error")
	}
	if calls > 1 {
		t.Error("Underlying client has been called")
	}
}

type cacheMock struct {
	items map[string]interface{}
}

func newCachedMock() *cacheMock {
	return &cacheMock{items: make(map[string]interface{})}
}

func (c *cacheMock) Get(key string, obj interface{}) error {
	_, ok := c.items[key]
	if !ok {
		return errors.New("not found")
	}
	return nil
}

func (c *cacheMock) Set(key string, obj interface{}) error {
	c.items[key] = obj
	return nil
}

func (c *cacheMock) Delete(key string) error {
	return nil
}

func (c *cacheMock) Append(key, obj interface{}) error {
	return nil
}
