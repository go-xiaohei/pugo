```toml
title = "Deploy to Nginx"
date = "2016-02-05 15:00:00"
slug = "en/guide/deploy-nginx"
hover = "guide"
lang = "en"
template = "guide.html"
```

There are two ways to deploy `PuGo` contents to public byself, or deploy to public cloud service such as git, Amazon Cloud via `deploy` command.

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

**In this case, please write navigation link with correct file path, such as `/page.html` not `/page`.**

### Proxy

Use server command, `PuGo` run http server by self, such as on `http://127.0.0.1:9899`. `Nginx` can run proxy to the site to make the site public.

In `server` block, insert this block:

```nginx
location / {
    proxy_pass   http://127.0.0.1:9899
}
```

But not recommended. You must add correct `Cache-Header` to get higher experience,

### Cloud

`PuGo` support [FTP](/en/docs/deploy/ftp-sftp.html), [SFTP](/en/docs/deploy/ftp-sftp.html), [Git](/en/docs/deploy/git.html), and [Storage Service](/en/docs/deploy/cloud.html). 

For example, the steps of deploying to github repository, in branch **gh-pages**:

- `git clone` your repository to [dir1]. Please use `git://` or `https://username:password@repository_url.git` . (PuGo can't type in username and password when pushing now).

- `git checkout gh-pages` in [dir1]

- `pugo build --dest="dir2"`, build files to [dir2]

- `pugo deploy git --local="dir2" --repo="dir1" --branch="gh-pages"`

You can read deployment document to learn further guides.