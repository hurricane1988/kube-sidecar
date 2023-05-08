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

package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/wonderivan/logger"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

const (
	// WorkDir = "/opt"
	WorkDir = "/Users/hurricane/github/kube-sidecar"
)

// Config 定义全局配置
var Config *config

// 初始化全局配置
func init() {
	f, err := LoadConfigFromFile()
	if err != nil {
		Config = nil
	}
	Config = f
}

// LoadConfigFromFile 初始化配置文件
func LoadConfigFromFile() (*config, error) {
	// 获取当前程序执行路径
	// workDir, _ := os.Getwd()
	workDir := filepath.Join(WorkDir, "config", "conf")
	// 创建一个Viper实例并读取初始配置文件
	viperInstance := viper.New()
	// 加载viper获取配置路径
	viperInstance.AddConfigPath(workDir)
	viperInstance.AddConfigPath(homedir.HomeDir())
	viperInstance.AddConfigPath(".")
	// 设置读取配置文件
	// 设置读取的文件名
	viperInstance.SetConfigName("config")
	// 设置读取的文件后缀
	viperInstance.SetConfigType("yaml")
	// 匹配环境变量
	viperInstance.AutomaticEnv()

	// 执行读取配置文件
	if err := viperInstance.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Error("配置文件" + workDir + "/conf.yaml 不存在!")
		} else {
			logger.Error(err.Error())
		}
	}
	// 配置文件加载成功
	logger.Info("加载配置文件" + workDir + "/conf.yaml 成功!")
	/* viper动态加载配置 */
	// 监视配置文件是否发生更改
	viperInstance.WatchConfig()
	// 处理配置变更事件
	viperInstance.OnConfigChange(func(in fsnotify.Event) {
		// TODO: 联调测试
		logger.Info("配置文件" + in.Name + "发生变更")
		Config = NewConfig(viperInstance)
	})
	return NewConfig(viperInstance), nil
}

// NewConfig 创建默认的config
func NewConfig(viper *viper.Viper) *config {
	// 设置默认值
	viper.SetDefault("version.version", "v0.0.1")
	return &config{
		LoggingConfig:       *NewLogging(viper),
		Sidecar:             *NewSidecar(viper),
		NamespacesWhiteList: *NewNamespacesWhiteList(viper),
		DeploymentWhiteList: *NewDeploymentWhiteList(viper),
		FluentBitConfig:     *NewFluentBitConfig(viper),
		Version:             viper.GetString("version.version"),
	}
}

// NewLogging 日志配置
func NewLogging(viper *viper.Viper) *logging {
	// 设置默认值
	viper.SetDefault("logging.maxSize", 10)
	viper.SetDefault("logging.maxBackups", 3)
	viper.SetDefault("logging.maxAge", 3)
	return &logging{
		LogPath:    viper.GetString("logging.logPath"),
		Filename:   viper.GetString("logging.filename"),
		WriteLog:   viper.GetBool("logging.writeLog"),
		MaxSize:    viper.GetInt("logging.maxSize"),
		MaxBackups: viper.GetInt("logging.maxBackups"),
		MaxAge:     viper.GetInt("logging.maxAge"),
	}
}

// NewSidecar sidecar容器配置
func NewSidecar(viper *viper.Viper) *sidecar {
	// 设置默认值
	viper.SetDefault("sidecar.name", "sidecar")
	viper.SetDefault("sidecar.imagePullPolicy", "Always")
	viper.SetDefault("sidecar.requestsCPU", "50m")
	viper.SetDefault("sidecar.limitsCPU", "500m")
	viper.SetDefault("requestsMemory", "128Mi")
	viper.SetDefault("sidecar.limitsMemory", "512Mi")
	viper.SetDefault("sidecar.volumeName", "sidecar-config")
	viper.SetDefault("sidecar.volumeMount", "/fluent-bit/config")
	viper.SetDefault("sidecar.readOnly", true)
	return &sidecar{
		Name:            viper.GetString("sidecar.name"),
		Image:           viper.GetString("sidecar.image"),
		ImagePullPolicy: viper.GetString("sidecar.imagePullPolicy"),
		RequestsCPU:     viper.GetString("sidecar.requestsCPU"),
		RequestsMemory:  viper.GetString("sidecar.requestsMemory"),
		LimitCPU:        viper.GetString("sidecar.limitsCPU"),
		LimitMemory:     viper.GetString("sidecar.limitsMemory"),
		VolumeName:      viper.GetString("sidecar.volumeName"),
		VolumeMount:     viper.GetString("sidecar.volumeMount"),
		ReadOnly:        viper.GetBool("sidecar.readOnly"),
	}
}

// NewVersion 版本信息
func NewVersion(viper *viper.Viper) {

}

// NewNamespacesWhiteList 获取NamespacesWhiteList配置方法
func NewNamespacesWhiteList(viper *viper.Viper) *namespacesWhiteList {
	return &namespacesWhiteList{Names: viper.GetStringSlice("namespacesWhiteList")}
}

// NewDeploymentWhiteList 获取deploymentWhiteList配置方法
func NewDeploymentWhiteList(viper *viper.Viper) *deploymentWhiteList {
	return &deploymentWhiteList{Names: viper.GetStringSlice("deploymentWhiteList")}
}

// NewFluentBitConfig 获取FluentBit配置方法
func NewFluentBitConfig(viper *viper.Viper) *fluentBitConfig {
	// 设置默认值
	viper.SetDefault("fluentBit.serviceLogLevel", "info")
	viper.SetDefault("fluentBit.inputMemBufLimit", "20MB")
	viper.SetDefault("fluentBit.inputRefreshInterval", 20)
	return &fluentBitConfig{
		ServiceLogLevel:      viper.GetString("fluentBit.serviceLogLevel"),
		InputMemBufLimit:     viper.GetString("fluentBit.inputMemBufLimit"),
		InputRefreshInterval: viper.GetInt("fluentBit.inputRefreshInterval"),
	}
}
