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

### Watch

`PuGo` can watch changes and re-build files immediately. It overwrites any html files and checks md5sum to replace static files that needed.

```bash
pugo build --watch
```

### Custom Source

Build files from custom source:

```bash
pugo build --source="your-source"
```

### Custom Destination

Build files to custom destination:

```bash
pugo build --dest="your-directory"
```

### Custom Theme

Build files with specific theme:

```bash
pugo build --theme="your-theme-directory"
```