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
	case ServiceAtlassian:
		return oauth2.Endpoint{
			AuthURL:   AtlassianAuthzURL,
			TokenURL:  AtlassianTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceEbay:
		return oauth2.Endpoint{
			AuthURL:   EbayAuthzURL,
			TokenURL:  EbayTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceEbaySandbox:
		return oauth2.Endpoint{
			AuthURL:   EbayAuthzURLSandbox,
			TokenURL:  EbayTokenURLSandbox,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceFacebook:
		return oauth2.Endpoint{
			AuthURL:   FacebookAuthzURL,
			TokenURL:  FacebookTokenURL,
			AuthStyle: oauth2.AuthStyleInParams}, nil
	case ServiceGithub:
		return oauth2.Endpoint{
			AuthURL:   GithubAuthzURL,
			TokenURL:  GithubTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
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
	case ServiceMailchimp:
		return oauth2.Endpoint{
			AuthURL:   MailchimpAuthzURL,
			TokenURL:  MailchimpTokenURL,
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
	case ServicePaypal:
		return oauth2.Endpoint{
			AuthURL:   PaypalAuthzURL,
			TokenURL:  PaypalTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServicePaypalSandbox:
		return oauth2.Endpoint{
			AuthURL:   PaypalAuthzURLSandbox,
			TokenURL:  PaypalTokenURLSandbox,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServicePipedrive:
		return oauth2.Endpoint{
			AuthURL:   PipedriveAuthzURL,
			TokenURL:  PipedriveTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServicePracticesuite:
		return oauth2.Endpoint{
			AuthURL:   "",
			TokenURL:  PracticesuiteTokenURL,
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
	case ServiceShippo:
		return oauth2.Endpoint{
			AuthURL:   ShippoAuthzURL,
			TokenURL:  ShippoTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceShopify:
		subdomain = strings.TrimSpace(subdomain)
		if len(subdomain) == 0 {
			return oauth2.Endpoint{}, fmt.Errorf("service [%s] requires subddomain", ServiceShopify)
		}
		return oauth2.Endpoint{
			AuthURL:   fmt.Sprintf(ShopifyAuthzURLFormat, subdomain),
			TokenURL:  fmt.Sprintf(ShopifyTokenURLFormat, subdomain),
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceSlack:
		return oauth2.Endpoint{
			AuthURL:   SlackAuthzURL,
			TokenURL:  SlackTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
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
	case ServiceTodoist:
		return oauth2.Endpoint{
			AuthURL:   TodoistAuthzURL,
			TokenURL:  TodoistTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceUber:
		return oauth2.Endpoint{
			AuthURL:   UberAuthzURL,
			TokenURL:  UberTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceWepay:
		return oauth2.Endpoint{
			AuthURL:   WepayAuthzURL,
			TokenURL:  WepayTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceWepaySandbox:
		return oauth2.Endpoint{
			AuthURL:   WepayAuthzURLSandbox,
			TokenURL:  WepayTokenURLSandbox,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceWrike:
		return oauth2.Endpoint{
			AuthURL:   WrikeAuthzURL,
			TokenURL:  WrikeTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	case ServiceWunderlist:
		return oauth2.Endpoint{
			AuthURL:   WunderlistAuthzURL,
			TokenURL:  WunderlistTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect}, nil
	}
	return oauth2.Endpoint{}, fmt.Errorf("service not found [%s]", serviceName)
}
