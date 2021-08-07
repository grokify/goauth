package endpoints

import (
	"fmt"
	"strings"

	"golang.org/x/oauth2"
)

func NewEndpoint(serviceName, subdomain string) (oauth2.Endpoint, error) {
	switch strings.ToLower(strings.TrimSpace(serviceName)) {
	case ServiceAha:
		subdomain = strings.TrimSpace(subdomain)
		if len(subdomain) == 0 {
			return oauth2.Endpoint{}, fmt.Errorf("service [%s] requires subddomain", ServiceAha)
		}
		return oauth2.Endpoint{
			AuthURL:   fmt.Sprintf(AhaAuthzURLFormat, subdomain),
			TokenURL:  fmt.Sprintf(AhaTokenURLFormat, subdomain),
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceAsana:
		return oauth2.Endpoint{
			AuthURL:   AsanaAuthzURL,
			TokenURL:  AsanaTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceFacebook:
		return oauth2.Endpoint{
			AuthURL:   FacebookAuthzURL,
			TokenURL:  FacebookTokenURL,
			AuthStyle: oauth2.AuthStyleInParams}, nil
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
	case ServiceInstagram:
		return oauth2.Endpoint{
			AuthURL:   InstagramAuthzURL,
			TokenURL:  InstagramTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceLyft:
		return oauth2.Endpoint{
			AuthURL:   LyftAuthzURL,
			TokenURL:  LyftTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceMonday:
		return oauth2.Endpoint{
			AuthURL:   MondayAuthzURL,
			TokenURL:  MondayTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServicePagerduty:
		return oauth2.Endpoint{
			AuthURL:   PagerdutyAuthzURL,
			TokenURL:  PagerdutyTokenURL,
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
	case ServiceStackoverflow:
		return oauth2.Endpoint{
			AuthURL:   StackoverflowAuthzURL,
			TokenURL:  StackoverflowTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceStripe:
		return oauth2.Endpoint{
			AuthURL:   StripeAuthzURL,
			TokenURL:  StripeTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	}
	return oauth2.Endpoint{}, fmt.Errorf("service not found [%s]", serviceName)
}
