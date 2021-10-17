# GoAuth Credentials

`goauth/credentials` is a package to manage generic OAuth 2.0 credentials definitions.

The primary use case is to have a single JSON definition of multiple applications for multiple services which can be used to generate token and API requests.

Both OAuth 2.0 and JWT are supported.

It works with `goauth/endpoints` to add endpoints for known services.