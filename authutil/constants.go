package authutil

const (
	GrantTypeAccountCredentials = "account_credentials" // used by only Zoom?
	GrantTypeAuthorizationCode  = "authorization_code"
	GrantTypeClientCredentials  = "client_credentials"
	GrantTypeJWTBearer          = "urn:ietf:params:oauth:grant-type:jwt-bearer" // #nosec G101
	GrantTypeSAML2Bearer        = "urn:ietf:params:oauth:grant-type:saml2-bearer"
	GrantTypePassword           = "password"
	GrantTypeRefreshToken       = "refresh_token"
	GrantTypeCustomStatic       = "custom_static"

	ParamAssertion    = "assertion"
	ParamClientID     = "client_id"
	ParamGrantType    = "grant_type"
	ParamScope        = "scope"
	ParamPassword     = "password"
	ParamUsername     = "usernamae"
	ParamRefreshToken = "refresh_token"

	TokenBasic  = "Basic"
	TokenBearer = "Bearer"

	OAuth2TokenPropAccessToken  = "access_token"
	OAuth2TokenPropExpiresIn    = "expires_in"
	OAuth2TokenPropRefreshToken = "refresh_token"
	OAuth2TokenPropTokenType    = "token_type"

	TestRedirectURL = "https://grokify.github.io/goauth/oauth2callback/"
)
