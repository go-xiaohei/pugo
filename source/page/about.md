-----ini
title = About Pugo Static
slug = about
date = 2015-11-12
author = pugo-robot
author_email =
; set nav to active status when this page
hover = about
; set template file to render this page
template =


-----markdown
### Introduction

`Pugo` is a simple static site generator by [Golang](https://golang.org). It compiled [markdown](#) content to site pages with beautiful theme. No dependencies, cross platform and very fast.

### Quick start

1. Download and extract from [Pugo Releases](#), Run `pugo[.exe]` directly.
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

- Theme: theme documentation is [wiki](#)

- Build from source: install and debug documentation is [wiki]($)