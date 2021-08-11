package lock

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sirupsen/logrus"
)

func CheckDuplicatedRequest(requestID string) bool {
	svc := dynamodb.New(session.New())
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"key": {S: aws.String(requestID)},
			"ttl": {N: aws.String(strconv.FormatInt(time.Now().Add(3*time.Hour).Unix(), 10))},
		},
		TableName:                aws.String("chimera-lock-request"),
		ConditionExpression:      aws.String("attribute_not_exists(#r)"),
		ExpressionAttributeNames: map[string]*string{"#r": aws.String("key")},
	}

	if _, err := svc.PutItem(input); err != nil {
		logrus.WithError(err).Info("try to lock failed")
		return true
	}
	return false
}
