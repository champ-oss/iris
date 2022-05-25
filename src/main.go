package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

type Event struct {
	HttpMethod            string            `json:"httpMethod"`
	Path                  string            `json:"path"`
	QueryStringParameters map[string]string `json:"queryStringParameters"`
	Body                  string            `json:"body"`
}

type Response struct {
	IsBase64Encoded   bool             `json:"isBase64Encoded"`
	StatusCode        int              `json:"statusCode"`
	StatusDescription string           `json:"statusDescription"`
	Headers           *ResponseHeaders `json:"headers"`
	Body              string           `json:"body"`
}

type ResponseHeaders struct {
	ContentType string `json:"Content-Type"`
}

// HandleRequest is the entry point for lambda events
func HandleRequest(ctx context.Context, event Event) (*Response, error) {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	logRequest(ctx, event)
	resp := &Response{Headers: &ResponseHeaders{ContentType: "text/plain"}, Body: ""}

	// only allow GET or HEAD requests
	if event.HttpMethod != http.MethodGet && event.HttpMethod != http.MethodHead {
		log.Warningf("method not allowed: %s", event.HttpMethod)
		resp.StatusCode = http.StatusForbidden
		resp.StatusDescription = http.StatusText(http.StatusForbidden)
		resp.Body = "method not allowed"
		return resp, nil
	}

	// Remove the "/" at the beginning of the requested path
	event.Path = strings.TrimPrefix(event.Path, "/")

	// Load comma separated list of allowed upstream URLs
	allowedURLs := getAllowedURLs("ALLOWED_URLS")

	if !isAllowedURL(event.Path, allowedURLs) {
		log.Warningf("Requested url is not allowed: %s", event.Path)
		log.Debugf("allowed urls: %v", allowedURLs)
		resp.StatusCode = http.StatusForbidden
		resp.StatusDescription = http.StatusText(http.StatusForbidden)
		resp.Body = "not allowed"
		return resp, nil
	}

	// Send the upstream request and pass along the returned status code
	status := httpGetReturnStatusCode(event.Path)
	resp.StatusCode = status
	resp.StatusDescription = http.StatusText(status)
	resp.Body = http.StatusText(status)
	return resp, nil
}

// logRequest logs request details from the load balancer
func logRequest(ctx context.Context, event Event) {
	log.Debugf("context: %s", ctx)
	log.Debugf("event: %s", event)
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

func main() {
	lambda.Start(HandleRequest)
}
