package server

import "errors"

type ServerOptions struct {
	Port              int
	AuthProvider      AuthProvider
	SimpleCredentials string
	HttpasswdFilename string
}

type AuthProvider string

const (
	NoAuth       AuthProvider = "none"
	HtpasswdAuth AuthProvider = "htpasswd"
	SimpleAuth   AuthProvider = "simple"
)

func (ap *AuthProvider) String() string {
	return string(*ap)
}

func (ap *AuthProvider) Set(v string) error {
	switch v {
	case "htpasswd", "simple":
		*ap = AuthProvider(v)
		return nil
	default:
		return errors.New("must be one of \"htpasswd\", \"simple\"")
	}
}

func (ap *AuthProvider) Type() string {
	return "AuthProvider"
}
