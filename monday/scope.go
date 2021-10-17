package monday

import "github.com/grokify/goauth"

func GetScopes() []goauth.Scope {
	return []goauth.Scope{
		{Name: "me:read", Description: "Read your basic personal details"},
		{Name: "boards:read", Description: "Read boards data"},
		{Name: "boards:write", Description: "Modify boards data"},
		{Name: "workspaces:read", Description: "Read user's workspaces data"},
		{Name: "workspaces:write", Description: "Modify user's workspaces data"},
		{Name: "users:read", Description: "Access data about your account's users"},
		{Name: "users:write", Description: "Modify your account's users data"},
		{Name: "account:read", Description: "Access information about your account"},
		{Name: "notifications:write", Description: "Send notifications to users"},
		{Name: "updates:read", Description: "Read updates data"},
		{Name: "updates:write", Description: "Modify updates data"},
		{Name: "assets:read", Description: "Read information of assets that the user has access to"},
		{Name: "tags:read", Description: "Read tags data"},
		{Name: "teams:read", Description: "Read teams data"},
		{Name: "webhooks:write", Description: "Create new webhooks"},
	}
}
