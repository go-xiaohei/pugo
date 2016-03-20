```toml
title = "Build"
date = "2016-02-04 15:00:00"
slug = "zh/docs/cmd/build"
hover = "docs"
lang = "zh"
template = "docs.html"
```

`build` 命令编译所有内容到静态 HTML 文件。

### Just Build

`build` 用法：

```go
pugo build --source="source" --dest="dest" --theme="theme/default" --watch --debug
```

`--source` 设置内容目录，默认是 `source`。

`--dest` 设置编译内容保存的目录, 默认 `dest`。

`--theme` 设置主题模板的目录， 默认 `theme/default` ( PuGo 在 `theme` 文件夹提供 3 个主题 )。

`--watch` 开启文件变化监测。如果发生变化，立刻重新编译最新内容。

`--debug` 打印更多调试信息。

