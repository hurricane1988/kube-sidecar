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
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    "kube-sidecar/utils/clients/k8s"
)

// CreateSecret 创建secret方法
func CreateSecret(client k8s.Client, secret *corev1.Secret) error {
    // 定义全局error
    var err error
    _, err = client.Kubernetes().
        CoreV1().
        Secrets(secret.Namespace).
        Create(context.Background(), secret, metav1.CreateOptions{})
    if err != nil {
        return err
    }
    return err
}

// CreateFluentBitSecret 创建fluentBit Secret
func CreateFluentBitSecret()  {
    fluentbitConf := FluentBitConf{
        ServiceLogLevel:      "info",
        InputLogPath:         logPath,
        InputAppName:         r.Name,
        InputAppTag:          r.Name + ".log",
        InputMemBufLimit:     "20MB",
        InputRefreshInterval: 15,
        OutputEsHost:         address,
        OutputEsPort:         port,
        OutputEsIndex:        index,
        OutputEsUser:         username,
        OutputEsPassword:     password,
    }
}
    data, err := FluentBitTemplate(FluentBitConf{})
    if err != nil {
        panic(err)
}