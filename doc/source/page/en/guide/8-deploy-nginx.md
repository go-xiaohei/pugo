```toml
title = "Deploy to Nginx"
date = "2016-02-05 15:00:00"
slug = "en/guide/deploy-nginx"
hover = "guide"
lang = "en"
template = "guide.html"
```

There are two ways to deploy `PuGo` contents to public.

### Static Website

If run build command, `PuGo` builds static contents to directory.

So `nginx` can serve a static web site in that directory.

In `nginx.conf`, use `server` block to define a website.

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

Set `[your-pugo-build-directory]` correct. Use absolute path better. Be careful of the permission.

### Proxy

Use server command, `PuGo` run http server by self, such as on `http://127.0.0.1:9899`. `Nginx` can run proxy to the site to make the site public.

In `server` block, insert this block:

```nginx
location / {
    proxy_pass   http://127.0.0.1:9899
}
```

But not recommended. You must add correct `Cache-Header` to get higher experience,