```toml
title = "Server"
date = "2016-02-04 15:00:00"
slug = "zh/docs/cmd/server"
hover = "docs"
lang = "zh"
template = "docs.html"
```

`server` 启动HTTP服务展示站点。

```go
pugo server --addr="0.0.0.0:9899" --source="source" --dest="dest" --theme="theme/default" --static --debug
```

`--addr` 设置 HTTP 服务的地址和端口，默认是 `0.0.0.0:9899`。

`--source`, `--dest` 和 `--theme` 设置内容、编译和主题目录，来源于 `build` 命令。

`--static` 只展示 `--dest` 静态内容，但是需要正确的 `--source` 加载必要数据。

`--debug` 打印更多调试信息。

### 注意

当执行 `server` 时， `PuGo` 会立刻编译内容，然后启动 HTTP 服务，同时监听文件修改，随时直接编译最新内容。因此 `server` 命令更适用于开发或正在写作的时候，预览修改的效果。

**我建议使用 server --static 启动 HTTP 对外 HTTP 服务**.

当然，我更期望直接使用 Web 服务器展示静态内容。