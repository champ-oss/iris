package test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

const retryDelaySeconds = 5
const retryAttempts = 36

func TestIris(t *testing.T) {

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/complete",
		BackendConfig: map[string]interface{}{
			"bucket": os.Getenv("TF_STATE_BUCKET"),
			"key":    os.Getenv("TF_VAR_git"),
		},
		Vars: map[string]interface{}{},
	}
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	functionUrl := terraform.Output(t, terraformOptions, "function_url")                   // Lambda that requires a special header
	functionUrlNoHeader := terraform.Output(t, terraformOptions, "function_url_no_header") // Lambda that doesnt require a special header
	headerKey := terraform.Output(t, terraformOptions, "expected_header_key")
	headerVal := terraform.Output(t, terraformOptions, "expected_header_value")

	t.Log("testing successful request to upstream google")
	err := checkHttpStatusAndBody(t, fmt.Sprintf("%s?url=%s", functionUrl, "about.google/google-in-america"), headerKey, headerVal, "OK", http.StatusOK)
	assert.NoError(t, err)

	t.Log("testing successful request to upstream amazon")
	err = checkHttpStatusAndBody(t, fmt.Sprintf("%s?url=%s", functionUrl, "aws.amazon.com/console"), headerKey, headerVal, "OK", http.StatusOK)
	assert.NoError(t, err)

	t.Log("testing failed request to not allowed url")
	err = checkHttpStatusAndBody(t, fmt.Sprintf("%s?url=%s", functionUrl, "www.example.com/foo/bar"), headerKey, headerVal, "Forbidden", http.StatusForbidden)
	assert.NoError(t, err)

	t.Log("testing failed request to unreachable url")
	err = checkHttpStatusAndBody(t, fmt.Sprintf("%s?url=%s", functionUrl, "1.com/foo"), headerKey, headerVal, "Internal Server Error", http.StatusInternalServerError)
	assert.NoError(t, err)

	t.Log("testing failed request to empty url")
	err = checkHttpStatusAndBody(t, fmt.Sprintf("%s?url=%s", functionUrl, ""), headerKey, headerVal, "Forbidden", http.StatusForbidden)
	assert.NoError(t, err)

	t.Log("testing failed request to missing url key")
	err = checkHttpStatusAndBody(t, functionUrl, headerKey, headerVal, "Forbidden", http.StatusForbidden)
	assert.NoError(t, err)

	t.Log("testing failed request with missing header")
	err = checkHttpStatusAndBody(t, functionUrl, "foo", "", "Forbidden", http.StatusForbidden)
	assert.NoError(t, err)

	t.Log("testing failed request with invalid header")
	err = checkHttpStatusAndBody(t, functionUrl, headerKey, "", "Forbidden", http.StatusForbidden)
	assert.NoError(t, err)

	t.Log("testing successful request to upstream amazon against the separate lambda with no header required")
	err = checkHttpStatusAndBody(t, fmt.Sprintf("%s?url=%s", functionUrlNoHeader, "aws.amazon.com/console"), "foo", "", "OK", http.StatusOK)
	assert.NoError(t, err)
}

func checkHttpStatusAndBody(t *testing.T, url, headerKey, headerVal, expectedBody string, expectedHttpStatus int) error {
	t.Logf("checking %s", url)
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set(headerKey, headerVal)

	for i := 0; ; i++ {
		resp, err := client.Do(request)
		if err != nil {
			t.Log(err)
		} else {
			t.Logf("StatusCode: %d", resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
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
