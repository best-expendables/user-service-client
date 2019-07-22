package userclient

import (
	"fmt"
)

type CacheClient struct {
	client Client
	cache  Cache
}

func NewCacheClient(client Client, cache Cache) *CacheClient {
	return &CacheClient{
		client: client,
		cache:  cache,
	}
}

func (c *CacheClient) Authenticate(username string, password string) (string, error) {
	return c.client.Authenticate(username, password)
}

func (c *CacheClient) Me(token string) (*User, error) {
	var user = &User{}
	cacheKey := c.cacheKeyMe(token)
	if err := c.cache.Get(cacheKey, user); err == nil {
		return user, nil
	}
	user, err := c.client.Me(token)
	if err != nil {
		return user, err
	}
	c.cache.Set(cacheKey, user)
	return user, nil
}

func (c *CacheClient) Logout(token string) error {
	c.cleanUp(token)
	return c.client.Logout(token)
}

func (c *CacheClient) FindById(token, userId string) (*User, error) {
	var user = &User{}
	cacheKey := c.cacheKeyFindById(token, userId)
	if err := c.cache.Get(cacheKey, user); err == nil {
		return user, nil
	}
	user, err := c.client.FindById(token, userId)
	if err != nil {
		return user, err
	}
	c.cache.Set(cacheKey, user)
	return user, nil
}

func (c *CacheClient) FindAll(token string) ([]*User, error) {
	return c.client.FindAll(token)
}

func (c *CacheClient) RevokedTokens(token string) ([]RevokedToken, error) {
	return c.client.RevokedTokens(token)
}

func (c *CacheClient) cacheKeyMe(token string) string {
	return fmt.Sprintf("user-middleware/%s/me", token)
}

func (c *CacheClient) cacheKeyFindById(token, userId string) string {
	key := fmt.Sprintf("user-middleware/%s/me/%s", token, userId)
	storedKeysKey := c.cacheKeyStoredKeys(token)
	var keys []string
	if err := c.cache.Get(storedKeysKey, keys); err != nil {
		keys = []string{}
	}
	keys = append(keys, key)
	c.cache.Set(storedKeysKey, keys)
	return key
}

func (c *CacheClient) cacheKeyStoredKeys(token string) string {
	return fmt.Sprintf("user-middleware/%s/stored-keys", token)
}

func (c *CacheClient) cleanUp(token string) {
	var keys []string
	c.cache.Delete(c.cacheKeyMe(token))

	storedKeysKey := c.cacheKeyStoredKeys(token)
	err := c.cache.Get(storedKeysKey, &keys)
	if err == nil {
		for _, key := range keys {
			c.cache.Delete(key)
		}
		c.cache.Delete(storedKeysKey)
	}
}
