package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Item struct {
	ID    string `json:"ID"`
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

var db *dynamodb.DynamoDB

func initDB() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		log.Fatalf("Failed to create session, %v", err)
	}
	db = dynamodb.New(sess)

	tableName := "ItemTable"
	err = ensureTableExists(tableName)
	if err != nil {
		log.Fatalf("Failed to ensure table exists, %v", err)
	}
}

func main() {
	initDB()
	lambda.Start(lambdaHandler)
}

func ensureTableExists(tableName string) error {
	_, err := db.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
			_, err := db.CreateTable(&dynamodb.CreateTableInput{
				TableName: aws.String(tableName),
				AttributeDefinitions: []*dynamodb.AttributeDefinition{
					{
						AttributeName: aws.String("ID"),
						AttributeType: aws.String("S"),
					},
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("ID"),
						KeyType:       aws.String("HASH"),
					},
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			})
			if err != nil {
				return err
			}

			err = db.WaitUntilTableExists(&dynamodb.DescribeTableInput{
				TableName: aws.String(tableName),
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func lambdaHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "POST":
		return handleCreateItem(request)
	case "PUT":
		return handleUpdateItem(request)
	case "GET":
		return handleGetItem(request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Method Not Allowed",
		}, nil
	}
}

func handleCreateItem(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var item Item
	err := json.Unmarshal([]byte(request.Body), &item)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("ItemTable"),
		Item:      av,
	}

	_, err = db.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Item created successfully",
	}, nil
}

func handleUpdateItem(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var item Item
	err := json.Unmarshal([]byte(request.Body), &item)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("ItemTable"),
		Item:      av,
	}

	_, err = db.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Item updated successfully",
	}, nil
}

func handleGetItem(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]

	input := &dynamodb.GetItemInput{
		TableName: aws.String("ItemTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
	}

	result, err := db.GetItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	if result.Item == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "Item not found",
		}, nil
	}

	var item Item
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	response, err := json.Marshal(item)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil
}
