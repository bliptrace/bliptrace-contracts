.PHONY: test tidy validate

test:
	go test ./...

validate: test

tidy:
	go mod tidy
