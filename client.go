package userclient

type Client interface {
	Authenticate(username, password string) (string, error)

	Me(token string) (*User, error)
	Logout(token string) error

	FindById(token, userId string) (*User, error)
	FindAll(token string) ([]*User, error)

	RevokedTokens(token string) ([]RevokedToken, error)
}
