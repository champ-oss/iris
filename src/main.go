package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	StatusCode        int              `json:"statusCode"`
	StatusDescription string           `json:"statusDescription"`
	Headers           *ResponseHeaders `json:"headers"`
	Body              string           `json:"body"`
}

type ResponseHeaders struct {
	ContentType string `json:"Content-Type"`
}

var allowedURLs map[string]struct{}
var expectedHeaderKey string
var expectedHeaderValue string

// init sets up the global state for the Lambda function
func init() {
	LoadSettings()
}

// LoadSettings sets log configuration and loads configuration settings
func LoadSettings() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

	// Load comma-separated list of allowed upstream URLs
	allowedURLs = getAllowedURLs("ALLOWED_URLS")

	// Load optional header to check for in each request
	expectedHeaderKey = os.Getenv("EXPECTED_HEADER_KEY")
	expectedHeaderValue = os.Getenv("EXPECTED_HEADER_VALUE")

	if expectedHeaderKey != "" {
		log.Debugf("expected header: %s=%s", expectedHeaderKey, expectedHeaderValue)
	}
}

// HandleRequest is the entry point for lambda events
func HandleRequest(ctx context.Context, event events.LambdaFunctionURLRequest) (*Response, error) {
	logAsJson("context", ctx)
	logAsJson("event", event)

	// Check if the expected header is present in the request
	headerValue := event.Headers[expectedHeaderKey]
	if headerValue != expectedHeaderValue {
		log.Warningf("invalid header value: %s", headerValue)
		log.Debugf("expected header: %s=%s", expectedHeaderKey, expectedHeaderValue)
		return respondForbidden()
	}

	// Check if the requested upstream url is allowed
	upstreamUrl := event.QueryStringParameters["url"]
	if !isAllowedURL(upstreamUrl, allowedURLs) {
		log.Warningf("Requested url is not allowed: %s", upstreamUrl)
		log.Debugf("allowed urls: %v", allowedURLs)
		return respondForbidden()
	}

	// Send the upstream request and pass along the returned status code
	status := httpGetReturnStatusCode(upstreamUrl)
	return &Response{
		StatusCode:        status,
		StatusDescription: http.StatusText(status),
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: http.StatusText(status),
	}, nil
}

// logAsJson logs the object as a JSON string
func logAsJson(name string, object interface{}) {
	data, _ := json.Marshal(object)
	log.Debugf("%s: %s", name, data)
}

// getAllowedURLs parses a comma separated list of allowed URLs from env variable
func getAllowedURLs(envKey string) map[string]struct{} {
	log.Infof("loading allowed URLs from env %s=%s", envKey, os.Getenv(envKey))

	// Use a map of empty structs for efficient lookups (https://yourbasic.org/golang/implement-set/)
	allowedURLs := map[string]struct{}{}

	if envValue := os.Getenv(envKey); envValue != "" {
		for _, url := range strings.Split(envValue, ",") {
			allowedURLs[url] = struct{}{}
		}
	}
	return allowedURLs
}

// isAllowedURL returns true/false if the requested URL should be allowed
func isAllowedURL(path string, allowedURLs map[string]struct{}) bool {
	_, present := allowedURLs[path]
	return present
}

// httpGetReturnStatusCode sends an HTTP request to the requested URL and returns the status code
func httpGetReturnStatusCode(url string) int {
	fullUrl := fmt.Sprintf("https://%s", url)
	log.Infof("sending GET request to: %s", fullUrl)
	resp, err := http.Get(fullUrl)
	if err != nil {
		log.Errorf("error calling %s - we will return http 500 status.", fullUrl)
		log.Error(err)
		return http.StatusInternalServerError
	}
	log.Infof("upstream service responded with %d", resp.StatusCode)
	return resp.StatusCode
}

// Generate a 403 Forbidden Response
func respondForbidden() (*Response, error) {
	return &Response{
		StatusCode:        http.StatusForbidden,
		StatusDescription: http.StatusText(http.StatusForbidden),
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "Forbidden",
	}, nil

}

func main() {
	lambda.Start(HandleRequest)
}
