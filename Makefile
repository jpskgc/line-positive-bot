.PHONY: deps clean build

deps:
	# go get -u ./...
	go get -u github.com/aws/aws-lambda-go/events
	go get -u github.com/aws/aws-lambda-go/lambda

clean:
	#rm -rf ./hello-world/hello-world
	rm -rf ./positive-line-bot/positive-line-bot

build:
	#GOOS=linux GOARCH=amd64 go build -o hello-world/hello-world ./hello-world
	GOOS=linux GOARCH=amd64 go build -o positive-line-bot/positive-line-bot ./positive-line-bot
