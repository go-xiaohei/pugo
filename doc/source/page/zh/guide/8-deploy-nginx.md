```toml
title = "部署到 Nginx"
date = "2016-02-05 15:00:00"
slug = "zh/guide/deploy-nginx"
hover = "guide"
lang = "zh"
template = "guide.html"
```

有两种方式部署 `PuGo` 的内容。

### 纯静态网站

编译之后，生成的文件会保存在编译目录中.

`Nginx` 可以将目录认为是静态网站提供服务。

`nginx.conf` 中 `server` 内容块定义网站。

```nginx
server {
        listen       80;
        server_name  domain;

        #charset utf-8;

        #access_log  logs/host.access.log  main;

        location / {
            root   [your-pugo-build-directory];
            index  index.html index.htm;
        }
}
```

修改 `[your-pugo-build-directory]` 到您编译的目录，建议使用绝对地址，并注意权限。

### 代理

`PuGo` 也可以启动 HTTP 服务在 `http://127.0.0.1:9899`. `Nginx` 代理到这个位置：

`server` 中添加这些内容：

```nginx
location / {
    proxy_pass   http://127.0.0.1:9899
}
```

你需要配置如 `Cache-Header` 去保证更好的静态性能。