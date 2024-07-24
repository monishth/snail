package auth

import (
	"github.com/charmbracelet/log"
	"github.com/tg123/go-htpasswd"
)

type CredentialMap map[string]string

type HtpasswdProvider struct {
	file *htpasswd.File
}

func NewHtpasswdProvider(filename string) HtpasswdProvider {
	file, err := htpasswd.New(filename, htpasswd.DefaultSystems, nil)
	if err != nil {
		panic(err)
	}
	log.Info("Loaded htpasswd file successfully")

	return HtpasswdProvider{file}
}

func (h *HtpasswdProvider) Verify(user, pass string) bool {
	return h.file.Match(user, pass)
}
