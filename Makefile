lint:
	golangci-lint run

test:
	go test -count=1 ./... -v

install-wampproto:
	sudo snap install wampproto --edge
