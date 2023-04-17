# Iris

An AWS Lambda function which can proxy healthcheck requests to your internal services

[![.github/workflows/gotest.yml](https://github.com/champ-oss/iris/actions/workflows/gotest.yml/badge.svg?branch=main)](https://github.com/champ-oss/iris/actions/workflows/gotest.yml)
[![.github/workflows/golint.yml](https://github.com/champ-oss/iris/actions/workflows/golint.yml/badge.svg?branch=main)](https://github.com/champ-oss/iris/actions/workflows/golint.yml)
[![.github/workflows/sonar.yml](https://github.com/champ-oss/iris/actions/workflows/sonar.yml/badge.svg)](https://github.com/champ-oss/iris/actions/workflows/sonar.yml)
[![.github/workflows/gotest.yml](https://github.com/champ-oss/iris/actions/workflows/gotest.yml/badge.svg?branch=main)](https://github.com/champ-oss/iris/actions/workflows/gotest.yml)

[![SonarCloud](https://sonarcloud.io/images/project_badges/sonarcloud-black.svg)](https://sonarcloud.io/summary/new_code?id=champ-oss_iris)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=champ-oss_iris&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=champ-oss_iris)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=champ-oss_iris&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=champ-oss_iris)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=champ-oss_iris&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=champ-oss_iris)

## Example Usage

Call the Lambda function endpoint and pass a `url` query parameter with the value of the URL to your internal service.

```shell
curl https://example123.lambda-url.us-east-1.on.aws/?url=myservice.lan/healthcheck
```

The Lambda function will make an HTTP GET request to the `url` provided and return the HTTP status code. 
No request body or headers will be passed to the internal service. And no response body or headers from your internal service will be returned from the lambda.


## Expected Headers
You can (optionally) set the Terraform variables: `expected_header_key` and `expected_header_value` to require the header
to be present on every request. For example:
```terraform
expected_header_key=X-MY-HEADER
expected_header_value=al1s9v8u210vn410vn
```

