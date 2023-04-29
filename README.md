## kubernetes deployment 便车容器自动注入
> kube-sidecar是以fluentBit为日志采集方案，结合kubernetes sidecar概念，实现创建根据如下deployment的注释，自动添加sidecar日志采集容器
- [x] 为deployment设置是否开启sidecar注释
```yaml
deployment.kubernetes.io/sidecar: 'true'
```
> FluentBit相关
- [x] [[FluentBit Github仓库]](https://github.com/fluent/fluent-bit)
- [x] [[FluentBit文档中心]](https://fluentbit.io/)
- [x] [[fluent-bit helm chart]](https://github.com/fluent/helm-charts/tree/main/charts/fluent-bit)
- [x] [[参考配置]](https://clickvisual.gocn.vip/clickvisual/07collect/fluent-bit-configuration-reference.html#_2-3-input-kubernetes-conf-%E9%85%8D%E7%BD%AE)