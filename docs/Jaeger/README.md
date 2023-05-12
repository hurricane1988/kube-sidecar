## Jaeger简介
Jaeger：开源、端到端的分布式链路链路追踪框架,在复杂分布式系统中进行监控和故障排除
> Jaeger 解决的问题
- [ ] 分布式事务监控
- [ ] 性能和延迟优化
- [ ] 根本原因分析
- [ ] 服务依赖分析
- [ ] 分布式上下文传播
- [Jaeger官网](https://www.jaegertracing.io/)
> Jaeger服务docker部署,部署采用官网的all-in-one的 镜像
```shell
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.35
```
> 执行命令后，打开http://127.0.0.1:16686/，出现下图则，安装成功
