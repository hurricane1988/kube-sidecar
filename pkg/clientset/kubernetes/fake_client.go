/*
Copyright 2023 QKP Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubernetes

import (
	promresourcesclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type FakeClient struct {
	// kubernetes clientset interface
	K8sClient kubernetes.Interface
	// discovery clientset
	DiscoveryClient *discovery.DiscoveryClient
	// generated kubernetes
	prometheusClient promresourcesclient.Interface
	MasterURL        string
	KubeConfig       *rest.Config
}

func NewFakeClientSets(
	k8sClient kubernetes.Interface,
	discoveryClient *discovery.DiscoveryClient,
	prometheusClient promresourcesclient.Interface,
	masterURL string, kubeConfig *rest.Config) Client {
	return &FakeClient{
		K8sClient:        k8sClient,
		DiscoveryClient:  discoveryClient,
		prometheusClient: prometheusClient,
		MasterURL:        masterURL,
		KubeConfig:       kubeConfig,
	}
}

func (n *FakeClient) Kubernetes() kubernetes.Interface {
	return n.K8sClient
}

func (n *FakeClient) Discovery() discovery.DiscoveryInterface {
	return n.DiscoveryClient
}

func (n *FakeClient) Prometheus() promresourcesclient.Interface {
	return n.prometheusClient
}

func (n *FakeClient) Master() string {
	return n.MasterURL
}

func (n *FakeClient) Config() *rest.Config {
	return n.KubeConfig
}
