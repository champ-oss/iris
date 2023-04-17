run:
	cd src && go run cmd/main.go

test:
	cd src && go test ./...

coverage:
	cd src && go test -json -coverprofile=cover.out ./... > result.json
	cd src && go tool cover -func cover.out
	cd src && go tool cover -html=cover.out

fmt:
	cd src && go fmt ./...
	terraform fmt -recursive -diff

tidy:
	cd src && go mod tidy
	cd terraform/test && go mod tidy