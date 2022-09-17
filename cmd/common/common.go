package network

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ListOptions = metav1.ListOptions{}

type Options struct {
	AllNamespaces bool
	Namespace     string
	Selector      string
	Details       bool
}

type CountCmd interface {
	Validate() error
	Run() error
}
