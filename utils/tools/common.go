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

package tools

import (
	appsv1 "k8s.io/api/apps/v1"
	lg "kube-sidecar/pkg/clientset/logging"
	"math/rand"
	"time"
)

// WhetherExists 检查某个字符是否在slice切片中
func WhetherExists(key string, slice []string) bool {
	m := make(map[interface{}]bool)
	for _, v := range slice {
		m[v] = true
	}
	return m[key]
}

// SetDefaultValueNotExist 实现读取annotations值并默认值
func SetDefaultValueNotExist(key string, defaultValue string) string {
	switch key {
	case "":
		return defaultValue
	default:
		return key
	}
}

// WorkloadContainerNames 获取工作负载Deployment\StatefulSet\DaemonSet的所有containers名称
func WorkloadContainerNames(objType string, object interface{}) []string {
	var (
		names []string
	)
	switch objType {
	case "DaemonSet":
		obj, ok := object.(*appsv1.DaemonSet)
		if !ok {
			lg.Logger.Error("对象转为Deployment失败,错误信息")
		}
		for _, dp := range obj.Spec.Template.Spec.Containers {
			names = append(names, dp.Name)
		}
	case "StatefulSet":
		obj, ok := object.(*appsv1.StatefulSet)
		if !ok {
			lg.Logger.Error("对象转为StatefulSet失败,错误信息")
		}
		for _, dp := range obj.Spec.Template.Spec.Containers {
			names = append(names, dp.Name)
		}
	default:
		obj, ok := object.(*appsv1.Deployment)
		if !ok {
			lg.Logger.Error("对象转为Deployment失败,错误信息")
		}
		for _, dp := range obj.Spec.Template.Spec.Containers {
			names = append(names, dp.Name)
		}
	}
	return names
}

// RandomInt64 生成随机int64
func RandomInt64() int64 {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(100)
}
