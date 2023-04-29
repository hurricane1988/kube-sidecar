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
