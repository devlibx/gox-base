package goxAuthSimpleS2S

import (
	_ "embed"
	"fmt"
	"github.com/devlibx/gox-base/v2/auth"
	"github.com/devlibx/gox-base/v2/serialization"
	goxJsonUtils "github.com/devlibx/gox-base/v2/serialization/utils/json"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

//go:embed test_s2s_simple_auth.yaml
var testS2SSimpleAuthYaml string

type tStruct struct {
	AuthSimpleS2S AuthConfig `yaml:"authSimpleS2S"`
}

func TestClientIdMissing(t *testing.T) {
	tS := tStruct{}
	err := serialization.ReadYamlFromString(testS2SSimpleAuthYaml, &tS)
	assert.NoError(t, err)
	authConfig := tS.AuthSimpleS2S
	assert.False(t, authConfig.Disabled)

	endpointCalled := false
	router := gin.New()
	router.GET("/test", authConfig.GinHandlerFuncWithAccessCheck([]string{"read"}, func(c *gin.Context) {
		endpointCalled = true
	}))
	server := httptest.NewServer(router)
	defer server.Close()

	r := resty.New()
	resp, err := r.R().SetHeader("Content-Type", "application/json").Get(fmt.Sprintf("%s/test", server.URL))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	assert.False(t, endpointCalled)
	fmt.Println(resp.String())
}

func TestWrongSecret(t *testing.T) {
	tS := tStruct{}
	err := serialization.ReadYamlFromString(testS2SSimpleAuthYaml, &tS)
	assert.NoError(t, err)
	authConfig := tS.AuthSimpleS2S
	assert.False(t, authConfig.Disabled)

	endpointCalled := false
	router := gin.New()
	router.GET("/test", authConfig.GinHandlerFuncWithAccessCheck([]string{"read"}, func(c *gin.Context) {
		endpointCalled = true
	}))
	server := httptest.NewServer(router)
	defer server.Close()

	r := resty.New()
	resp, err := r.R().SetHeader("Content-Type", "application/json").
		SetHeader(auth.HeaderClientId, "102").
		SetHeader(auth.HeaderClientAccessToken, "bad_secret").
		Get(fmt.Sprintf("%s/test", server.URL))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	assert.False(t, endpointCalled)
	fmt.Println(resp.String())
}

func TestGoodSecret(t *testing.T) {
	tS := tStruct{}
	err := serialization.ReadYamlFromString(testS2SSimpleAuthYaml, &tS)
	assert.NoError(t, err)
	authConfig := tS.AuthSimpleS2S
	assert.False(t, authConfig.Disabled)

	endpointCalled := false
	router := gin.New()
	router.GET("/test", authConfig.GinHandlerFuncWithAccessCheck([]string{"read"}, func(c *gin.Context) {
		endpointCalled = true
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}))
	server := httptest.NewServer(router)
	defer server.Close()

	r := resty.New()
	resp, err := r.R().SetHeader("Content-Type", "application/json").
		SetHeader(auth.HeaderClientId, "102").
		SetHeader(auth.HeaderClientAccessToken, "user_2_123").
		Get(fmt.Sprintf("%s/test", server.URL))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.True(t, endpointCalled)
	fmt.Println(resp.String())
}

func TestGoodSecretSecond(t *testing.T) {
	tS := tStruct{}
	err := serialization.ReadYamlFromString(testS2SSimpleAuthYaml, &tS)
	assert.NoError(t, err)
	authConfig := tS.AuthSimpleS2S
	assert.False(t, authConfig.Disabled)

	endpointCalled := false
	router := gin.New()
	router.GET("/test", authConfig.GinHandlerFuncWithAccessCheck([]string{"read"}, func(c *gin.Context) {
		endpointCalled = true
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}))
	server := httptest.NewServer(router)
	defer server.Close()

	r := resty.New()
	resp, err := r.R().SetHeader("Content-Type", "application/json").
		SetHeader(auth.HeaderClientId, "102").
		SetHeader(auth.HeaderClientAccessToken, "user_2_abc").
		Get(fmt.Sprintf("%s/test", server.URL))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.True(t, endpointCalled)
	fmt.Println(resp.String())
}

func TestGoodSecretButBadPermission(t *testing.T) {
	tS := tStruct{}
	err := serialization.ReadYamlFromString(testS2SSimpleAuthYaml, &tS)
	assert.NoError(t, err)
	authConfig := tS.AuthSimpleS2S
	assert.False(t, authConfig.Disabled)

	endpointCalled := false
	router := gin.New()
	router.GET("/test", authConfig.GinHandlerFuncWithAccessCheck([]string{"write"}, func(c *gin.Context) {
		endpointCalled = true
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}))
	server := httptest.NewServer(router)
	defer server.Close()

	r := resty.New()
	resp, err := r.R().SetHeader("Content-Type", "application/json").
		SetHeader(auth.HeaderClientId, "102").
		SetHeader(auth.HeaderClientAccessToken, "user_2_123").
		Get(fmt.Sprintf("%s/test", server.URL))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	assert.False(t, endpointCalled)
	fmt.Println(resp.String())
}

func TestGoodSecretGoodPermission(t *testing.T) {
	tS := tStruct{}
	err := serialization.ReadYamlFromString(testS2SSimpleAuthYaml, &tS)
	assert.NoError(t, err)
	authConfig := tS.AuthSimpleS2S
	assert.False(t, authConfig.Disabled)

	endpointCalled := false
	router := gin.New()
	router.GET("/test", authConfig.GinHandlerFuncWithAccessCheck([]string{"write"}, func(c *gin.Context) {
		endpointCalled = true
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}))
	server := httptest.NewServer(router)
	defer server.Close()

	r := resty.New()
	resp, err := r.R().SetHeader("Content-Type", "application/json").
		SetHeader(auth.HeaderClientId, "101").
		SetHeader(auth.HeaderClientAccessToken, "user_1_123").
		Get(fmt.Sprintf("%s/test", server.URL))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.True(t, endpointCalled)
	fmt.Println(resp.String())
}

func TestGoodSecretWithStarPermission(t *testing.T) {
	tS := tStruct{}
	err := serialization.ReadYamlFromString(testS2SSimpleAuthYaml, &tS)
	assert.NoError(t, err)
	authConfig := tS.AuthSimpleS2S
	assert.False(t, authConfig.Disabled)

	endpointCalled := false
	router := gin.New()
	router.GET("/test", authConfig.GinHandlerFuncWithAccessCheck([]string{"write"}, func(c *gin.Context) {
		endpointCalled = true
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}))
	server := httptest.NewServer(router)
	defer server.Close()

	r := resty.New()
	resp, err := r.R().SetHeader("Content-Type", "application/json").
		SetHeader(auth.HeaderClientId, "103").
		SetHeader(auth.HeaderClientAccessToken, "user_3_123").
		Get(fmt.Sprintf("%s/test", server.URL))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.True(t, endpointCalled)
	fmt.Println(resp.String())
}

func TestGoodSecretAdmin(t *testing.T) {
	tS := tStruct{}
	err := serialization.ReadYamlFromString(testS2SSimpleAuthYaml, &tS)
	assert.NoError(t, err)
	authConfig := tS.AuthSimpleS2S
	assert.False(t, authConfig.Disabled)

	endpointCalled := false
	router := gin.New()
	router.GET("/test", authConfig.GinHandlerFuncWithAccessCheck([]string{"and_random"}, func(c *gin.Context) {
		endpointCalled = true
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}))
	server := httptest.NewServer(router)
	defer server.Close()

	r := resty.New()
	resp, err := r.R().SetHeader("Content-Type", "application/json").
		SetHeader(auth.HeaderClientId, "0").
		SetHeader(auth.HeaderClientAccessToken, "admin_123").
		Get(fmt.Sprintf("%s/test", server.URL))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.True(t, endpointCalled)
	fmt.Println(resp.String())
}

func TestAdhocUser(t *testing.T) {
	tS := tStruct{}
	err := serialization.ReadYamlFromString(testS2SSimpleAuthYaml, &tS)
	assert.NoError(t, err)
	authConfig := tS.AuthSimpleS2S
	assert.False(t, authConfig.Disabled)

	// Good case
	t.Run("GoodCase", func(t *testing.T) {
		endpointCalled := false
		router := gin.New()
		router.GET("/test", authConfig.GinHandlerFuncWithAccessCheck([]string{"some_random"}, func(c *gin.Context) {
			endpointCalled = true
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}))
		server := httptest.NewServer(router)
		defer server.Close()

		token := GenerateAdhocUserSecret(time.Duration(1)*time.Minute, "harish", authConfig.SecretSalt)
		r := resty.New()
		resp, err := r.R().SetHeader("Content-Type", "application/json").
			SetHeader(auth.HeaderClientId, "1").
			SetHeader(auth.HeaderClientAccessToken, token).
			Get(fmt.Sprintf("%s/test", server.URL))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		assert.True(t, endpointCalled)
		fmt.Println(resp.String())
	})

	// Expired
	t.Run("Expired", func(t *testing.T) {
		endpointCalled := false
		router := gin.New()
		router.GET("/test", authConfig.GinHandlerFuncWithAccessCheck([]string{"some_random"}, func(c *gin.Context) {
			endpointCalled = true
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}))
		server := httptest.NewServer(router)
		defer server.Close()

		token := GenerateAdhocUserSecret(time.Duration(-10)*time.Minute, "harish", authConfig.SecretSalt)
		r := resty.New()
		resp, err := r.R().SetHeader("Content-Type", "application/json").
			SetHeader(auth.HeaderClientId, "1").
			SetHeader(auth.HeaderClientAccessToken, token).
			Get(fmt.Sprintf("%s/test", server.URL))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
		assert.False(t, endpointCalled)
		fmt.Println(resp.String())

		authError := goxJsonUtils.BytesToObjectSuppressError[*AuthError](resp.Body())
		assert.NotNil(t, authError)
		assert.Equal(t, "adhoc_user_secret_expired", authError.Reason)
		assert.Equal(t, "1", authError.ClientId)
	})
}
