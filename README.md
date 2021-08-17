# Chimera

## Architect

![components](https://raw.githubusercontent.com/liminaab/chimera/main/nocode/system-Components.png)

| Components | Description |
|--|--|
| Slack App Bot <-> Slack Server| Chimera communicate with client through Slack bot chat. It's handled by Slack Client and Slack Server system |
| Slack Server <-> AWS Lambda | Slack server's set the event to call to AWS Lambda (through AWS gateway) to trigger event everytime slack chat bot receive any command/chat |
| AWS Lambda <-> AWS EKS | AWS Lambda connect and control AWS EKS/k8s through AWS API and K8s API |

## Infrastructure setup

### Slack

### AWS

### Build application

### Deploy

## Environments

--- 

## Development

### Requirement

 - Install SAM CLI: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html
 - Go 1.16
 - Setup local AWS Credential: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-getting-started-set-up-credentials.html

### Local

1. Build & Push image to AWS

```
docker build -t xxxxx.dkr.ecr.eu-north-1.amazonaws.com/chimera:test .;
aws ecr get-login-password --region eu-north-1 | docker login --username AWS --password-stdin xxxxx.dkr.ecr.eu-north-1.amazonaws.com;
docker push xxxxx.dkr.ecr.eu-north-1.amazonaws.com/chimera:test;
AWS_PAGER="" aws lambda update-function-code --function-name Chimera --image-uri "xxxxx.dkr.ecr.eu-north-1.amazonaws.com/chimera:test" --no-paginate;
```

### Environments

```
SLACK_VERIFICATION_TOKEN
AUTHS
SLACK_OAUTH_TOKEN
```
