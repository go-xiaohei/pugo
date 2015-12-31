```ini
title = About PuGo
slug = about
desc = some words about pugo
date = 2015-12-24
author = pugo
; set nav to active status when this page
hover = about
; set template file to render this page
template =

[meta]
metadata = this is meta data
```

### Introduction

`PuGo` is a simple static site generator by [Golang](https://golang.org). It compiled [markdown](https://help.github.com/articles/markdown-basics/) content to site pages with beautiful theme. No dependencies, cross platform and very fast.

![golang](@media/golang.png)

### Quick start

1. Download and extract from [PuGo Releases](http://pugo.io), Run `pugo[.exe]` server directly.
2. Open `http://localhost:9899` to visit.


### Commands

Run a command when run `pugo` executable file:

- `pugo new` create new site, post or page.
- `pugo build` build static files.
- `pugo server` build and serve static files.

More details in [Commands](http://pugo.io/docs/command.html).

### Writing

`PuGo` support two kinds of content, `post` and `page`. you can create any `.md` file in proper directories in `source` directory. Read the [wiki](http://pugo.io/docs/write.html) to learn the layout format and more details.

### Publish

After you change your source `.md` files, run

```bash
$ ./pugo build
```

To build static files.

If you want to deploy your site, read this [documentation](http://pugo.io/docs/deploy-platform.html).

### Customize

- Theme: theme documentation is [Here](http://pugo.io/docs/theme.html)
