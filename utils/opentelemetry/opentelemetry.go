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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"log"
	"strings"
	"time"

	"kube-sidecar/config"
	lg "kube-sidecar/utils/logging"
)

var (
	jaegerConfig = config.Config.JaegerConfig
)

// TracerProvider 设置tracerProvider方法
func TracerProvider(service, environment string, id int64) (*tracesdk.TracerProvider, error) {
	endpoint := strings.Join([]string{jaegerConfig.Scheme, ":/", jaegerConfig.Host, string(rune(jaegerConfig.Port)), jaegerConfig.Path}, "/")
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
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
func RegisterGlobalTracerProvider(tracerName, spanName, service, environment string, id int64, obj interface{}, Func func(obj ...interface{})) {
	tp, err := TracerProvider(service, environment, id)
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

	ctx, span := tr.Start(ctx, spanName)
	defer span.End()
	// Context 向下传递
	Func(ctx, tracerName, spanName, obj)
}
