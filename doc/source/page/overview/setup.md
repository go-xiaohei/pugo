```ini
title = Setup
slug = docs/setup
date = 2015-11-11
update_date = 2015-12-20
author = fuxiaohei
author_url = http://fuxiaohei.me/
hover = docs
template =

[meta]
Source = "https://github.com/go-xiaohei/pugo-io/blob/master/source/page/prolog/setup.md"
Version = ">=0.8.5"
```

# Install

You can install `Pugo.Static` from **binary** file or **source code**.

##### Binary

Binary files are released in [github release](https://github.com/go-xiaohei/pugo-static/releases) that marked as `Alpha` or `Stable` or [homepage](/) . Download the `zip` file and extract, it contains:

- **pugo[.exe]** - executable file, just run it
- **source** - where put your writing contents
- **template** - save template files

Then just run, `./pugo[.exe]`.

**One-File Install**

If no assets, only `pugo[.exe]` file can install a whole site. run :

    $ ./pugo[.exe] new site

It extracts binding assets to current directory.

Read [the doc](/docs/commands) to get more information.

##### Source Code

Install from source code requires you have setup `Go` environment and set `$GOPATH` variable correctly. It needs **Go 1.4+**.

Then just `go get` the github repository:

    go get -u github.com/go-xiaohei/pugo-static

You can then use following command to check which version of `Pugo.Static` is installed on your system (Suppose `$GOPATH/bin` has been added to your `$PATH`):

```bash
$ pugo[.exe] -v
pugo version 0.8.5(2015-12-20)
```

If you want to build it in manual, go to the `pugo-static` directory and run:

```bash
$ go generate -x
$ go build -v pugo.go
```

This file contains assets files. You can run it by **one file**.

# Getting Started

Let's explore the `Pugo.Static` project.

- **conf.ini** - configuration file, read more details in [Configuration](/docs/config).
- **source** - put your markdown contents into the directory, support subdirectories.
- **source/post** - posts are in the directory, support subdirectories.
- **source/page** - pages are in the directory, support subdirectories.
- **source/media** - media files usd in posts and pages
- **template** - put themes in the directory
- **template/default** - the default theme of `Pugo.Static`, read more in [Template & Theme](/docs/templates).

After you prepare contents and templates well, just run `./pugo[.exe]`, it build and serve site on `0.0.0.0:9899`. It watches the `md` files in **`source`** directory. If you change something, it rebuilds contents to serve new files.

#### Run it up

Because of the assets extracted from zip, you can run the default server directly.

    $ ./pugo.[exe] server

Then preview `http://localhost:9899`.

