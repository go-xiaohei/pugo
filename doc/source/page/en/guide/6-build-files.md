```toml
title = "Build Files"
date = "2016-02-05 15:00:00"
slug = "en/guide/build-files"
hover = "guide"
lang = "en"
template = "guide.html"
```

Use `build` command to build files:

```bash
pugo build 
```

For example:

![guide-3-build.jpeg](@media/guide-3-build.jpeg)

### Watch

`PuGo` can watch changes and re-build files immediately. It overwrites any html files and checks md5sum to replace static files that needed.

```bash
pugo build --watch
```

### Custom Source and Destination

Build files from custom source:

```bash
pugo build --from="your-source"
```

Build files to custom destination:

```bash
pugo build --to="your-directory"
```

### Migrate and Deploy

`build` command do `migrate` and `deploy` together.

```bash
pugo build --from="rss+http://source-of-rss.xml"
pugo build --to="git://local-git-repository-directory"
```

Read more details in [Migrate](#) & [Deploy](#) documentations.