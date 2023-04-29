```shell
docker run -d \
    -p 0.0.0.0:9200:9200 \
    -e "discovery.type=single-node" \
    -e "ELASTIC_PASSWORD=admin" \
    -e "xpack.security.enabled=true" \
    -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
    -e "ELASTIC_USERNAME=admin" \
    --name elasticsearch \
    elasticsearch:7.10.2
```
```shell
docker run -d \
    -p 9200:9200 -p 9300:9300 \
    -e "discovery.type=single-node" \
    -e "ELASTIC_USERNAME=elastic" \
    -e "ELASTIC_PASSWORD=Harbor12345" \
    -e "xpack.security.enabled=true" \
    -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
    --name elasticsearch \
    elasticsearch:7.10.2
```
```shell
curl -u elastic:Harbor12345 http://localhost:9200/
```