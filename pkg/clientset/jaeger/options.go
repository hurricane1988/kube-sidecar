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

package jaeger

// Options 链路跟踪相关、
type Options struct {
	Enable bool   `json:"enable,omitempty" yaml:"enable,omitempty" xml:"enable,omitempty"`
	Scheme string `yaml:"scheme,omitempty" xml:"scheme,omitempty" json:"scheme,omitempty"`
	Host   string `json:"host,omitempty" yaml:"host,omitempty" xml:"host,omitempty"`
	Port   string `json:"port,omitempty" xml:"port,omitempty" yaml:"port,omitempty"`
	Path   string `json:"path,omitempty" xml:"path,omitempty" yaml:"path,omitempty"`
}

// NewJaegerOptions jaeger链路跟踪配置
func NewJaegerOptions() *Options {
	return &Options{
		Enable: true,
		Scheme: "http",
		Host:   "127.0.0.1",
		Port:   "16686",
		Path:   "/api/traces",
	}
}
