build:
	docker build -t qeek-dev/casbin-server .

test:
	go test -v -race -tags test $(shell go list ./... | grep -v 'vendor\|cmd')