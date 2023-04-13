package ringcentral

import (
	"strings"

	"github.com/grokify/goauth"
	"github.com/grokify/mogo/crypto/argon2util"
)

func UsernameExtensionPasswordToString(username, password string) string {
	return strings.Join([]string{
		strings.TrimSpace(username),
		strings.TrimSpace(password)}, "\t")
}

func UsernameExtensionPasswordToHash(username, extension, password string, salt []byte) string {
	return argon2util.HashSimpleBase36(
		[]byte(UsernameExtensionPasswordToString(username, password)),
		salt)
}

func PasswordCredentialsToHash(pwdCreds goauth.CredentialsOAuth2, salt []byte) string {
	return argon2util.HashSimpleBase36(
		[]byte(UsernameExtensionPasswordToString(
			pwdCreds.Username, pwdCreds.Password)),
		salt)
}
