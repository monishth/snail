package auth

type SimpleAuthProvider struct {
	User string
	Pass string
}

func (s *SimpleAuthProvider) Verify(user, pass string) bool {
	// TODO: Should be constant time compare
	return user == s.User && pass == s.Pass
}
