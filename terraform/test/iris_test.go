package test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

// Upstream urls to test with
const googleUrl = "about.google/google-in-america"
const amazonUrl = "aws.amazon.com/console"
const badUrl = "1.com/foo"
const notAllowedUrl = "www.example.com/foo/bar"

// TestIris tests the application in an ephemeral environment
func TestIris(t *testing.T) {

	terraformOptions := &terraform.Options{
		TerraformDir:  "../examples/complete",
		BackendConfig: map[string]interface{}{},
		Vars: map[string]interface{}{
			"commit_sha": os.Getenv("GITHUB_SHA"),
		},
	}
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	dns := terraform.Output(t, terraformOptions, "dns")

	t.Log("testing successful request to upstream google")
	status, body := getHttpStatusAndBody(t, dns, googleUrl)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "OK", body)

	t.Log("testing successful request to upstream amazon")
	status, body = getHttpStatusAndBody(t, dns, amazonUrl)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "OK", body)

	t.Log("testing failed request to not allowed url")
	status, body = getHttpStatusAndBody(t, dns, notAllowedUrl)
	assert.Equal(t, http.StatusForbidden, status)
	assert.Equal(t, "not allowed", body)

	t.Log("testing failed request using POST")
	status, body = postHttpStatusAndBody(t, dns, googleUrl)
	assert.Equal(t, http.StatusForbidden, status)
	assert.Equal(t, "method not allowed", body)

	t.Log("testing failed request to unreachable url")
	status, body = getHttpStatusAndBody(t, dns, badUrl)
	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Equal(t, "Internal Server Error", body)

	t.Log("delaying destroy to allow for troubleshooting time")
	time.Sleep(10 * time.Minute)
}

func getHttpStatusAndBody(t *testing.T, dns, upstreamUrl string) (int, string) {
	resp, err := http.Get(fmt.Sprintf("https://%s/%s", dns, upstreamUrl))
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	return resp.StatusCode, string(body)
}

func postHttpStatusAndBody(t *testing.T, dns, upstreamUrl string) (int, string) {
	resp, err := http.Post(fmt.Sprintf("https://%s/%s", dns, upstreamUrl), "text/plain", nil)
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	return resp.StatusCode, string(body)
}
