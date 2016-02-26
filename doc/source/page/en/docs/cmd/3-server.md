```toml
title = "Server"
date = "2016-02-04 15:00:00"
slug = "en/docs/cmd/server"
hover = "docs"
lang = "en"
template = "docs.html"
```

`server` starts a HTTP server to display website.

```go
pugo server [--addr="0.0.0.0:9899" --from="source" --to="public" --theme="theme/default"]
```

`--addr` set the address and port that http server listen on, default is `0.0.0.0:9899`

`--from` and `--to` set source and destination directory, same to `build` command. **But these flags are not supported to migrate and deploy operations**.

`--theme` set the theme to build

### Notice

When `server` runs, `PuGo` builds contents immediately, then start http server. At the same time, `PuGo` watches file changes to rebuild soon.

So `server` command is better when developing or writing new contents. You can preview the new post or page. **But I don't recommend to use for public**.

Please use web server to serve static files after building website.