```ini
title = Commands
slug = docs/commands
date = 2015-11-15
update_date = 2015-12-20
author = fuxiaohei
author_email = fuxiaohei@vip.qq.com
hover = docs
template =

[meta]
Source = "https://github.com/go-xiaohei/pugo-io/blob/master/source/page/prolog/command.md"
Version = ">=0.8.5"
```

`Pugo.Static` provides `new`, `build` , `server` and `deploy`.

### New

`new` command can create new `site`,empty `post` or `page` with pre-filled data.

    $ ./pugo new site

Extract bundle assets to init site with default data in **current** directory.

If you download `zip` from releases, the assets are included. `new site` extracts same assets.

    $ ./pugo new post
    $ ./pugo new page

It creates a empty post or page `.md` file named as `2015-11-30-12-20.md` as time-string format. You can fill content and rename it following your idea.

The created page's default template is `page.html`.

### Build

`build` command provides more options to build static files.

    $ ./pugo build
    $ ./pugo build --dest="dest" --theme="default" --nowatch --debug

Set `--dest` to make the built files in given directory.

Use `--theme` to change the theme used to build.

If add `--nowatch`, it just builds and exist, not watches changes to rebuild.

`--debug` prints more logs.

### Server

`server` command starts building and serving together.

    $ ./pugo server
    $ ./pugo server --dest="dest" --theme="default" --addr="0.0.0.0:9899" --debug

The usage of `--dest`, `--theme` and `--debug` are same to `build` command.

`--addr` changes the http server address, default value is `0.0.0.0:9899`.

The `server` command always watches changes of sources and templates and rebuild it immediately.

### Deploy

`deploy` command deploys files to other platform after building.

    $ ./pugo deploy
    $ ./pugo deploy --theme="default" --watch --debug

The usage of `--watch`, `--theme` and `--debug` are same to `build` command.

Read more in [Deploy](/docs/deploy).