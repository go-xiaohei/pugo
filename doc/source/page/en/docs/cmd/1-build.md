```toml
title = "Build"
date = "2016-02-04 15:00:00"
slug = "en/docs/cmd/build"
hover = "docs"
lang = "en"
template = "docs.html"
```

`build` command can not only build all contents to website html files, but also migrate other data to source and deploy built data to third party system.

### Just Build

`build` command basic usage:

```go
pugo build [--from="source"] [--to="public"] [--theme="theme/default"] [--watch]
```

`--from` set the source directory, default is `source`.

`--to` set the directory that PuGo builds contents to, default is `public`.

`--theme` set the directory of theme, default is `theme/default` ( PuGo provides 3 themes in `theme` ).

`--watch` set flag to watching changes and rebuild site.

### From & To

`--from` & `--to` provides rich supports to enable `migrate` and `deploy`, such as

```go
pugo build --from="rss+http://fuxiaohei.me/feed.xml"
pugo build --to="ftp://user:password@127.0.0.0.1:2121/pugo-data"
```

Those all meet **URL** scheme. More usages can be found in `Migrate` and `Deploy` documents.