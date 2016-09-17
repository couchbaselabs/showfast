build:
	go build -v 

fmt:
	find . -name "*.go" -not -path "./vendor/*" | xargs gofmt -w -s

docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a 
	docker build --rm -t perflab/showfast .

clean:
	rm -fr showfast
