package goxAuthSimpleS2S

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type authHeaders struct {
	ClientId     string `header:"X-Client-Id" binding:"required"`
	ClientSecret string `header:"X-Client-Secret" binding:"required"`
	ClientName   string `header:"X-Client-Name" binding:"-"`
}

func (auth *AuthConfig) GinHandlerFuncWithAccessCheck(accessToCheck []string, handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := auth.ginValidate(c, accessToCheck); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
		} else {
			handlerFunc(c)
		}
	}
}

func (auth *AuthConfig) ginValidate(c *gin.Context, accessToCheck []string) error {
	if auth.Disabled {
		return nil
	}

	headers := &authHeaders{}
	if err := c.ShouldBindHeader(headers); err != nil {
		return &AuthError{
			Reason:   "missing_auth_headers",
			ClientId: headers.ClientId,
			Err:      err,
		}
	}

	// Make sure client id is valid
	if _, ok := auth.Clients[headers.ClientId]; !ok {
		return &AuthError{
			Reason:   "invalid_client_id",
			ClientId: headers.ClientId,
		}
	}

	// Make sure client secret is valid
	if !auth.checkClientSecret(c, headers.ClientId, headers.ClientSecret) {
		return &AuthError{
			Reason:   "client_secret_mismatch",
			ClientId: headers.ClientId,
		}
	}

	// If action * is allowed then no need to check anything
	if auth.Clients[headers.ClientId].AllowedActions != nil {
		for _, action := range auth.Clients[headers.ClientId].AllowedActions {
			if action == "*" {
				return nil
			}
		}
	}

	// Make sure client has access to the action
	for _, action := range accessToCheck {
		for _, allowedAction := range auth.Clients[headers.ClientId].AllowedActions {
			if action == allowedAction {
				return nil
			}
		}
	}

	return &AuthError{
		Reason:   fmt.Sprintf("client_does_not_have_access_to_actions__%s", strings.Join(accessToCheck, "_")),
		ClientId: headers.ClientId,
	}
}

func (auth *AuthConfig) checkClientSecret(c *gin.Context, clientId string, inputClientSecret string) bool {
	if _, ok := auth.Clients[clientId]; ok {
		for _, secret := range strings.Split(auth.Clients[clientId].ClientSecret, ",") {
			if inputClientSecret == secret {
				return true
			}
		}
	}
	return false
}
