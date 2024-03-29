package endpoints

const (
	// Service names must be lower case alpha-numeric.
	ServiceAha                = "aha"
	ServiceAsana              = "asana"
	ServiceAtlassian          = "atlassian"
	ServiceBackblaze          = "backblaze" // basic auth
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
	ServiceZoom               = "zoom"

	AhaAuthzURLFormat           = "https://%s.aha.io/oauth/authorize"
	AhaTokenURLFormat           = "https://%s.aha.io/oauth/token" // #nosec G101
	AhaServerURLFormat          = "https://%s.aha.io"
	AsanaAuthzURL               = "https://app.asana.com/-/oauth_authorize"
	AsanaTokenURL               = "https://app.asana.com/-/oauth_token" // #nosec G101
	AtlassianAuthzURL           = "https://auth.atlassian.com/authorize"
	AtlassianTokenURL           = "https://auth.atlassian.com/oauth/token" // #nosec G101
	EbayAuthzURL                = "https://auth.ebay.com/oauth2/authorize"
	EbayTokenURL                = "https://api.ebay.com/identity/v1/oauth2/token" // #nosec G101
	EbaySandboxAuthzURL         = "https://auth.sandbox.ebay.com/oauth2/authorize"
	EbaySandboxTokenURL         = "https://api.sandbox.ebay.com/identity/v1/oauth2/token" // #nosec G101
	FacebookAuthzURL            = "https://www.facebook.com/v3.2/dialog/oauth"
	FacebookTokenURL            = "https://graph.facebook.com/v3.2/oauth/access_token" // #nosec G101
	GithubAuthzURL              = "https://github.com/login/oauth/authorize"
	GithubTokenURL              = "https://github.com/login/oauth/access_token" // #nosec G101
	GithubServerURL             = "https://api.github.com"
	GoogleAuthzURL              = "https://accounts.google.com/o/oauth2/auth"
	GoogleTokenURL              = "https://oauth2.googleapis.com/token" // #nosec G101
	HubspotAuthzURL             = "https://app.hubspot.com/oauth/authorize"
	HubspotTokenURL             = "https://api.hubapi.com/oauth/v1/token" // #nosec G101
	HubspotServerURL            = "https://api.hubapi.com"
	InstagramAuthzURL           = "https://api.instagram.com/oauth/authorize"
	InstagramTokenURL           = "https://api.instagram.com/oauth/access_token" // #nosec G101
	LyftAuthzURL                = "https://www.lyft.com/oauth/authorize"
	LyftTokenURL                = "https://api.lyft.com/oauth/token" // #nosec G101
	MailchimpAuthzURL           = "https://login.mailchimp.com/oauth2/authorize"
	MailchimpTokenURL           = "https://login.mailchimp.com/oauth2/token" // #nosec G101
	MondayAuthzURL              = "https://auth.monday.com/oauth2/authorize"
	MondayTokenURL              = "https://auth.monday.com/oauth2/token" // #nosec G101
	MondayServerURL             = "https://api.monday.com/v2"
	PagerdutyAuthzURL           = "https://app.pagerduty.com/oauth/authorize"
	PagerdutyTokenURL           = "https://app.pagerduty.com/oauth/token" // #nosec G101
	PagerdutyServerURL          = "https://api.pagerduty.com"
	PaypalAuthzURL              = "https://www.paypal.com/webapps/auth/protocol/openidconnect/v1/authorize"
	PaypalTokenURL              = "https://api.paypal.com/v1/identity/openidconnect/tokenservice" // #nosec G101
	PaypalSandboxAuthzURL       = "https://www.sandbox.paypal.com/webapps/auth/protocol/openidconnect/v1/authorize"
	PaypalSandboxTokenURL       = "https://api.sandbox.paypal.com/v1/identity/openidconnect/tokenservice" // #nosec G101
	PipedriveAuthzURL           = "https://oauth.pipedrive.com/oauth/authorize"
	PipedriveTokenURL           = "https://oauth.pipedrive.com/oauth/token"           // #nosec G101
	PracticesuiteTokenURL       = "https://staging.practicesuite.com/uaa/oauth/token" // #nosec G101
	RingcentralAuthzURL         = "https://platform.ringcentral.com/restapi/oauth/authorize"
	RingcentralTokenURL         = "https://platform.ringcentral.com/restapi/oauth/token" // #nosec G101
	RingcentralServerURL        = "https://platform.ringcentral.com"
	RingcentralSandboxAuthzURL  = "https://platform.devtest.ringcentral.com/restapi/oauth/authorize"
	RingcentralSandboxTokenURL  = "https://platform.devtest.ringcentral.com/restapi/oauth/token" // #nosec G101
	RingcentralSandboxServerURL = "https://platform.devtest.ringcentral.com"
	SalesforceAuthzURL          = "https://login.salesforce.com/services/oauth2/authorize"
	SalesforceTokenURL          = "https://login.salesforce.com/services/oauth2/token" // #nosec G101
	SalesforceRevokeURL         = "https://login.salesforce.com/services/oauth2/revoke"
	ShippoAuthzURL              = "https://goshippo.com/oauth/authorize"
	ShippoTokenURL              = "https://goshippo.com/oauth/access_token" // #nosec G101
	ShopifyAuthzURLFormat       = "https://%s.myshopify.com/admin/oauth/authorize"
	ShopifyTokenURLFormat       = "https://%s.myshopify.com/admin/oauth/access_token" // #nosec G101
	SlackAuthzURL               = "https://slack.com/oauth/authorize"
	SlackTokenURL               = "https://slack.com/api/oauth.access" // #nosec G101
	SlackServerURL              = "https://slack.com/api"
	StackoverflowAuthzURL       = "https://stackoverflow.com/oauth"
	StackoverflowTokenURL       = "https://stackoverflow.com/oauth/access_token" // #nosec G101
	StackoverflowServerURL      = "https://api.stackexchange.com/2.2"
	StripeAuthzURL              = "https://connect.stripe.com/oauth/authorize"
	StripeTokenURL              = "https://connect.stripe.com/oauth/token" // #nosec G101
	TodoistAuthzURL             = "https://todoist.com/oauth/authorize"
	TodoistTokenURL             = "https://todoist.com/oauth/access_token" // #nosec G101
	TwitterTokenURL             = "https://api.twitter.com/oauth2/token"   // #nosec G101
	UberAuthzURL                = "https://login.uber.com/oauth/v2/authorize"
	UberTokenURL                = "https://login.uber.com/oauth/v2/token" // #nosec G101
	WepayAuthzURL               = "https://www.wepay.com/v2/oauth2/authorize"
	WepayTokenURL               = "https://wepayapi.com/v2/oauth2/token" // #nosec G101
	WepaySandboxAuthzURL        = "https://stage.wepay.com/v2/oauth2/authorize"
	WepaySandboxTokenURL        = "https://stage.wepayapi.com/v2/oauth2/token" // #nosec G101
	WrikeAuthzURL               = "https://login.wrike.com/oauth2/authorize/v4"
	WrikeTokenURL               = "https://login.wrike.com/oauth2/token" // #nosec G101
	WunderlistAuthzURL          = "https://www.wunderlist.com/oauth/authorize"
	WunderlistTokenURL          = "https://www.wunderlist.com/oauth/access_token" // #nosec G101
	ZoomAuthzURL                = "https://zoom.us/oauth/authorize"
	ZoomTokenURL                = "https://zoom.us/oauth/token" // #nosec G101
	ZoomServerURL               = "https://api.zoom.us/v2"
	ZoomJWTSigningMethod        = "HS256"
)
