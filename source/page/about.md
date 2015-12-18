```ini
title = About Pugo.Static
slug = about-pugo
desc = some words about pugo.static
date = 2015-12-20
author = pugo
; set nav to active status when this page
hover = about
; set template file to render this page
template =

[meta]
metadata = this is meta data
```

### Introduction

`Pugo` is a simple static site generator by [Golang](https://golang.org). It compiled [markdown](https://help.github.com/articles/markdown-basics/) content to site pages with beautiful theme. No dependencies, cross platform and very fast.

![golang](/static/media/golang.png)

### Quick start

1. Download and extract from [Pugo Releases](http://pugo.io), Run `pugo[.exe]` directly.
2. Open `http://localhost:9899` to visit.


### Configuration

Add flags when run `pugo` executable file:

- `--addr=0.0.0.0:9999` set http server address
- `--theme=abc` set theme directory in template dir
- `--debug` set debug mode to print more information when running

### Writing

`Pugo` support two kinds of content, `post` and `page`. you can create any `.md` file in proper directories in `source` directory. Read the [wiki](#) to learn the layout format and more details.

### Publish

After you change your source `.md` files, just restart the server.

### Customize

- Theme: theme documentation is [wiki](http://pugo.io/docs/templates.html)
