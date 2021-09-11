package endpoints

const (
	// Service names must be lower case alpha-numeric.
	ServiceAha                = "aha"
	ServiceAsana              = "asana"
	ServiceAtlassian          = "atlassian"
	ServiceEbay               = "ebay"
	ServiceEbaySandbox        = "ebaysandbox"
	ServiceFacebook           = "facebook"
	ServiceGithub             = "github"
	ServiceGoogle             = "google"
	ServiceHubspot            = "hubspot"
	ServiceInstagram          = "instagram"
	ServiceLyft               = "lyft"
	ServiceMailchimp          = "mailchimp"
	ServiceMonday             = "monday"
	ServicePagerduty          = "pagerduty"
	ServicePaypal             = "paypal"
	ServicePaypalSandbox      = "paypalsandbox"
	ServicePipedrive          = "pipedrive"
	ServicePracticesuite      = "practicesuite"
	ServiceRingcentral        = "ringcentral"
	ServiceRingcentralSandbox = "ringcentralsandbox"
	ServiceShippo             = "shippo"
	ServiceShopify            = "shopify"
	ServiceSlack              = "slack"
	ServiceStackoverflow      = "stackoverflow"
	ServiceStripe             = "stripe"
	ServiceTodoist            = "todoist"
	ServiceTwitter            = "twitter"
	ServiceUber               = "uber"
	ServiceWepay              = "wepay"
	ServiceWepaySandbox       = "wepaysandbox"
	ServiceWrike              = "wrike"
	ServiceWunderlist         = "wunderlist"

	AhaAuthzURLFormat           = "https://%s.aha.io/oauth/authorize"
	AhaTokenURLFormat           = "https://%s.aha.io/oauth/token"
	AhaServerURLFormat          = "https://%s.aha.io"
	AsanaAuthzURL               = "https://app.asana.com/-/oauth_authorize"
	AsanaTokenURL               = "https://app.asana.com/-/oauth_token"
	AtlassianAuthzURL           = "https://auth.atlassian.com/authorize"
	AtlassianTokenURL           = "https://auth.atlassian.com/oauth/token"
	EbayAuthzURL                = "https://auth.ebay.com/oauth2/authorize"
	EbayTokenURL                = "https://api.ebay.com/identity/v1/oauth2/token"
	EbaySandboxAuthzURL         = "https://auth.sandbox.ebay.com/oauth2/authorize"
	EbaySandboxTokenURL         = "https://api.sandbox.ebay.com/identity/v1/oauth2/token"
	FacebookAuthzURL            = "https://www.facebook.com/v3.2/dialog/oauth"
	FacebookTokenURL            = "https://graph.facebook.com/v3.2/oauth/access_token"
	GithubAuthzURL              = "https://github.com/login/oauth/authorize"
	GithubTokenURL              = "https://github.com/login/oauth/access_token"
	GithubServerURL             = "https://api.github.com"
	GoogleAuthzURL              = "https://accounts.google.com/o/oauth2/auth"
	GoogleTokenURL              = "https://oauth2.googleapis.com/token"
	HubspotAuthzURL             = "https://app.hubspot.com/oauth/authorize"
	HubspotTokenURL             = "https://api.hubapi.com/oauth/v1/token"
	InstagramAuthzURL           = "https://api.instagram.com/oauth/authorize"
	InstagramTokenURL           = "https://api.instagram.com/oauth/access_token"
	LyftAuthzURL                = "https://www.lyft.com/oauth/authorize"
	LyftTokenURL                = "https://api.lyft.com/oauth/token"
	MailchimpAuthzURL           = "https://login.mailchimp.com/oauth2/authorize"
	MailchimpTokenURL           = "https://login.mailchimp.com/oauth2/token"
	MondayAuthzURL              = "https://auth.monday.com/oauth2/authorize"
	MondayTokenURL              = "https://auth.monday.com/oauth2/token"
	MondayServerURL             = "https://api.monday.com/v2"
	PagerdutyAuthzURL           = "https://app.pagerduty.com/oauth/authorize"
	PagerdutyTokenURL           = "https://app.pagerduty.com/oauth/token"
	PagerdutyServerURL          = "https://api.pagerduty.com"
	PaypalAuthzURL              = "https://www.paypal.com/webapps/auth/protocol/openidconnect/v1/authorize"
	PaypalTokenURL              = "https://api.paypal.com/v1/identity/openidconnect/tokenservice"
	PaypalSandboxAuthzURL       = "https://www.sandbox.paypal.com/webapps/auth/protocol/openidconnect/v1/authorize"
	PaypalSandboxTokenURL       = "https://api.sandbox.paypal.com/v1/identity/openidconnect/tokenservice"
	PipedriveAuthzURL           = "https://oauth.pipedrive.com/oauth/authorize"
	PipedriveTokenURL           = "https://oauth.pipedrive.com/oauth/token"
	PracticesuiteTokenURL       = "https://staging.practicesuite.com/uaa/oauth/token"
	RingcentralAuthzURL         = "https://platform.ringcentral.com/restapi/oauth/authorize"
	RingcentralTokenURL         = "https://platform.ringcentral.com/restapi/oauth/token"
	RingcentralServerURL        = "https://platform.ringcentral.com"
	RingcentralSandboxAuthzURL  = "https://platform.devtest.ringcentral.com/restapi/oauth/authorize"
	RingcentralSandboxTokenURL  = "https://platform.devtest.ringcentral.com/restapi/oauth/token"
	RingcentralSandboxServerURL = "https://platform.devtest.ringcentral.com"
	ShippoAuthzURL              = "https://goshippo.com/oauth/authorize"
	ShippoTokenURL              = "https://goshippo.com/oauth/access_token"
	ShopifyAuthzURLFormat       = "https://%s.myshopify.com/admin/oauth/authorize"
	ShopifyTokenURLFormat       = "https://%s.myshopify.com/admin/oauth/access_token"
	SlackAuthzURL               = "https://slack.com/oauth/authorize"
	SlackTokenURL               = "https://slack.com/api/oauth.access"
	StackoverflowAuthzURL       = "https://stackoverflow.com/oauth"
	StackoverflowTokenURL       = "https://stackoverflow.com/oauth/access_token"
	StripeAuthzURL              = "https://connect.stripe.com/oauth/authorize"
	StripeTokenURL              = "https://connect.stripe.com/oauth/token"
	TodoistAuthzURL             = "https://todoist.com/oauth/authorize"
	TodoistTokenURL             = "https://todoist.com/oauth/access_token"
	TwitterTokenURL             = "https://api.twitter.com/oauth2/token"
	UberAuthzURL                = "https://login.uber.com/oauth/v2/authorize"
	UberTokenURL                = "https://login.uber.com/oauth/v2/token"
	WepayAuthzURL               = "https://www.wepay.com/v2/oauth2/authorize"
	WepayTokenURL               = "https://wepayapi.com/v2/oauth2/token"
	WepaySandboxAuthzURL        = "https://stage.wepay.com/v2/oauth2/authorize"
	WepaySandboxTokenURL        = "https://stage.wepayapi.com/v2/oauth2/token"
	WrikeAuthzURL               = "https://login.wrike.com/oauth2/authorize/v4"
	WrikeTokenURL               = "https://login.wrike.com/oauth2/token"
	WunderlistAuthzURL          = "https://www.wunderlist.com/oauth/authorize"
	WunderlistTokenURL          = "https://www.wunderlist.com/oauth/access_token"
)
