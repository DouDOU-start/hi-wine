# 配置模板

```yaml
# https://goframe.org/docs/web/server-config-file-template
server:
  address:     "0.0.0.0:8000"
  openapiPath: "/api.json"
  swaggerPath: "/swagger"
  url:         "http://localhost:8000"

# https://goframe.org/docs/core/glog-config
logger:
  level : "all"
  stdout: true

# https://goframe.org/docs/core/gdb-config-file
database:
  default:
    link: "mysql:root:root@tcp(127.0.0.1:3306)/hi-wine"

wechat:
  appid: xxxxxxxxxxxxxxxx
  secret: xxxxxxxxxxxxxxxx

# MinIO配置
minio:
  endpoint: "127.0.0.1:9000"
  accessKey: "minio"
  secretKey: "Minio8888"
  useSSL: false
  bucket: "hi-wine"
  region: ""
  domain: "http://127.0.0.1:9000"
```