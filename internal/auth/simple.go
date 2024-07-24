package auth

type SimpleAuthProvider struct {
	User string
	Pass string
}

func (s *SimpleAuthProvider) Verify(user, pass string) bool {
	return user == s.User && pass == s.Pass
}
