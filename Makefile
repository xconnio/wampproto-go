lint:
	golangci-lint run

test:
	go test -count=1 ./... -v --race

install-wampproto:
	sudo snap install wampproto --edge
