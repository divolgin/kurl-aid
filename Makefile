.PHONY: fmt
fmt:
	go fmt ./pkg/...

.PHONY: vet
vet:
	go vet ./pkg/...

.PHONY: test
test:
	go test ./pkg/... -coverprofile cover.out

.PHONY: build
build: fmt vet test
	go build -o bin/kurl-aid main.go
