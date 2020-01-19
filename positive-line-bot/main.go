package main

import (
	"encoding/json"
	"fmt"
	//"log"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/ssm"
	//"github.com/line/line-bot-sdk-go/linebot"
)

type Positive struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

func UnmarshalLineRequest(data []byte) (LineRequest, error) {
	var r LineRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

type LineRequest struct {
	Events      []Event `json:"events"`
	Destination string  `json:"destination"`
}

type Event struct {
	Type       string  `json:"type"`
	ReplyToken string  `json:"replyToken"`
	Source     Source  `json:"source"`
	Timestamp  int64   `json:"timestamp"`
	Message    Message `json:"message"`
}

type Message struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Text string `json:"text"`
}

type Source struct {
	UserID string `json:"userId"`
	Type   string `json:"type"`
}

func getSSMParameterStore(parameter string) string {
	svc := ssm.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})

	res, _ := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(parameter),
		WithDecryption: aws.Bool(true),
	})
	return *res.Parameter.Value
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// myLineRequest, err := UnmarshalLineRequest([]byte(request.Body))
	// if err != nil {
	// 	log.Fatal(err)
	// 	return events.APIGatewayProxyResponse{}, err
	// }

	// bot, err := linebot.New(
	// 	// os.Getenv("LINE_CHANNEL_SECRET"),
	// 	// os.Getenv("LINE_CHANNEL_TOKEN"),
	// 	getSSMParameterStore("LINE_CHANNEL_SECRET"),
	// 	getSSMParameterStore("LINE_CHANNEL_TOKEN"),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// 	return events.APIGatewayProxyResponse{}, err
	// }

	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")

	sess := session.Must(session.NewSession())
	config := aws.NewConfig().WithRegion("ap-northeast-1")
	if len(endpoint) > 0 {
		config = config.WithEndpoint(endpoint)
	}

	db := dynamodb.New(sess, config)

	result, err := db.Scan(&dynamodb.ScanInput{
		TableName:              aws.String(tableName),
		ConsistentRead:         aws.Bool(true),
		ReturnConsumedCapacity: aws.String("NONE"),
	})

	var positives []Positive

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &positives)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{}, err
	}

	var words []string
	for _, positive := range positives {
		words = append(words, positive.Name)
	}

	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(words))
	word := words[i]

	// var tmpReplyMessage string
	// tmpReplyMessage = word
	// if _, err = bot.ReplyMessage(myLineRequest.Events[0].ReplyToken, linebot.NewTextMessage(tmpReplyMessage)).Do(); err != nil {
	// 	log.Fatal(err)
	// 	return events.APIGatewayProxyResponse{}, err
	// }

	return events.APIGatewayProxyResponse{
		Body:       string(word),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
