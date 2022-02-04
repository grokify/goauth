package ringcentral

import (
	"strings"

	"github.com/grokify/goauth/credentials"
	"github.com/grokify/mogo/crypto/hash/argon2"
)

func UsernameExtensionPasswordToString(username, password string) string {
	return strings.Join([]string{
		strings.TrimSpace(username),
		strings.TrimSpace(password)}, "\t")
}

func UsernameExtensionPasswordToHash(username, extension, password string, salt []byte) string {
	return argon2.HashSimpleBase36(
		[]byte(UsernameExtensionPasswordToString(username, password)),
		salt)
}

func PasswordCredentialsToHash(pwdCreds credentials.CredentialsOAuth2, salt []byte) string {
	return argon2.HashSimpleBase36(
		[]byte(UsernameExtensionPasswordToString(
			pwdCreds.Username, pwdCreds.Password)),
		salt)
}
