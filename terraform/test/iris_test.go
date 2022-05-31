package test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

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
	terraform.Init(t, terraformOptions)

	// recursively set prevent destroy to false
	cmd := exec.Command("bash", "-c", "find . -type f -name '*.tf' -exec sed -i'' -e 's/prevent_destroy = true/prevent_destroy = false/g' {} +")
	cmd.Dir = "../../"
	_ = cmd.Run()

	terraform.ApplyAndIdempotent(t, terraformOptions)

	functionUrl := terraform.Output(t, terraformOptions, "function_url")

	t.Log("testing successful request to upstream google")
	err := checkHttpStatusAndBody(t, fmt.Sprintf("%s?url=%s", functionUrl, "about.google/google-in-america"), "OK", http.StatusOK)
	assert.NoError(t, err)

	t.Log("testing successful request to upstream amazon")
	err = checkHttpStatusAndBody(t, fmt.Sprintf("%s?url=%s", functionUrl, "aws.amazon.com/console"), "OK", http.StatusOK)
	assert.NoError(t, err)

	t.Log("testing failed request to not allowed url")
	err = checkHttpStatusAndBody(t, fmt.Sprintf("%s?url=%s", functionUrl, "www.example.com/foo/bar"), "not allowed", http.StatusForbidden)
	assert.NoError(t, err)

	t.Log("testing failed request to unreachable url")
	err = checkHttpStatusAndBody(t, fmt.Sprintf("%s?url=%s", functionUrl, "1.com/foo"), "Internal Server Error", http.StatusInternalServerError)
	assert.NoError(t, err)

	t.Log("testing failed request to empty url")
	err = checkHttpStatusAndBody(t, fmt.Sprintf("%s?url=%s", functionUrl, ""), "not allowed", http.StatusForbidden)
	assert.NoError(t, err)

	t.Log("testing failed request to missing url key")
	err = checkHttpStatusAndBody(t, functionUrl, "not allowed", http.StatusForbidden)
	assert.NoError(t, err)

	t.Log("testing successful request to upstream amazon with upper case url")
	err = checkHttpStatusAndBody(t, fmt.Sprintf("%s?URL=%s", functionUrl, "aws.amazon.com/console"), "OK", http.StatusOK)
	assert.NoError(t, err)
}

func checkHttpStatusAndBody(t *testing.T, url, expectedBody string, expectedHttpStatus int) error {
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
