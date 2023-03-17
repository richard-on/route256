package kube

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Client is a wrapper for *kubernetes.Clientset.
type Client struct {
	clientSet  *kubernetes.Clientset
	kubeConfig config.Kubernetes
}

// NewClient creates new Kubernetes gRPC client.
func NewClient(clientSet *kubernetes.Clientset, kubeConfig config.Kubernetes) *Client {
	return &Client{
		clientSet:  clientSet,
		kubeConfig: kubeConfig,
	}
}

// GetReplicaCount returns current number of this pod replicas in a Kubernetes cluster.
func (c *Client) GetReplicaCount(ctx context.Context) int {
	pods, err := c.clientSet.CoreV1().
		Pods(c.kubeConfig.Namespace).
		List(ctx, metav1.ListOptions{
			LabelSelector: c.kubeConfig.LabelSelector,
		})
	if err != nil || len(pods.Items) < 1 {
		return 1
	}

	return len(pods.Items)
}
