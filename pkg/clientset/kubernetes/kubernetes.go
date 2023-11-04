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
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

// Client 定义kubernetes接口
type Client interface {
	Kubernetes() kubernetes.Interface
	Discovery() discovery.DiscoveryInterface
	Master() string
	Config() *rest.Config
}

// 定义kubernetes客户端结构体
type kubernetesClient struct {
	// kubernetes client interface
	k8s kubernetes.Interface
	// discovery client
	discoveryClient *discovery.DiscoveryClient
	prometheus      promresourcesclient.Interface
	master          string
	config          *rest.Config
}

// NewKubernetesClientOrDie 实例化kubernetes 客户端creates KubernetesClient and panic if there is an error
func NewKubernetesClientOrDie(options *KubernetesOptions) Client {
	config, err := clientcmd.BuildConfigFromFlags("", options.KubeConfig)
	if err != nil {
		panic(err)
	}

	config.QPS = options.QPS
	config.Burst = options.Burst
	k := &kubernetesClient{
		k8s:             kubernetes.NewForConfigOrDie(config),
		discoveryClient: discovery.NewDiscoveryClientForConfigOrDie(config),
		prometheus:      promresourcesclient.NewForConfigOrDie(config),
		master:          config.Host,
		config:          config,
	}

	if options.Master != "" {
		k.master = options.Master
	}
	if !strings.HasPrefix(k.master, "http://") && !strings.HasPrefix(k.master, "https://") {
		k.master = "https://" + k.master
	}
	return k
}

// NewKubernetesClient 创建kubernetesClient客户端
func NewKubernetesClient(options *KubernetesOptions) (Client, error) {
	// 定义局部kubernetes客户端
	var k kubernetesClient
	config, err := clientcmd.BuildConfigFromFlags("", options.KubeConfig)
	if err != nil {
		return nil, err
	}
	config.QPS = options.QPS
	config.Burst = options.Burst
	k.k8s, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	k.discoveryClient, err = discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}
	k.prometheus, err = promresourcesclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	k.master = options.Master
	k.config = config

	return &k, nil
}

// Kubernetes 实例化Kubernetes方法
func (k *kubernetesClient) Kubernetes() kubernetes.Interface {
	return k.k8s
}

// Discovery 实例化Discovery()方法
func (k *kubernetesClient) Discovery() discovery.DiscoveryInterface {
	return k.discoveryClient
}

// Master 实例化Master()方法
func (k *kubernetesClient) Master() string {
	return k.master
}

// Config 实例化rest.config客户端
func (k *kubernetesClient) Config() *rest.Config {
	return k.config
}
