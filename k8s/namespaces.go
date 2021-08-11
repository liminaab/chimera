package k8s

import (
	"context"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNamespaces(clusterName string) ([]string, error) {
	cfg := &ClusterConfig{ClusterName: clusterName}

	clientset, _, err := NewAuthClient(cfg)
	if err != nil {
		return []string{}, err
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return []string{}, errors.Wrap(err, "Can't query namespaces k8s")
	}

	result := make([]string, 0, len(namespaces.Items))
	for _, namespace := range namespaces.Items {
		result = append(result, namespace.Name)
	}
	return result, nil
}
