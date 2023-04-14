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

package main

import (
    "context"
    "fmt"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    "kube-sidecar/utils/clients/k8s"
)

func main() {
    options := k8s.NewKubernetesOptions()
    client, _ := k8s.NewKubernetesClient(options)
    podList, _ := client.Kubernetes().CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

    for _, pod := range podList.Items {
        fmt.Println(pod.Name, pod.Namespace)
    }
}
