package k8s

import (
	"bytes"
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/remotecommand"
)

const ServiceCerberus = "cerberus"

func Exec(clusterName, namespace, serviceName string, command string) (string, error) {
	cfg := &ClusterConfig{ClusterName: clusterName}
	clientset, token, err := NewAuthClient(cfg)
	if err != nil {
		return "", errors.Wrap(err, "create authentication client")
	}

	pods, err := clientset.CoreV1().
		Pods(namespace).
		List(context.TODO(), metav1.ListOptions{LabelSelector: "service=" + serviceName})
	if err != nil {
		return "", errors.Wrap(err, "check service available")
	}

	if len(pods.Items) == 0 {
		return "", errors.New("system's service isn't ready to execute the task")
	}

	restCfg, err := GetK8sRestConfig(clusterName, token)
	if err != nil {
		return "", errors.Wrap(err, "Get k8s config")
	}
	coreClient, err := corev1client.NewForConfig(restCfg)
	if err != nil {
		return "", errors.Wrap(err, "create rest client")
	}

	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	request := coreClient.RESTClient().
		Post().
		Namespace(pods.Items[0].Namespace).
		Resource("pods").
		Name(pods.Items[0].Name).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: pods.Items[0].Spec.Containers[0].Name,
			Command:   []string{"/bin/sh", "-c", command},
			Stdout:    true,
			Stderr:    true,
		}, scheme.ParameterCodec)

	logrus.WithField("command", command).Info("Run command")

	exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", request.URL())
	if err != nil {
		return "", errors.Wrap(err, "sending command request")
	}
	// oldState, err := terminal.MakeRaw(0)
	// if err != nil {
	// 	return "", errors.Wrap(err, "get response")
	// }
	// defer terminal.Restore(0, oldState)
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: buf,
		Stderr: errBuf,
	})
	if err != nil {
		return "", errors.Wrap(err, "read response")
	}

	return strings.TrimSpace(buf.String()) + strings.TrimSpace(errBuf.String()), nil
}
