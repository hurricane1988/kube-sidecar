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

package deploy

import (
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/tools/cache"

    dp "kube-sidecar/pkg/model/deploy"

    "kube-sidecar/utils/clients/k8s"
)

// DeploymentWatcher 定义deployment的watcher
func DeploymentWatcher(clientSet k8s.Client, container corev1.Container) {
    const (
        annotationKey = "deployment.kubernetes.io/sidecar"
    )
    // 初始化sidecar
    deploy := dp.Sidecar
    watcher := cache.NewListWatchFromClient(
        clientSet.Kubernetes().AppsV1().RESTClient(),
        "deployment", metav1.NamespaceAll,
        nil)
    _, controller := cache.NewInformer(
        watcher,
        &appsv1.Deployment{},
        0,
        cache.ResourceEventHandlerFuncs{
            UpdateFunc: func(oldObj, newObj interface{}) {
                deployment := newObj.(*appsv1.Deployment)
                annotations := deployment.Annotations
                if val, ok := annotations[annotationKey]; ok && val == "true" {
                    deploy.AddDeploymentSidecar(&container, deployment)
                }
            },
        },
    )
    // 运行Controller
    stop := make(chan struct{})
    defer close(stop)
    controller.Run(stop)
}
