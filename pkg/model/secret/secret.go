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

package secret

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"

	lg "kube-sidecar/utils/logging"

	"kube-sidecar/utils/clients/k8s"
)

// CreateFluentBitSecret 创建secret方法
func CreateFluentBitSecret(backend, name, namespace string, client k8s.Client, fluent FluentBitConf) error {
	// 生成fluentBit secret []byte数据
	data, err := GenerateFluentBitConfig(backend, fluent)
	if err != nil {
		lg.Logger.Error("生成fluentBit配置文件失败,错误信息" + err.Error())
		return nil
	}
	// 创建一个新的secret对象
	newSecret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      strings.Join([]string{name, "sidecar"}, "-"),
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"fluent-bit.conf": data,
		},
	}
	// 创建secret
	createdSecret, err := client.Kubernetes().CoreV1().Secrets(namespace).Create(context.Background(), newSecret, v1.CreateOptions{})
	if err != nil {
		lg.Logger.Error(err.Error())
	}
	lg.Logger.Info("namespace " + createdSecret.Namespace + "创建secret " + createdSecret.Name + "成功!")
	return nil
}

// GenerateFluentBitConfig 创建fluentBit配置文件模版,输出位[]byte
func GenerateFluentBitConfig(backend string, fluent FluentBitConf) ([]byte, error) {
	data, err := FluentBitTemplate(backend, fluent)
	if err != nil {
		lg.Logger.Error(err.Error())
		return nil, err
	}
	return data, nil
}
