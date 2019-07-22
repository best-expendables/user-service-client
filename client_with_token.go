package userclient

type ClientWithToken interface {
	Me() (*User, error)
	FindById(userId string) (*User, error)
	FindAll() ([]*User, error)
	RevokedTokens() ([]RevokedToken, error)
}

type ClientWithTokenHolder struct {
	client      Client
	tokenHolder TokenHolder
}

func NewClientWithTokenHolder(c Client, h TokenHolder) *ClientWithTokenHolder {
	return &ClientWithTokenHolder{c, h}
}

func (c *ClientWithTokenHolder) Me() (*User, error) {
	token, err := c.tokenHolder.GetToken(c.client)
	if err != nil {
		return nil, err
	}
	u, err := c.client.Me(token)
	if err == ErrUnauthorized {
		c.tokenHolder.Invalidate()
	}
	return u, err
}

func (c *ClientWithTokenHolder) FindById(userId string) (*User, error) {
	token, err := c.tokenHolder.GetToken(c.client)
	if err != nil {
		return nil, err
	}
	u, err := c.client.FindById(token, userId)
	if err == ErrUnauthorized {
		c.tokenHolder.Invalidate()
	}
	return u, err
}

func (c *ClientWithTokenHolder) FindAll() ([]*User, error) {
	token, err := c.tokenHolder.GetToken(c.client)
	if err != nil {
		return nil, err
	}
	users, err := c.client.FindAll(token)
	if err == ErrUnauthorized {
		c.tokenHolder.Invalidate()
	}
	return users, err
}

func (c *ClientWithTokenHolder) RevokedTokens() ([]RevokedToken, error) {
	token, err := c.tokenHolder.GetToken(c.client)
	if err != nil {
		return nil, err
	}
	tokens, err := c.client.RevokedTokens(token)
	if err == ErrUnauthorized {
		c.tokenHolder.Invalidate()
	}
	return tokens, err
}
