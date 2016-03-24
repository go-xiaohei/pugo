```toml
title = "New"
date = "2016-02-04 15:00:00"
slug = "zh/docs/cmd/new"
hover = "docs"
lang = "zh"
template = "docs.html"
```

`new` 命令可以创建新的 `site`( 站点 ), `post`( 文章 ) 或 `page`( 页面 )。

### 站点

`PuGo` 内嵌有静态资源，因而可以直接新建站点，不需要下载任何附加内容：

```go
pugo new site
```

释放好的新站点内有默认的配置，`source` 文件夹有起始文章和页面，`theme` 文件夹有三个主题。

```go
pugo new site --doc
```

当设置 `--doc`，将只释放 `doc` 目录带有所有的文档资源。你可以使用 [Doc](/zh/docs/cmd/doc) 命令浏览文档内容。


### 文章

直接创建文章：

```go
pugo new post
```

空文章会生成在 `source/[year]/[day-month-hour-minute-second].md` 文件。

创建带标题的文章：

```go
pugo new post "this is new post"
```

新文章保存在 `source/this-is-new-post.md` 文件。

### 页面

和 `new post` 操作一样：

```go
pugo new page
pugo new page "this is new page"
```