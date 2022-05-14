package kid

// Http header
const (
	HeaderRequestId                     = "X-Request-ID"
	HeaderContentDisposition            = "Content-Disposition"
	HeaderContentType                   = "Content-Type"
	HeaderLocation                      = "Location"
	HeaderAuthorization                 = "Authorization"
	HeaderWWWAuthenticate               = "WWW-Authenticate"
	HeaderOrigin                        = "Origin"
	HeaderVary                          = "Vary"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
)

const (
	CtxRequestId         = "Ctx-Request-ID"
	CtxBasicAuthUsername = "Ctx-Basic-Auth-Username"
	CtxBasicAuthPassword = "Ctx-Basic-Auth-Password"
)
