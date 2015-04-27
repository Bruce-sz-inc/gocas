package authenticator

import "strings"

type Dummy struct{}

func (Dummy) Auth(u string, p string) (bool, string) {
	if strings.TrimSpace(u) != "" && u == p {
		return true, ""
	}
	return false, ""
}
