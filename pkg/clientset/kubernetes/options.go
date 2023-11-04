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
	"github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
	"os"
	"os/user"
	"path"
	"reflect"
)

// KubernetesOptions 定义kubernetes配置结构体
type KubernetesOptions struct {
	// kubeconfig path, if not specified, will use
	// in cluster way to create clientset
	KubeConfig string `json:"kubeconfig" yaml:"kubeconfig"`

	// kubernetes apiserver public address, used to generate kubeconfig
	// for downloading, default to host defined in kubeconfig
	// +optional
	Master string `json:"master,omitempty" yaml:"master"`

	// kubernetes clientset qps
	// +optional
	QPS float32 `json:"qps,omitempty" yaml:"qps"`

	// kubernetes clientset burst
	// +optional
	Burst int `json:"burst,omitempty" yaml:"burst"`
}

// NewKubernetesOptions returns a `zero` instance
func NewKubernetesOptions() (option *KubernetesOptions) {
	option = &KubernetesOptions{
		// 在Go语言中，1e6是科学计数法表示的1乘以10的6次方，也就是1000000（一百万）。因此，float32类型的1e6是1000000.0
		QPS:   1e6,
		Burst: 1e6,
	}

	homePath := homedir.HomeDir()
	if homePath == "" {
		// 获取当前用户信息（Uid、Gid、Username、Name、HomeDir）try os/user.HomeDir when $HOME is unset.
		if u, err := user.Current(); err == nil {
			homePath = u.HomeDir
		}
	}
	// 定义当前用户的kubeconfig文件
	userHomeConfig := path.Join(homePath, ".kube/config")
	// 检查userHomeConfig文件是否能返回文件信息(Name、Size、Mode、ModTime、IsDir、Sys),如果err为空则正常
	if _, err := os.Stat(userHomeConfig); err == nil {
		option.KubeConfig = userHomeConfig
	}
	return
}

// Validate 校验kubernetesOptions方法
func (k *KubernetesOptions) Validate() []error {
	var errors []error
	if k.KubeConfig != "" {
		if _, err := os.Stat(k.KubeConfig); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// ApplyTo 重新生效KubernetesOptions
func (k *KubernetesOptions) ApplyTo(options *KubernetesOptions) {
	Override(options, k)
}

// AddFlags 增加flag方法
func (k *KubernetesOptions) AddFlags(fs *pflag.FlagSet, c *KubernetesOptions) {
	fs.StringVar(
		&k.KubeConfig, "kubeconfig", c.KubeConfig, ""+
			"Path for kubernetes kubeconfig file, if left blank, will use "+
			"in cluster way.")
	fs.StringVar(&k.Master, "master", c.Master, ""+
		"Used to generate kubeconfig for downloading, if not specified, will use host in kubeconfig.")
}

// Override 覆盖对象方法
func Override(left interface{}, right interface{}) {
	if reflect.ValueOf(left).IsNil() || reflect.ValueOf(right).IsNil() {
		return
	}

	if reflect.ValueOf(left).Type().Kind() != reflect.Ptr ||
		reflect.ValueOf(right).Type().Kind() != reflect.Ptr ||
		reflect.ValueOf(left).Kind() != reflect.ValueOf(right).Kind() {
		return
	}

	oldVal := reflect.ValueOf(left).Elem()
	newVal := reflect.ValueOf(right).Elem()

	for i := 0; i < oldVal.NumField(); i++ {
		val := newVal.Field(i).Interface()
		if !reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface()) {
			oldVal.Field(i).Set(reflect.ValueOf(val))
		}
	}
}
