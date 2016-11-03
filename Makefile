build:
	go build -v

fmt:
	find . -name "*.go" -not -path "./vendor/*" | xargs gofmt -w -s

docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a --ldflags "-s" && upx -q6 showfast
	docker build --rm -t perflab/showfast .

clean:
	rm -fr showfast

test:
	go test -v -cover -race
