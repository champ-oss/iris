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
const retryDelaySeconds = 5
const retryAttempts = 36

func TestIris(t *testing.T) {

	terraformOptions := &terraform.Options{
		TerraformDir:  "../examples/complete",
		BackendConfig: map[string]interface{}{},
		Vars: map[string]interface{}{
			"docker_tag": os.Getenv("GITHUB_SHA"),
		},
	}
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	dns := terraform.Output(t, terraformOptions, "dns")

	t.Log("testing successful request to upstream google")
	err := checkHttpStatusAndBody(t, dns, googleUrl, "OK", http.StatusOK)
	assert.NoError(t, err)

	t.Log("testing successful request to upstream amazon")
	err = checkHttpStatusAndBody(t, dns, amazonUrl, "OK", http.StatusOK)
	assert.NoError(t, err)

	t.Log("testing failed request to not allowed url")
	err = checkHttpStatusAndBody(t, dns, notAllowedUrl, "not allowed", http.StatusForbidden)
	assert.NoError(t, err)

	t.Log("testing failed request to unreachable url")
	err = checkHttpStatusAndBody(t, dns, badUrl, "Internal Server Error", http.StatusInternalServerError)
	assert.NoError(t, err)
}

func checkHttpStatusAndBody(t *testing.T, dns, upstreamUrl, expectedBody string, expectedHttpStatus int) error {
	url := fmt.Sprintf("https://%s/%s", dns, upstreamUrl)
	t.Logf("checking %s", url)

	for i := 0; ; i++ {
		resp, err := http.Get(url)
		if err != nil {
			t.Log(err)
		} else {
			t.Logf("StatusCode: %d", resp.StatusCode)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Log(err)
			} else {
				t.Logf("body: %s", body)
				if resp.StatusCode == expectedHttpStatus && string(body) == expectedBody {
					return nil
				}
			}
		}

		if i >= (retryAttempts - 1) {
			panic("Timed out while retrying")
		}

		t.Logf("Retrying in %d seconds...", retryDelaySeconds)
		time.Sleep(time.Second * retryDelaySeconds)
	}
}
