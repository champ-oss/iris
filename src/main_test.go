package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"os"
	"testing"
)

func Test_getAllowedURLs(t *testing.T) {
	key := "TEST_KEY"
	_ = os.Setenv(key, "test1,test2")
	results := getAllowedURLs(key)
	_, test1Present := results["test1"]
	assert.True(t, test1Present)
	_, test2Present := results["test2"]
	assert.True(t, test2Present)
}

func Test_checkIfAllowedURL(t *testing.T) {
	allowedURLs := map[string]struct{}{
		"www.google.com/foo":   {},
		"www.facebook.com/bar": {},
		"www.outlook.com/tee":  {},
	}
	assert.True(t, isAllowedURL("www.google.com/foo", allowedURLs))
	assert.True(t, isAllowedURL("www.outlook.com/tee", allowedURLs))

	assert.False(t, isAllowedURL("www.example.com/blah", allowedURLs))
	assert.False(t, isAllowedURL("www.outlook.com/tee/1", allowedURLs))
}

func Test_httpGetReturnStatusCode(t *testing.T) {
	t.Run("upstream 200 response", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://foo.com").Get("/bar").Reply(200)
		code := httpGetReturnStatusCode("foo.com/bar")
		assert.Equal(t, 200, code)
		assert.True(t, gock.IsDone())
	})

	t.Run("upstream 400 response", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://foo.com").Get("/bar").Reply(400)
		code := httpGetReturnStatusCode("foo.com/bar")
		assert.Equal(t, 400, code)
		assert.True(t, gock.IsDone())
	})

	t.Run("internal http error", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://foo.com").Get("/bar").Reply(200)
		code := httpGetReturnStatusCode("none.com")
		assert.Equal(t, 500, code)
	})
}

func TestHandleRequest(t *testing.T) {
	t.Run("upstream 200 response from GET", func(t *testing.T) {
		envKey := "ALLOWED_URLS"
		_ = os.Setenv(envKey, "www.google.com/foo,www.facebook.com/bar")
		defer gock.Off()
		gock.New("https://www.facebook.com").Get("/bar").Reply(200)
		resp, err := HandleRequest(context.Background(), Event{
			QueryStringParameters: map[string]string{"url": "www.facebook.com/bar"},
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "OK", resp.StatusDescription)
		assert.Equal(t, "OK", resp.Body)
		assert.True(t, gock.IsDone())
	})

	t.Run("upstream 200 response from HEAD", func(t *testing.T) {
		envKey := "ALLOWED_URLS"
		_ = os.Setenv(envKey, "www.google.com/foo,www.facebook.com/bar")
		defer gock.Off()
		gock.New("https://www.facebook.com").Get("/bar").Reply(200)
		resp, err := HandleRequest(context.Background(), Event{
			QueryStringParameters: map[string]string{"url": "www.facebook.com/bar"},
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "OK", resp.StatusDescription)
		assert.Equal(t, "OK", resp.Body)
		assert.True(t, gock.IsDone())
	})

	t.Run("upstream 400 response", func(t *testing.T) {
		envKey := "ALLOWED_URLS"
		_ = os.Setenv(envKey, "www.google.com/foo,www.facebook.com/bar")
		defer gock.Off()
		gock.New("https://www.facebook.com").Get("/bar").Reply(400)
		resp, err := HandleRequest(context.Background(), Event{
			QueryStringParameters: map[string]string{"url": "www.facebook.com/bar"},
		})
		assert.Nil(t, err)
		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, "Bad Request", resp.StatusDescription)
		assert.Equal(t, "Bad Request", resp.Body)
		assert.True(t, gock.IsDone())
	})

	t.Run("internal http error", func(t *testing.T) {
		envKey := "ALLOWED_URLS"
		_ = os.Setenv(envKey, "www.google.com/foo,www.facebook.com/bar")
		defer gock.Off()
		gock.New("https://foo.com").Get("/bar").Reply(200)
		resp, err := HandleRequest(context.Background(), Event{
			QueryStringParameters: map[string]string{"url": "www.facebook.com/bar"},
		})
		assert.Nil(t, err)
		assert.Equal(t, 500, resp.StatusCode)
		assert.Equal(t, "Internal Server Error", resp.StatusDescription)
		assert.Equal(t, "Internal Server Error", resp.Body)
	})

	t.Run("not an allowed url", func(t *testing.T) {
		envKey := "ALLOWED_URLS"
		_ = os.Setenv(envKey, "www.google.com/foo,www.facebook.com/bar")
		resp, err := HandleRequest(context.Background(), Event{
			QueryStringParameters: map[string]string{"url": "www.foo.com/bar"},
		})
		assert.Nil(t, err)
		assert.Equal(t, 403, resp.StatusCode)
		assert.Equal(t, "Forbidden", resp.StatusDescription)
		assert.Equal(t, "not allowed", resp.Body)
	})

	t.Run("url not set", func(t *testing.T) {
		envKey := "ALLOWED_URLS"
		_ = os.Setenv(envKey, "www.google.com/foo,www.facebook.com/bar")
		resp, err := HandleRequest(context.Background(), Event{})
		assert.Nil(t, err)
		assert.Equal(t, 403, resp.StatusCode)
		assert.Equal(t, "Forbidden", resp.StatusDescription)
		assert.Equal(t, "not allowed", resp.Body)
	})
}

func Test_logRequest(t *testing.T) {
	logRequest(context.Background(), Event{
		Url: "foo",
	})
}
