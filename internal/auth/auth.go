package auth

type AuthProvider interface {
	Verify(user, pass string) bool
}
