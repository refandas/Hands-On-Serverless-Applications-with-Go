package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Movie struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func update(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var movie Movie
	err := json.Unmarshal([]byte(request.Body), &movie)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid payload",
		}, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while retrieving AWS credentials",
		}, nil
	}

	svc := dynamodb.NewFromConfig(cfg)

	// Update an item by overwriting it with a new item
	_, err = svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{
				Value: movie.ID,
			},
			"Name": &types.AttributeValueMemberS{
				Value: movie.Name,
			},
		},
	})

	// Update an item by modifying the existing item
	//update := expression.Set(expression.Name("Name"), expression.Value(movie.Name))
	//expr, err := expression.NewBuilder().WithUpdate(update).Build()
	//
	//_, err = svc.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
	//	TableName: aws.String(os.Getenv("TABLE_NAME")),
	//	Key: map[string]types.AttributeValue{
	//		"ID": &types.AttributeValueMemberS{
	//			Value: movie.ID,
	//		},
	//	},
	//	ExpressionAttributeNames:  expr.Names(),
	//	ExpressionAttributeValues: expr.Values(),
	//	UpdateExpression:          expr.Update(),
	//})

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while updating the movie",
		}, nil
	}

	response, err := json.Marshal(movie)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while decoding to string value",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(update)
}
