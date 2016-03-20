```toml
title = "Doc"
date = "2016-02-04 15:00:00"
slug = "en/docs/cmd/doc"
hover = "docs"
lang = "en"
template = "docs.html"
```

`doc` command run documentation site on local:

```go
pugo new site
pugo doc --addr="0.0.0.0:9899"
```

After `new site`, it extracts `doc` data together. So you can run `doc` to compile them as documentation website.

`--addr` set the address and port that http server listen on, default is `0.0.0.0:9899`