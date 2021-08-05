package endpoints

import (
	"fmt"
	"strings"

	"golang.org/x/oauth2"
)

const (
	ServiceGoogle             = "google"
	ServiceRingCentral        = "ringcentral"
	ServiceRingCentralSandbox = "ringcentralsandbox"
)

func NewEndpoint(serviceName string) (oauth2.Endpoint, error) {
	switch strings.ToLower(strings.TrimSpace(serviceName)) {
	case ServiceRingCentral:
		return oauth2.Endpoint{
			AuthURL:   RingCentralAuthURL,
			TokenURL:  RingCentralTokenURL,
			AuthStyle: oauth2.AuthStyleInHeader}, nil
	case ServiceRingCentralSandbox:
		return oauth2.Endpoint{
			AuthURL:   RingCentralAuthURLSandbox,
			TokenURL:  RingCentralTokenURLSandbox,
			AuthStyle: oauth2.AuthStyleInHeader}, nil
	}
	return oauth2.Endpoint{}, fmt.Errorf("service not found [%s]", serviceName)
}
