package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/quimnieto/akeneo-serverless-connector/aws/internal/handler"
	"github.com/quimnieto/akeneo-serverless-connector/aws/internal/logger"
	"github.com/quimnieto/akeneo-serverless-connector/aws/internal/processor"
	"github.com/quimnieto/akeneo-serverless-connector/aws/internal/publisher"
)

func main() {
	topicARN := os.Getenv("SNS_TOPIC_ARN")
	if topicARN == "" {
		log.Fatal("SNS_TOPIC_ARN is not set")
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	snsClient := sns.NewFromConfig(cfg)

	lambdaHandler := handler.NewLambdaHandler(
		processor.NewProcessor(),
		publisher.NewSNSPublisher(snsClient, topicARN),
		logger.New(),
	)

	lambda.Start(lambdaHandler.Handle)
}
