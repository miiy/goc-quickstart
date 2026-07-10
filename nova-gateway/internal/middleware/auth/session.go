package auth

import (
	"net/http"
	"strings"

	"github.com/miiy/goc-quickstart/nova-gateway/internal/transport"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/middleware/csrf"
	"github.com/miiy/goc/gin/sessions"
	"google.golang.org/grpc/codes"
)

const (
	sessionAccessTokenKey       = "access_token"
	sessionAuthenticationSource = "nova-gateway.auth.session"
)

// SessionCSRF protects unsafe requests authenticated through the browser's
// session cookie. Explicit Authorization credentials are not ambient browser
// credentials and continue through the bearer-token flow without CSRF checks.
func SessionCSRF() gin.HandlerFunc {
	validate := csrf.Middleware(csrf.WithUnauthorized(func(c *gin.Context) {
		transport.WriteOpenAPIError(c, http.StatusForbidden, int32(codes.PermissionDenied), "invalid CSRF token")
	}))

	return func(c *gin.Context) {
		fromSession, _ := c.Get(sessionAuthenticationSource)
		if fromSession != true {
			c.Next()
			return
		}

		validate(c)
	}
}

// SessionTokenToBearer bridges nova-web's session-backed access token into the
// existing bearer-token auth flow. Explicit Authorization headers take priority.
func SessionTokenToBearer() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.TrimSpace(c.GetHeader("Authorization")) == "" {
			if token := sessionAccessToken(c); token != "" {
				c.Set(sessionAuthenticationSource, true)
				c.Request.Header.Set("Authorization", "Bearer "+token)
			}
		}
		c.Next()
	}
}

func sessionAccessToken(c *gin.Context) string {
	token, _ := sessions.Default(c).Get(sessionAccessTokenKey).(string)
	return strings.TrimSpace(token)
}
