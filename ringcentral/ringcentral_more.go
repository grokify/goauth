package ringcentral

import (
	"strings"

	"github.com/grokify/gotilla/crypto/hash/argon2"
)

func UsernameExtensionPasswordToString(username, extension, password string) string {
	return strings.Join([]string{
		strings.TrimSpace(username),
		strings.TrimSpace(extension),
		strings.TrimSpace(password)}, "\t")
}

func UsernameExtensionPasswordToHash(username, extension, password string, salt []byte) string {
	return argon2.HashSimpleBase36(
		[]byte(UsernameExtensionPasswordToString(username, extension, password)),
		salt)
}

func PasswordCredentialsToHash(pwdCreds PasswordCredentials, salt []byte) string {
	return argon2.HashSimpleBase36(
		[]byte(UsernameExtensionPasswordToString(
			pwdCreds.Username, pwdCreds.Extension, pwdCreds.Password)),
		salt)
}
