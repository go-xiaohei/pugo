```toml
title = "Deploy"
date = "2016-02-04 15:00:00"
slug = "en/docs/cmd/deploy"
hover = "docs"
lang = "en"
template = "docs.html"
```

`deploy` command deploys static files to other platform after building.

```go
pugo build
pugo deploy [method] [--options]
```

`PuGo` can deploy via `FTP`, `SFTP`, `Git` and `AWS S3`, `Qiniu Storage` methods.

Read [Deploy](/en/docs/deploy/standalone.html) doc to get more details for each method.