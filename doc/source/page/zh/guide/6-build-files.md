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

### 监听变化

`PuGo` 可以监听内容和模板的变化，并立即重新编译最新内容。这将会覆盖所有生成的 HTML，并根据 md5 值判断是否需要更新静态文件。

```bash
pugo build --watch
```

### 自定义内容

从自定义目录读取内容编译：

```bash
pugo build --source="your-source"
```

### 自定义编译目录

编译到自定义的目录去：

```bash
pugo build --dest="your-directory"
```

### 自定义主题目录

使用特定的主题编译内容：

```bash
pugo build --theme="your-theme"
```
