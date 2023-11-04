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

package version

type Options struct {
	GitVersion   string `json:"gitVersion,omitempty" xml:"gitVersion,omitempty" yaml:"gitVersion,omitempty"`
	GoVersion    string `json:"goVersion,omitempty" xml:"goVersion,omitempty" yaml:"goVersion,omitempty"`
	BuildVersion string `json:"buildVersion,omitempty" xml:"buildVersion,omitempty" yaml:"buildVersion,omitempty"`
}

// NewVersionOptions 版本信息
func NewVersionOptions() *Options {
	return &Options{
		GitVersion:   "0.0.0",
		GoVersion:    "1.21",
		BuildVersion: "v1.0.0",
	}
}
