package endpoints

const (
	// Service names must be lower case alpha-numeric.
	ServiceGoogle             = "google"
	ServiceHubspot            = "hubspot"
	ServiceMonday             = "monday"
	ServiceRingcentral        = "ringcentral"
	ServiceRingcentralSandbox = "ringcentralsandbox"

	GoogleAuthzURL             = "https://accounts.google.com/o/oauth2/auth"
	GoogleTokenURL             = "https://oauth2.googleapis.com/token"
	HubspotAuthzURL            = "https://app.hubspot.com/oauth/authorize"
	HubspotTokenURL            = "https://api.hubapi.com/oauth/v1/token"
	MondayAuthzURL             = "https://auth.monday.com/oauth2/authorize"
	MondayTokenURL             = "https://auth.monday.com/oauth2/token"
	RingcentralAuthzURL        = "https://platform.ringcentral.com/restapi/oauth/authorize"
	RingcentralTokenURL        = "https://platform.ringcentral.com/restapi/oauth/token"
	RingcentralAuthzURLSandbox = "https://platform.devtest.ringcentral.com/restapi/oauth/authorize"
	RingcentralTokenURLSandbox = "https://platform.devtest.ringcentral.com/restapi/oauth/token"
)
