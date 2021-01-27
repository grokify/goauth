package scim

import (
	"fmt"
	"strings"
)

type UserSet struct {
	Users []User
}

func (set *UserSet) Count() int {
	return len(set.Users)
}

func (set *UserSet) GetByEmail(emailAddress string) *User {
	emailAddress = strings.TrimSpace(emailAddress)
	matches := []User{}
	for _, user := range set.Users {
		if emailAddress == user.EmailAddress() {
			matches = append(matches, user)
		}
	}
	if len(matches) > 1 {
		panic(fmt.Sprintf("non-unique email address in set [%s]", emailAddress))
	} else if len(matches) == 1 {
		return &matches[0]
	}
	return nil
}
