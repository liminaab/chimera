package k8s

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	clientset "k8s.io/client-go/kubernetes"
)

func IsValidNameSpace(namespace string) bool {
	if len(os.Getenv("CLUSTERS")) > 0 {
		configs, err := ParseAuthsFromEnv()
		if err != nil {
			logrus.Warn("ParseAuthsFromEnv failed", err)
			return false
		}

		_, err = configs.GetByName(namespace)
		if err != nil {
			logrus.Warn("configs.GetByName failed", err)
		}
		return err == nil
	}

	_, _, err := GetEKSConfig(namespace)
	if err != nil {
		logrus.Warn("GetEKSConfig failed", err)
		return false
	}
	return true
}

// NewAuthClient creates a new EKS authenticated clientset.
func NewAuthClient(config *ClusterConfig) (*clientset.Clientset, string, error) {
	// Start new AWS session if not specified
	if config.Session == nil {
		config.Session = newSession()
	}

	// Load the rest from AWS using SDK
	err := config.loadConfig()
	if err != nil {
		return nil, "", errors.Wrap(err, "Unable to load Kubernetes Client Config")
	}

	// Create the Kubernetes client
	client, err := config.NewClientConfig()
	if err != nil {
		return nil, "", errors.Wrap(err, "Unable to create Kubernetes Client Config")
	}

	clientset, token, err := client.NewClientSetWithEmbeddedToken()
	if err != nil {
		return nil, "", errors.Wrap(err, "Unable to create Kubernetes Client Set")
	}

	return clientset, token, nil
}

// Load k8s cluster config from cache.
func (c *ClusterConfig) loadConfig() error {
	if c.ClusterName == "" {
		return errors.New("ClusterName cannot be empty")
	}

	if os.Getenv("CLUSTERS") != "" {
		configs, err := ParseAuthsFromEnv()
		if err != nil {
			return errors.Errorf("Can't load error %v", err)
		}
		config, err := configs.GetByName(c.ClusterName)
		if err != nil {
			return errors.Errorf("ClusterName %s is not cached. Contact admin to check", c.ClusterName)
		}
		c.MasterEndpoint = config.Endpoint
		c.CertificateAuthorityData = config.CertificateAuthority
	} else {
		endpoint, certificateAuthorityData, err := GetEKSConfig(c.ClusterName)
		if err != nil {
			return errors.Errorf("Error to get EKS config %s", err)
		}
		c.MasterEndpoint = endpoint
		c.CertificateAuthorityData = certificateAuthorityData
	}

	return nil
}

// Real check cluster config
func GetEKSConfig(clusterName string) (string, string, error) {
	if clusterName == "" {
		errors.New("ClusterName cannot be empty")
	}

	ss := newSession()
	svc := eks.New(ss)
	input := &eks.DescribeClusterInput{Name: aws.String(clusterName)}
	logrus.WithField("cluster", clusterName).Info(time.Now(), "Looking up EKS cluster")

	result, err := svc.DescribeCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			logrus.WithField("cluster", clusterName).Error(aerr.Error())
			return "", "", errors.Wrap(err, aerr.Error())
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			logrus.WithField("cluster", clusterName).Error(err.Error())
			return "", "", errors.Wrap(err, err.Error())
		}
	}

	logrus.WithFields(logrus.Fields{
		"cluster":               clusterName,
		"endpoint":              *result.Cluster.Endpoint,
		"certificate_authority": *result.Cluster.CertificateAuthority.Data,
	}).Info(time.Now(), "Found cluster")
	return *result.Cluster.Endpoint, *result.Cluster.CertificateAuthority.Data, nil
}

func (c *ClusterConfig) NewClientConfig() (*ClientConfig, error) {

	stsAPI := sts.New(c.Session)

	iamRoleARN, err := checkAuth(stsAPI)
	if err != nil {
		return nil, err
	}
	contextName := fmt.Sprintf("%s@%s", getUsername(iamRoleARN), c.ClusterName)

	data, err := base64.StdEncoding.DecodeString(c.CertificateAuthorityData)
	if err != nil {
		return nil, errors.Wrap(err, "decoding certificate authority data")
	}

	logrus.Info("Creating Kubernetes client config")
	clientConfig := &ClientConfig{
		Client: &clientcmdapi.Config{
			Clusters: map[string]*clientcmdapi.Cluster{
				c.ClusterName: {
					Server:                   c.MasterEndpoint,
					CertificateAuthorityData: data,
				},
			},
			Contexts: map[string]*clientcmdapi.Context{
				contextName: {
					Cluster:  c.ClusterName,
					AuthInfo: contextName,
				},
			},
			AuthInfos: map[string]*clientcmdapi.AuthInfo{
				contextName: &clientcmdapi.AuthInfo{},
			},
			CurrentContext: contextName,
		},
		ClusterName: c.ClusterName,
		ContextName: contextName,
		roleARN:     iamRoleARN,
		sts:         stsAPI,
	}

	return clientConfig, nil

}

func newSession() *session.Session {
	config := aws.NewConfig()
	config = config.WithCredentialsChainVerboseErrors(true)

	opts := session.Options{
		Config:                  *config,
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	}

	stscreds.DefaultDuration = 30 * time.Minute

	return session.Must(session.NewSessionWithOptions(opts))
}

func checkAuth(stsAPI stsiface.STSAPI) (string, error) {
	input := &sts.GetCallerIdentityInput{}
	output, err := stsAPI.GetCallerIdentity(input)
	if err != nil {
		return "", errors.Wrap(err, "checking AWS STS access – cannot get role ARN for current session")
	}
	iamRoleARN := *output.Arn
	return iamRoleARN, nil
}

type ClusterConfig struct {
	ClusterName              string
	MasterEndpoint           string
	CertificateAuthorityData string
	Session                  *session.Session
}

type ClientConfig struct {
	Client      *clientcmdapi.Config
	ClusterName string
	ContextName string
	roleARN     string
	sts         stsiface.STSAPI
}

func getUsername(iamRoleARN string) string {
	usernameParts := strings.Split(iamRoleARN, "/")
	if len(usernameParts) > 1 {
		return usernameParts[len(usernameParts)-1]
	}
	return "iam-root-account"
}

func (c *ClientConfig) WithEmbeddedToken() (*ClientConfig, string, error) {
	clientConfigCopy := *c

	logrus.Info("Generating token")

	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return nil, "", errors.Wrap(err, "could not get token generator")
	}

	tok, err := gen.GetWithSTS(c.ClusterName, c.sts.(*sts.STS))
	if err != nil {
		return nil, "", errors.Wrap(err, "could not get token")
	}

	x := c.Client.AuthInfos[c.ContextName]
	x.Token = tok.Token
	return &clientConfigCopy, tok.Token, nil
}

func (c *ClientConfig) NewClientSetWithEmbeddedToken() (*clientset.Clientset, string, error) {
	clientConfig, token, err := c.WithEmbeddedToken()
	if err != nil {
		return nil, "", errors.Wrap(err, "creating Kubernetes client config with embedded token")
	}
	clientSet, err := clientConfig.NewClientSet()
	if err != nil {
		return nil, "", errors.Wrap(err, "creating Kubernetes client")
	}
	return clientSet, token, nil
}

func (c *ClientConfig) NewClientSet() (*clientset.Clientset, error) {
	clientConfig, err := clientcmd.NewDefaultClientConfig(*c.Client, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create API client configuration from client config")
	}

	client, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create API client")
	}
	return client, nil
}

func GetK8sRestConfig(clusterName, token string) (*rest.Config, error) {

	var endpoint string
	if os.Getenv("CLUSTERS") != "" {
		configs, err := ParseAuthsFromEnv()
		if err != nil {
			return nil, errors.Errorf("Can't load error %v", err)
		}

		config, err := configs.GetByName(clusterName)
		if err != nil {
			return nil, errors.Errorf("ClusterName %s is not cached. Contact admin to check", clusterName)
		}
		endpoint = config.Endpoint
	} else {
		var err error
		endpoint, _, err = GetEKSConfig(clusterName)
		if err != nil {
			return nil, errors.Errorf("Error to get EKS config %s", err)
		}
	}

	ret := &rest.Config{
		Host:            endpoint,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
		APIPath:         "/",
		BearerToken:     token,
	}

	return ret, nil
}
