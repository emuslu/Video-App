.PHONY: dev
dev:
	@go get -u github.com/cosmtrek/air
	@air -c .air.toml
alias air='$(go env GOPATH)/bin/air'