package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"

	"github.com/monorepo/companion_services/chimera/commands"
	"github.com/monorepo/companion_services/chimera/lock"
	"github.com/monorepo/companion_services/chimera/slackutils"
)

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log info request
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.WithFields(logrus.Fields{
		"env":               os.Environ(),
		"request_id":        request.RequestContext.RequestID,
		"method":            request.HTTPMethod,
		"path":              request.Path,
		"res":               request.Resource,
		"body":              request.Body,
		"is_base64_encoded": request.IsBase64Encoded,
	}).Info("New request")

	// Parse request
	body := request.Body
	generalRequest, err := slackutils.ParseGeneralRequest([]byte(body))
	if err != nil {
		logrus.WithError(err).Info("Error")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent}, nil
	}
	logrus.WithField("data", generalRequest).Info("GeneralRequest")

	if !generalRequest.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		logrus.WithError(fmt.Errorf("Invalid token")).Info("Error")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized}, nil
	}

	// Challange
	if generalRequest.Type == slackutils.TypeURLVerification {
		challengeRequest, err := slackutils.ParseURLVerificationRequest([]byte(body))
		if err != nil {
			logrus.WithError(err).Info("Error")
		}
		logrus.WithField("data", challengeRequest).Info("ChallengeRequest")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent, Body: challengeRequest.Challenge}, nil
	}

	// handle request
	messageRequest, err := slackutils.ParseMessageEventRequest([]byte(body))
	if err != nil {
		logrus.WithError(err).Info("Error")
		slackutils.SendSimple(slackutils.GetOAuthToken(), messageRequest.Event.Channel, fmt.Sprintf("Your command got error: parsing request %s", err))
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent}, nil
	}

	logrus.WithField("data", messageRequest).Info("MessageRequest")
	if messageRequest.IsMessageFromBot() {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent}, nil
	}

	if lock.CheckDuplicatedRequest(messageRequest.EventID) {
		logrus.Warnf("Duplicate request event ID %s", messageRequest.EventID)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent}, nil
	}

	if err := handleCommands(messageRequest); err != nil {
		logrus.WithError(err).Info("Handle commands")
		slackutils.SendSimple(slackutils.GetOAuthToken(), messageRequest.Event.Channel, fmt.Sprintf("Your command got error: %s", err))
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent}, nil
}

func handleCommands(request *slackutils.MessageEventRequest) error {
	command, args := commands.ParseIntent(request)
	switch command {
	case commands.IntentBackupDB:
		return commands.BackupDB(request, args)
	case commands.IntentListEnvironments:
		return commands.ListEnvironments(request)
	case commands.IntentListClients:
		return commands.ListClients(request, args)
	case commands.IntentInvalidCommand:
		return commands.InvalidCommand(request)
	default:
		return commands.RootHelp(request)
	}
	return nil
}

func main() {
	// k8s.GetEKSConfig("test")
	// k8s.TestNamespaces()
	// fmt.Println(k8s.Exec("dev", "qa", "cerberus", "mkdir /root/scripts/test > /dev/null 2> /dev/null &"))
	lambda.Start(HandleRequest)
}
