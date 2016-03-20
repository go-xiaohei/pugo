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
pugo server --addr="0.0.0.0:9899" --source="source" --dest="dest" --theme="theme/default" --static --debug
```

`--addr` set the address and port that http server listen on, default is `0.0.0.0:9899`

`--source`, `--dest` and `--theme` set source, destination and theme directory, same to `build` command.

`--static` serve dest static files, but need correct `source` to load

`--debug` print more logs when running command.

### Notice

When `server` runs, `PuGo` builds contents immediately, then start http server. At the same time, `PuGo` watches file changes to rebuild soon.

So `server` command is better when developing or writing new contents. You can preview the new post or page. **But I recommend to use for public with --static flag**.

It's better to use web server to serve static files after building website.