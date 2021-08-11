package k8s

import (
	"bytes"
	"context"
	"time"

	"golang.org/x/crypto/ssh/terminal"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const testClusterName = "dev"
const testNamespace = "qa"
const testCommand = "ls db-dump"

func TestNamespaces() {
	logrus.Info("DEBUG: start prod")
	// Setup the basic EKS cluster info
	cfg := &ClusterConfig{ClusterName: testClusterName}
	logrus.Info("DEBUG: done config")

	logrus.Info(time.Now(), "DEBUG: start auth")
	clientset, _, err := NewAuthClient(cfg)
	if err != nil {
		logrus.Info("ERROR: ", err.Error())
		return
	}
	logrus.Info("DEBUG: done auth")

	// Call Kubernetes API here
	logrus.Info("DEBUG: start list namespaces")
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Info(time.Now(), "ERROR: ", err.Error())
		return
	}
	logrus.Info("DEBUG: done list namespaces")

	logrus.Info("DEBUG: start print namespaces")
	for i, namespace := range namespaces.Items {
		logrus.Info("DEBUG: [", i, "] ", namespace.Name, " ", namespace.Labels)
	}
	logrus.Info("DEBUG: done print namespaces")
}

func TestPods() {
	logrus.Info("DEBUG: start config")
	// Setup the basic EKS cluster info
	cfg := &ClusterConfig{ClusterName: testClusterName}
	logrus.Info("DEBUG: done config")

	logrus.Info(time.Now(), "DEBUG: start auth")
	clientset, _, err := NewAuthClient(cfg)
	if err != nil {
		logrus.Info("ERROR: ", err.Error())
		return
	}
	logrus.Info("DEBUG: done auth")

	// Call Kubernetes API here
	logrus.Info("DEBUG: start list pods")
	namespace := testNamespace
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Info(time.Now(), "ERROR: ", err.Error())
		return
	}
	logrus.Info("DEBUG: done list pods")

	logrus.Info("DEBUG: start print pods")
	for i, pod := range pods.Items {
		logrus.Info("DEBUG: [", i, "] ", pod.Name)
	}
	logrus.Info("DEBUG: done print pods")
}

func TestExecPod() {
	logrus.Info("DEBUG: start config")
	cfg := &ClusterConfig{ClusterName: testClusterName}
	logrus.Info("DEBUG: done config")

	logrus.Info(time.Now(), "DEBUG: start auth")
	clientset, token, err := NewAuthClient(cfg)
	if err != nil {
		logrus.Info("ERROR: ", err.Error())
		return
	}
	logrus.Info("DEBUG: done auth")

	logrus.Info("DEBUG: start list pods")
	namespace := testNamespace
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "service=cerberus"})
	if err != nil {
		logrus.Info(time.Now(), "ERROR: ", err.Error())
		return
	}
	logrus.Info("DEBUG: done list pods")

	logrus.Info("DEBUG: start print pods ", len(pods.Items))
	for i, pod := range pods.Items {
		logrus.Info("DEBUG: [", i, "] ", pod.Name)
	}
	logrus.Info("DEBUG: done print pods")

	if len(pods.Items) == 0 {
		logrus.Info("DEBUG: Can't exec because no pod")
		return
	}

	restCfg, err := GetK8sRestConfig(testClusterName, token)
	if err != nil {
		logrus.Info("DEBUG: Can't get k8s rest config: ", err)
		return
	}

	coreClient, err := corev1client.NewForConfig(restCfg)
	if err != nil {
		logrus.Info("DEBUG: Can't create clientCOnfig: ", err)
		return
	}

	logrus.Info("DEBUG: create coreClient successful")

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
			Command:   []string{"/bin/sh", "-c", testCommand},
			Stdout:    true,
			Stderr:    true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", request.URL())
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: buf,
		Stderr: errBuf,
	})

	if err != nil {
		logrus.Info("DEBUG: Can't exec: ", errors.Wrapf(err, "Failed executing command %s on %v/%v", testCommand, pods.Items[0].Namespace, pods.Items[0].Name))
		return
	}

	logrus.Info("DEBUG: Run successful: ", buf.String(), errBuf.String())

}
