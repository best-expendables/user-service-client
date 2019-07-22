package userclient

import (
	"testing"
)

func TestUser_HasRole(t *testing.T) {
	cases := []struct {
		Roles  []string
		User   User
		Expect bool
	}{
		{
			[]string{"admin", "tester"},
			User{Roles: []string{"admin"}},
			true,
		},
		{
			[]string{"admin", "tester"},
			User{Roles: []string{"pm", "qa", "gamer"}},
			false,
		},
	}

	for _, c := range cases {
		result := c.User.HasRole(c.Roles...)
		if result != c.Expect {
			t.Errorf("Expect %v, Got %v", c.Expect, result)
		}
	}
}
