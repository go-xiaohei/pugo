```toml
title = "独立部署"
date = "2016-02-04 15:00:00"
slug = "zh/docs/deploy/standalone"
hover = "docs"
lang = "zh"
template = "docs.html"
```

`PuGo` 能独立启动 HTTP 服务。 发布时，运行 `server` 命令：

```bash
pugo build
pugo server --static
```

如果还在写作中，直接使用 `server` 命令及时编译最新内容：

```bash
pugo server
```