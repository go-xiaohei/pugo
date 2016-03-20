```toml
title = "内置 HTTP 服务"
date = "2016-02-05 15:00:00"
slug = "zh/guide/standalone-server"
hover = "guide"
lang = "zh"
template = "guide.html"
```

`server` 命令可以启动内置的 HTTP 服务：

```bash
pugo server [--addr=0.0.0.0:9899]
```

默认监听在 `http://0.0.0.0:9899`. `--addr` 可以自定义监听的 IP 和端口。

开启 HTTP 服务时会立刻编译最新的内容，并在生成的目录上启动静态 HTTP 服务。

你可以使用 `--source`、`--dest` 和 `--theme` 修改编译的设置。

```bash
pugo server --source="your-source" --dest="your-destination"
```

### 静态服务

如果只是启动静态文件服务，使用`--flag`：

```bash
pugo server --source="source" --static
```

你需要设置正确的`--source`目录，以便静态服务可以得到正确的 URL 地址。