package cmd

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var listOptions = metav1.ListOptions{}

type Options struct {
	allNamespaces bool
	namespace     string
	selector      string
}

type CountCmd interface {
	validate() error
	run() error
}
