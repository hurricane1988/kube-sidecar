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

// TODO: https://github.com/open-telemetry/opentelemetry-go

package opentelemetry

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"kube-sidecar/pkg/clientset/fluent"
	jg "kube-sidecar/pkg/clientset/jaeger"
	"kube-sidecar/pkg/clientset/kubernetes"
	lg "kube-sidecar/pkg/clientset/logging"
	"kube-sidecar/pkg/clientset/sidecar"
	"kube-sidecar/pkg/clientset/workload"
	"log"
	"net/url"
	"time"

	"kube-sidecar/pkg/controller/deploy"
)

type openTelemetry struct {
	K8sClient kubernetes.Client
	FluentBit fluent.Options
	Sidecar   sidecar.Options
	Jeager    jg.Options
	WhiteList workload.Options
}

type OpenTelemetry interface {
	TracerProvider(service, environment string, id int64) (*tracesdk.TracerProvider, error)
}

func NewOpenTelemetry(k8sClient kubernetes.Client, fluentBit fluent.Options, sidecar sidecar.Options, jeager jg.Options, whiteList workload.Options) OpenTelemetry {
	return &openTelemetry{
		K8sClient: k8sClient,
		FluentBit: fluentBit,
		Sidecar:   sidecar,
		Jeager:    jeager,
		WhiteList: whiteList,
	}
}

// TracerProvider 设置tracerProvider方法
func (o *openTelemetry) TracerProvider(service, environment string, id int64) (*tracesdk.TracerProvider, error) {
	endpoint, _ := url.Parse(fmt.Sprintf("%s://%s", o.Jeager.Scheme, o.Jeager.Host))
	// 添加端口号
	endpoint.Host = fmt.Sprintf("%s:%s", o.Jeager.Host, o.Jeager.Port)
	// 拼接路径
	endpoint.Path += o.Jeager.Path
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint.String())))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}

// RegisterGlobalTracerProvider 注册全局的tracerProvider
func (o *openTelemetry) RegisterGlobalTracerProvider(tracerName, spanName, service, environment string, id int64) {
	tp, err := o.TracerProvider(service, environment, id)
	if err != nil {
		lg.Logger.Error("初始化TracerProvider失败,错误信息" + err.Error())
	}
	// 注册全局TracerProvider
	otel.SetTracerProvider(tp)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	tr := tp.Tracer(tracerName)
	_, span := tr.Start(ctx, spanName)
	defer span.End()
	// Context 向下传递
	deploy.NewDeployment(o.K8sClient, o.FluentBit, o.Sidecar, o.Jeager, o.WhiteList).Watch(ctx, tracerName, spanName)
}
