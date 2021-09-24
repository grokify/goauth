package scim

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type UserSet struct {
	Users []User `json:"users,omitempty"`
}

func NewUserSet() UserSet {
	return UserSet{Users: []User{}}
}

func ReadFileUserSet(filename string) (UserSet, error) {
	set := NewUserSet()
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return set, err
	}
	return set, json.Unmarshal(bytes, &set)
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
