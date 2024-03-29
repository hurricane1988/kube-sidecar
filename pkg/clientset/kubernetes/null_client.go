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

type nullClient struct {
}

func NewNullClient() Client {
	return &nullClient{}
}

func (n *nullClient) Kubernetes() kubernetes.Interface {
	return nil
}

func (n *nullClient) Discovery() discovery.DiscoveryInterface {
	return nil
}

func (n *nullClient) Prometheus() promresourcesclient.Interface {
	return nil
}

func (n *nullClient) Master() string {
	return ""
}

func (n *nullClient) Config() *rest.Config {
	return nil
}
