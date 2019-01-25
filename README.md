# nginx_mirror
Use mirror upstream to statstic http info.


## 功能
* 统计nginx流量中的TPS
* 记录流量到influxdb
* 提供接口用于prometheus获取tps排行

## 依赖
* nginx_http_mirror_module(nginx >= 1.13.4)
* 在nginx配置中使用proxy_set_header 见 *nginx配置*
* influxdb

## nginx 配置
```
upstream web {
        server 127.0.0.1:12345;
    }

    upstream mirror {
        server 127.0.0.1:9999;
    }


    server {
        listen       80;
        server_name  localhost;

        location / {
            mirror /mirror;
            proxy_pass http://web/mirror;
        }

        location /mirror{
            internal;
            proxy_pass http://mirror/mirror;
            proxy_set_header X-Original-URI $request_uri;
            proxy_set_header Host $host;
            proxy_set_header  X-Real-IP  $remote_addr;
            proxy_set_header  X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto  $scheme;
            proxy_set_header Nginx-IP 1.1.1.1;
            proxy_set_header SSL-Protocil $ssl_protocol;

        }

    }
```