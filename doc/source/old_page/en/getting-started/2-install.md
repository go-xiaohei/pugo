```ini
title = Install
slug = start/install
date = 2015-11-11
update_date = 2015-12-20
author = fuxiaohei
author_url ＝ http://fuxiaohei.me/
hover = start
lang = en
template = start.html

[meta]
Source = "https://github.com/go-xiaohei/pugo-io/blob/master/doc/source/page/docs.md"
Version = "0.9.0"
```

Binary file of `PuGo` is released on [http://pugo.io](http://pugo.io).

- Select one by your system and download:

![install-1](@media/s-2-install-1.png)

- Extract zip and create default site:

![install-2](@media/s-2-install-2.png)

    use command line tool (in Windows, `Shift + Right-Click`):

```bash
$ pugo new site
```

![install-3](@media/s-2-install-3.png)

    assets are extracted:

![install-4](@media/s-2-install-4.png)

- Run default site without any changes:

```bash
$ pugo server

// print log
INFO[01-13|11:54:32] Dest.dest
WARN[01-13|11:54:32] Dest.dest.Existed
INFO[01-13|11:54:32] Server.Start.0.0.0.0:9899
INFO[01-13|11:54:32] Lang.en
INFO[01-13|11:54:33] Build.Finish                             duration=14.0125ms count=1
INFO[01-13|11:54:33] Watch.Start

```

Visit `http://localhost:9899`, it's running.






