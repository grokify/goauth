package endpoints

const (
	// Service names must be lower case alpha-numeric.
	ServiceAha                = "aha"
	ServiceAsana              = "asana"
	ServiceFacebook           = "facebook"
	ServiceGoogle             = "google"
	ServiceHubspot            = "hubspot"
	ServiceInstagram          = "instagram"
	ServiceLyft               = "lyft"
	ServiceMonday             = "monday"
	ServicePagerduty          = "pagerduty"
	ServiceRingcentral        = "ringcentral"
	ServiceRingcentralSandbox = "ringcentralsandbox"
	ServiceStackoverflow      = "stackoverflow"
	ServiceStripe             = "stripe"

	AhaAuthzURLFormat          = "https://%s.aha.io/oauth/authorize"
	AhaTokenURLFormat          = "https://%s.aha.io/oauth/token"
	AsanaAuthzURL              = "https://app.asana.com/-/oauth_authorize"
	AsanaTokenURL              = "https://app.asana.com/-/oauth_token"
	FacebookAuthzURL           = "https://www.facebook.com/v3.2/dialog/oauth"
	FacebookTokenURL           = "https://graph.facebook.com/v3.2/oauth/access_token"
	GoogleAuthzURL             = "https://accounts.google.com/o/oauth2/auth"
	GoogleTokenURL             = "https://oauth2.googleapis.com/token"
	HubspotAuthzURL            = "https://app.hubspot.com/oauth/authorize"
	HubspotTokenURL            = "https://api.hubapi.com/oauth/v1/token"
	InstagramAuthzURL          = "https://api.instagram.com/oauth/authorize"
	InstagramTokenURL          = "https://api.instagram.com/oauth/access_token"
	LyftAuthzURL               = "https://www.lyft.com/oauth/authorize"
	LyftTokenURL               = "https://api.lyft.com/oauth/token"
	MondayAuthzURL             = "https://auth.monday.com/oauth2/authorize"
	MondayTokenURL             = "https://auth.monday.com/oauth2/token"
	PagerdutyAuthzURL          = "https://app.pagerduty.com/oauth/authorize"
	PagerdutyTokenURL          = "https://app.pagerduty.com/oauth/token"
	RingcentralAuthzURL        = "https://platform.ringcentral.com/restapi/oauth/authorize"
	RingcentralTokenURL        = "https://platform.ringcentral.com/restapi/oauth/token"
	RingcentralAuthzURLSandbox = "https://platform.devtest.ringcentral.com/restapi/oauth/authorize"
	RingcentralTokenURLSandbox = "https://platform.devtest.ringcentral.com/restapi/oauth/token"
	StackoverflowAuthzURL      = "https://stackoverflow.com/oauth"
	StackoverflowTokenURL      = "https://stackoverflow.com/oauth/access_token"
	StripeAuthzURL             = "https://connect.stripe.com/oauth/authorize"
	StripeTokenURL             = "https://connect.stripe.com/oauth/token"
)
