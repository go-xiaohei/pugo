```toml
title = "下载和编译"
date = "2016-02-04 15:00:00"
slug = "zh/docs/dl-compile"
hover = "docs"
lang = "zh"
template = "docs.html"
```

安装`PuGo`之前，你需要从 [golang.org](https://golang.org) 下载 Go 语言。

`PuGo` 需要 **Go 1.4+**.

Go 安装好并功能正常后，获取 `PuGo` 的源码：

```go
go get github.com/go-xiaohei/pugo
```

然后你可以修改代码，编译自己的版本。

### 静态内容

`PuGo` 使用 [go-bindata](https://github.com/jteeuwen/go-bindata) 将静态资源嵌入源代码中。

```go
go get github.com/jteeuwen/go-bindata
```

添加 `$GOPATH/bin` 到系统 `PATH`, 然后:

```go
go generate -x
```

默认会编码 `doc`, `source`, `theme` 文件夹的内容到 `app/asset/asset.go` 文件。

### 注意

- `PuGo` 依赖 `fsnotify` 监听文件修改，但它并没有在 **arm** 平台充分测试.
- `PuGo` 使用 `log15.v2` 打印日志，但是该库无法在 **openbsd** 和 **netbsd** 平台编译。