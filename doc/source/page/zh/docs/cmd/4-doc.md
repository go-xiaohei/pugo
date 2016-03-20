```toml
title = "Doc"
date = "2016-02-04 15:00:00"
slug = "zh/docs/cmd/doc"
hover = "docs"
lang = "zh"
template = "docs.html"
```

`doc` 命令会在本地启动文档 HTTP 服务：

```go
pugo new site
pugo doc --addr="0.0.0.0:9899"
```

执行 `new site` 后, `PuGo` 会释放文档数据到 `doc` 文件夹。 你可以使用 `doc` 命令编译和浏览文档内容。

`--addr` 设置 HTTP 服务的地址和端口，默认是 `0.0.0.0:9899`。