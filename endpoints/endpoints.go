package endpoints

import (
	"fmt"
	"strings"

	"golang.org/x/oauth2"
)

func NewEndpoint(serviceName string) (oauth2.Endpoint, error) {
	switch strings.ToLower(strings.TrimSpace(serviceName)) {
	case ServiceGoogle:
		return oauth2.Endpoint{
			AuthURL:   GoogleAuthzURL,
			TokenURL:  GoogleTokenURL,
			AuthStyle: oauth2.AuthStyleInParams}, nil
	case ServiceHubspot:
		return oauth2.Endpoint{
			AuthURL:   HubspotAuthzURL,
			TokenURL:  HubspotTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceMonday:
		return oauth2.Endpoint{
			AuthURL:   MondayAuthzURL,
			TokenURL:  MondayTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceRingcentral:
		return oauth2.Endpoint{
			AuthURL:   RingcentralAuthzURL,
			TokenURL:  RingcentralTokenURL,
			AuthStyle: oauth2.AuthStyleInHeader}, nil
	case ServiceRingcentralSandbox:
		return oauth2.Endpoint{
			AuthURL:   RingcentralAuthzURLSandbox,
			TokenURL:  RingcentralTokenURLSandbox,
			AuthStyle: oauth2.AuthStyleInHeader}, nil
	}
	return oauth2.Endpoint{}, fmt.Errorf("service not found [%s]", serviceName)
}
