package goxAuthSimpleS2S

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/devlibx/gox-base/v2/errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type authHeaders struct {
	ClientId     string `header:"X-Client-Id" binding:"required"`
	ClientSecret string `header:"X-Client-Secret" binding:"required"`
	ClientName   string `header:"X-Client-Name" binding:"-"`
}

func (auth *AuthConfig) GinHandlerFuncWithAccessCheck(accessToCheck []string, handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := auth.ginValidate(c, accessToCheck); err != nil {
			if e, ok := errors.AsTyped[*AuthError](err); ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, e)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				c.Abort()
			}
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
	if headers.ClientId != "1" && !auth.checkClientSecret(c, headers.ClientId, headers.ClientSecret) {
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

	if checkClientSecretForAdhocUser, err := auth.checkClientSecretForAdhocUser(c, headers.ClientId, headers.ClientSecret); checkClientSecretForAdhocUser {
		return nil
	} else if err != nil {
		return err
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

func (auth *AuthConfig) checkClientSecretForAdhocUser(c *gin.Context, clientId string, inputClientSecret string) (bool, error) {
	// If client id is 1 then it is the adhoc client
	if clientId != "1" {
		return false, nil
	} else if _, ok := auth.Clients[clientId]; !ok {
		return false, nil
	}

	tokens := strings.Split(inputClientSecret, "-")
	if len(tokens) == 3 {

		// Make sure token is not expired
		t, err := strconv.ParseInt(tokens[1], 10, 64)
		if err != nil {
			return false, &AuthError{
				Reason:   "adhoc_user_timestamp_invalid",
				ClientId: "1",
			}
		}
		timestamp := time.Unix(t, 0)
		if timestamp.Before(time.Now()) {
			return false, &AuthError{
				Reason:   "adhoc_user_secret_expired",
				ClientId: "1",
			}
		}

		// Create hash
		input := fmt.Sprintf("%s-%d-%s", tokens[0], t, auth.SecretSalt)
		hash := md5.Sum([]byte(input))
		encodedString := hex.EncodeToString(hash[:])
		if encodedString != tokens[2] {
			return false, &AuthError{
				Reason:   "adhoc_user_secret_mismatch",
				ClientId: "1",
			}
		}
	} else {
		return false, &AuthError{
			Reason:   "adhoc_user_secret_invalid",
			ClientId: "1",
		}
	}

	slog.Info("Adhoc user access", "client_id", "1", "secret_part", fmt.Sprintf("%s-%s", tokens[0], tokens[1]))
	return true, nil
}

func GenerateAdhocUserSecret(duration time.Duration, userName string, salt string) string {
	t := time.Now().Add(duration)
	input := fmt.Sprintf("%s-%d-%s", userName, t.Unix(), salt)
	hash := md5.Sum([]byte(input))
	return fmt.Sprintf("%s-%d-%s", userName, t.Unix(), hex.EncodeToString(hash[:]))
}
