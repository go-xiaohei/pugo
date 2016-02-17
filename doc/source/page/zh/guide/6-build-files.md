```toml
title = "编译站点"
date = "2016-02-05 15:00:00"
slug = "zh/guide/build-files"
hover = "guide"
lang = "zh"
template = "guide.html"
```

`build` 命令可以编译站点：

```bash
pugo build 
```

For example:

![guide-3-build.jpeg](@media/guide-3-build.jpeg)

### 监听变化

`PuGo` 可以监听内容和模板的变化，并立即重新编译最新内容。这将会覆盖所有生成的 HTML，并根据 md5 值判断是否需要更新静态文件。

```bash
pugo build --watch
```

### 自定义内容和编译文件夹

从自定义目录读取内容编译：

```bash
pugo build --from="your-source"
```

编译到自定义的目录去：

```bash
pugo build --to="your-directory"
```

### 迁移和部署

`build` 命令支持直接进行 `迁移` 和 `部署`。

```bash
pugo build --from="rss+http://source-of-rss.xml"
pugo build --to="git://local-git-repository-directory"
```

更多的内容请阅读 [迁移](#) & [部署](#) 的文档.