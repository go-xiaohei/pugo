```ini
title = Deploy to Platform
slug = docs/deploy-platform
date = 2015-11-28
update_date = 2015-12-22
author = fuxiaohei
author_url = http://fuxiaohei.me/
hover = docs
template =

[meta]
Source = "https://github.com/go-xiaohei/pugo-io/blob/master/source/page/deploy/platform.md"
Version = ">=0.8.5"
```

`Pugo.Static` support to sync static files with third-party platforms, `git`, `ftp` and `sftp`.

### Git

You need clone the repository with `git://` or `https://{user}:{password@github.com/username/repo-name`. So that `Pugo` can push commits without password requirement.

Then checkout to your proper branch ( such as `gh-pages` for github pages).

```bash

$ git clone https://{user}:{password@github.com/username/repo-name
$ git checkout gh-pages
$ pugo --dest="git://repo-name?commit=your-commit-message"

```

`Pugo` builds static files to your `repo-name` and commit with optional message `commit`.


### Ftp

Upload to ftp server:

```bash

$ pugo --dest="ftp://{user}:{password}@{host}/{directory}
$ pugo --dest="ftp://admin:admin@127.0.0.1:2121/dest

```

The `directory` is based on your ftp home directory.

### Sftp

Upload via ssh:

```bash

$ pugo --dest="sftp://{user}:{password}@{host}/{directory}
$ pugo --dest="sftp://admin:admin@127.0.0.1:22/home/user/dest
$ pugo --dest="sftp://admin:admin@127.0.0.1:22/~/dest

```

Now SSH only supports user-password authorization.

The `directory` is based on root `/`. You need use `~` to set to user's home directory.