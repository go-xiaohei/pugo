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

你可以使用 `--from` 和 `--to` 修改编译的设置，但是不支持迁移和部署.

```bash
pugo server --from="your-source" --to="your-destination"
```