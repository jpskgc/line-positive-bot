.PHONY: deps clean build

deps:
	# go get -u ./...
	go get -u github.com/aws/aws-lambda-go/events
	go get -u github.com/aws/aws-lambda-go/lambda
	go get -u github.com/aws/aws-sdk-go/aws
	go get -u github.com/aws/aws-sdk-go/aws/session
	go get -u github.com/aws/aws-sdk-go/service/dynamodb
	go get -u github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute
	go get -u github.com/line/line-bot-sdk-go/linebot
	go get -u github.com/aws/aws-sdk-go/service/ssm

clean:
	rm -rf ./positive-line-bot/positive-line-bot

build:
	GOOS=linux GOARCH=amd64 go build -o positive-line-bot/positive-line-bot ./positive-line-bot
