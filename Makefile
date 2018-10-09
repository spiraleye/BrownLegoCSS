.PHONY: check
check: lint vet test

.PHONY: lint
lint:
	golint -set_exit_status ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test ./...

.PHONY: coverage
coverage:
	go test -cover ./...
