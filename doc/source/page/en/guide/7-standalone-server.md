```toml
title = "Standalone Server"
date = "2016-02-05 15:00:00"
slug = "en/guide/standalone-server"
hover = "guide"
lang = "en"
template = "guide.html"
sort = 7
```

Use `server` command to run built-in server in `PuGo`:

```bash
pugo server --addr=0.0.0.0:9899
```

It listens `http://0.0.0.0:9899`. `--addr` can change the address and port that listens on.

It builds the source code right now, then start http server on destination directory as static file server. 

So you can set `--source`, `--dest` and `--theme` as `build` command.

```bash
pugo server --source="your-source" --dest="your-destination"
```

### Static

If just serve static files, use `--static` flag:

```bash
pugo server --source="source" --static
```

But It need read source to set correct url for http server.