```ini
title = Setup
slug = docs/setup
date = 2015-11-11
update_date = 2015-12-31
author = fuxiaohei
author_url = http://fuxiaohei.me/
hover = docs
template =

[meta]
Source = "https://github.com/go-xiaohei/pugo/blob/master/doc/source/page/prolog/setup.md"
Version = "0.9.0"
```

# Install

You can install `PuGo` from **binary** file or **source code**.

**# Binary**

Binary files are released in [github release](https://github.com/go-xiaohei/pugo/releases) that marked as `Alpha` or `Stable` or [homepage](/) . Download the `zip` file and extract, it contains:

- **pugo[.exe]** - executable file, just run it
- **source** - where put your writing contents
- **template** - save template files
- **doc** - documentation assets, use for `pugo doc` command

Then just run, `./pugo[.exe]`.

**# Source Code**

Install from source code requires you have setup `Go` environment and set `$GOPATH` variable correctly. It needs **Go 1.4+**.

Then just `go get` the github repository:

    go get -u github.com/go-xiaohei/pugo

Build it in manual:

```bash
$ go generate -x
$ go build -v pugo.go
```

Then use following command to check which version of `PuGo` is installed on your system (Suppose `$GOPATH/bin` has been added to your `$PATH`):

```bash
$ pugo[.exe] -v
pugo version 0.9.(2016-01-01)
```

This file contains assets files. You can run it by **one file**.

**# One-File Install**

If no assets, only `pugo[.exe]` file can install a whole site. run :

    $ ./pugo[.exe] new site

It extracts binding assets to current directory.

Read [the doc](/docs/commands) to get more information.

# Getting Started

Let's explore the `PuGo` project.

- **source** - put your markdown contents into the directory, support subdirectories.
- **source/post** - posts are in the directory, support subdirectories.
- **source/page** - pages are in the directory, support subdirectories.
- **source/media** - media files usd in posts and pages
- **template** - put themes in the directory
- **template/default** - the default theme of `PuGo`, read more in [Template & Theme](/docs/templates).

After you prepare contents and templates well, just run `./pugo[.exe]`, it build and serve site on `0.0.0.0:9899`. It watches the `md` files in **`source`** directory. If you change something, it rebuilds contents to serve new files.

#### Run it up

Because of the assets extracted from zip, you can run the default server directly.

    $ ./pugo.[exe] server

Then preview `http://localhost:9899`.

