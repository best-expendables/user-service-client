package userclient

type UserClientMock struct {
	AuthenticateMock func(username string, password string) (string, error)
	MeMock           func(token string) (*User, error)
	LogoutMock       func(token string) error

	FindByIdMock      func(token, userId string) (*User, error)
	FindAllMock       func(token string) ([]*User, error)
	RevokedTokensMock func(token string) ([]RevokedToken, error)
}

func (c *UserClientMock) Authenticate(username string, password string) (string, error) {
	return c.AuthenticateMock(username, password)
}

func (c *UserClientMock) Me(token string) (*User, error) {
	return c.MeMock(token)
}

func (c *UserClientMock) Logout(token string) error {
	return c.LogoutMock(token)
}

func (c *UserClientMock) FindById(token, userId string) (*User, error) {
	return c.FindByIdMock(token, userId)
}

func (c *UserClientMock) FindAll(token string) ([]*User, error) {
	return c.FindAllMock(token)
}

func (c *UserClientMock) RevokedTokens(token string) ([]RevokedToken, error) {
	return c.RevokedTokens(token)
}
