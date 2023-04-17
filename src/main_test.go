package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"os"
	"testing"
)

func Test_HandleRequest_WithSingleReachableAllowedUrl_ShouldReturnHttp200(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "www.facebook.com:1234/bar")

	// Mock the expected http reply
	defer gock.Off()
	gock.New("https://www.facebook.com:1234").Get("/bar").Reply(200)

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"url": "www.facebook.com:1234/bar"}}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        200,
		StatusDescription: "OK",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "OK",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithMultipleReachableAllowedUrl_ShouldReturnHttp200(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "www.google.com/foo,www.facebook.com/bar,yahoo.com")

	// Mock the expected http reply
	defer gock.Off()
	gock.New("https://www.facebook.com").Get("/bar").Reply(200)

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"url": "www.facebook.com/bar"}}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        200,
		StatusDescription: "OK",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "OK",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithUnreachableUpstreamUrl_ShouldReturnHttp500(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "localhost:60578")

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"url": "localhost:60578"}}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        500,
		StatusDescription: "Internal Server Error",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "Internal Server Error",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithMissingUpstreamUrl_ShouldReturnHttp403(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "www.google.com/foo,www.facebook.com/bar,yahoo.com")

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"url": ""}}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        403,
		StatusDescription: "Forbidden",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "Forbidden",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithInvalidUpstreamUrl_ShouldReturnHttp403(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "www.google.com/foo,www.facebook.com/bar,yahoo.com")

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"url": "www.google.com/fo"}}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        403,
		StatusDescription: "Forbidden",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "Forbidden",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithUpperCaseUpstreamUrl_ShouldReturnHttp403(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "www.facebook.com:1234/bar")

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"URL": "www.facebook.com:1234/bar"}}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        403,
		StatusDescription: "Forbidden",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "Forbidden",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithEmptyAllowedUrls_ShouldReturnHttp403(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "")

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"url": "www.facebook.com:1234/bar"}}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        403,
		StatusDescription: "Forbidden",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "Forbidden",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithUpstreamUrl400Response_ShouldReturnHttp400(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "www.facebook.com:1234/bar")

	// Mock the expected http reply
	defer gock.Off()
	gock.New("https://www.facebook.com:1234").Get("/bar").Reply(400)

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"url": "www.facebook.com:1234/bar"}}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        400,
		StatusDescription: "Bad Request",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "Bad Request",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithUpstreamUrl500Response_ShouldReturnHttp500(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "www.facebook.com:1234/bar")

	// Mock the expected http reply
	defer gock.Off()
	gock.New("https://www.facebook.com:1234").Get("/bar").Reply(500)

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"url": "www.facebook.com:1234/bar"}}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        500,
		StatusDescription: "Internal Server Error",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "Internal Server Error",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithEmptyExpectedHeaderKey_ShouldReturnHttp200(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "www.facebook.com:1234/bar")
	_ = os.Setenv("EXPECTED_HEADER_KEY", "")
	_ = os.Setenv("EXPECTED_HEADER_VALUE", "")

	// Mock the expected http reply
	defer gock.Off()
	gock.New("https://www.facebook.com:1234").Get("/bar").Reply(200)

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"url": "www.facebook.com:1234/bar"}}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        200,
		StatusDescription: "OK",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "OK",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithEmptyExpectedHeaderValue_ShouldReturnHttp200(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "www.facebook.com:1234/bar")
	_ = os.Setenv("EXPECTED_HEADER_KEY", "TEST-HEADER")
	_ = os.Setenv("EXPECTED_HEADER_VALUE", "")

	// Mock the expected http reply
	defer gock.Off()
	gock.New("https://www.facebook.com:1234").Get("/bar").Reply(200)

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"url": "www.facebook.com:1234/bar"}}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        200,
		StatusDescription: "OK",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "OK",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithEmptyExpectedHeaderKeyAndSetExpectedValue_ShouldReturnHttp403(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "www.facebook.com:1234/bar")
	_ = os.Setenv("EXPECTED_HEADER_KEY", "")
	_ = os.Setenv("EXPECTED_HEADER_VALUE", "test-value-123")

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"url": "www.facebook.com:1234/bar"}}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        403,
		StatusDescription: "Forbidden",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "Forbidden",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithMissingHeaderKey_ShouldReturnHttp403(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "www.facebook.com:1234/bar")
	_ = os.Setenv("EXPECTED_HEADER_KEY", "TEST-HEADER")
	_ = os.Setenv("EXPECTED_HEADER_VALUE", "test-value-123")

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"url": "www.facebook.com:1234/bar"},
		Headers:               map[string]string{"": "test-value-123"},
	}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        403,
		StatusDescription: "Forbidden",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "Forbidden",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithEmptyHeaderValue_ShouldReturnHttp403(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "www.facebook.com:1234/bar")
	_ = os.Setenv("EXPECTED_HEADER_KEY", "TEST-HEADER")
	_ = os.Setenv("EXPECTED_HEADER_VALUE", "test-value-123")

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"url": "www.facebook.com:1234/bar"},
		Headers:               map[string]string{"TEST-HEADER": ""},
	}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        403,
		StatusDescription: "Forbidden",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "Forbidden",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithInvalidHeaderValue_ShouldReturnHttp403(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "www.facebook.com:1234/bar")
	_ = os.Setenv("EXPECTED_HEADER_KEY", "TEST-HEADER")
	_ = os.Setenv("EXPECTED_HEADER_VALUE", "test-value-123")

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"url": "www.facebook.com:1234/bar"},
		Headers:               map[string]string{"TEST-HEADER": "foo"},
	}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        403,
		StatusDescription: "Forbidden",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "Forbidden",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithLowerCaseHeaderKey_ShouldReturnHttp403(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "www.facebook.com:1234/bar")
	_ = os.Setenv("EXPECTED_HEADER_KEY", "TEST-HEADER")
	_ = os.Setenv("EXPECTED_HEADER_VALUE", "test-value-123")

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"url": "www.facebook.com:1234/bar"},
		Headers:               map[string]string{"test-header": "test-value-123"},
	}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        403,
		StatusDescription: "Forbidden",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "Forbidden",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}

func Test_HandleRequest_WithCorrectHeader_ShouldReturnHttp200(t *testing.T) {
	_ = os.Setenv("ALLOWED_URLS", "www.facebook.com:1234/bar")
	_ = os.Setenv("EXPECTED_HEADER_KEY", "TEST-HEADER")
	_ = os.Setenv("EXPECTED_HEADER_VALUE", "test-value-123")

	// Mock the expected http reply
	defer gock.Off()
	gock.New("https://www.facebook.com:1234").Get("/bar").Reply(200)

	// Run the Lambda function
	LoadSettings()
	request := events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"url": "www.facebook.com:1234/bar"},
		Headers:               map[string]string{"TEST-HEADER": "test-value-123"},
	}
	response, err := HandleRequest(context.TODO(), request)

	// Validate the response
	assert.Equal(t, &Response{
		StatusCode:        200,
		StatusDescription: "OK",
		Headers: &ResponseHeaders{
			ContentType: "text/plain",
		},
		Body: "OK",
	}, response)
	assert.Nil(t, err)
	assert.True(t, gock.IsDone())
}
