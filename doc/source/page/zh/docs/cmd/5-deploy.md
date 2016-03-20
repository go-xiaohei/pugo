```toml
title = "Deploy"
date = "2016-02-04 15:00:00"
slug = "zh/docs/cmd/deploy"
hover = "docs"
lang = "zh"
template = "docs.html"
```

`deploy` 命令部署静态内容到别的平台。

```go
pugo build
pugo deploy [method] [--options]
```

`PuGo` 支持通过 `FTP`, `SFTP`, `Git` 和 `AWS S3`, `Qiniu Storage` 部署.

阅读 [Deploy](/zh/docs/deploy/standalone.html) ，了解各种部署方式的相关内容。