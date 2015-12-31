```ini
title = Site Meta
slug = docs/meta
date = 2015-11-11
update_date = 2015-12-30
author = fuxiaohei
author_url = http://fuxiaohei.me/
hover = docs
template =

[meta]
Source = "https://github.com/go-xiaohei/pugo/blob/master/doc/source/page/writing/meta.md"
Version = "0.9.0"
```

The site information are saved in file `meta.ini` in source directory including site meta, navigation.

### Meta

```ini
; site title, show in <title>
title = PuGo

; subtitle, words after title, in description
subtitle = Static Site Generator

; print in html <meta>
keyword = pugo,golang,static,site,generator

; print in html <meta>
desc = PuGo is a Simple Static Site Generator

; build links for feed, sitemap
domain = localhost

; root path for site; if empty, build as http://{domain}/
root = http://pugo.io/

; cover page for homepage
cover = @media/cover.jpg

; global language
lang = en
```

`meta` data can use in template via go template syntax:

```html
<meta name="keywords" content="{{.Meta.Keyword}}"/>
<meta name="description" content="{{.Meta.Desc}}"/>
```

Be careful about `root` value in meta. If your site is built in subdirectory, such as :

```ini
root = http://pugo.io/blog
```

You need use `{{.Root}}` to fix your url:

```html
<h1><a href="{{.Root}}/">homepage</a></h1>
<!-- now the href value is "/blog/" -->
```

### Navigation

```ini

[nav]
; reference to [nav.doc]
-:doc
-:github

[nav.doc]
; nav link
link = /docs

; text for href link
title = Documentation

; i18n keyword, not implemented
i18n = documentation

; set nav to active status
hover = docs

[nav.github]
link = https://github.com/go-xiaohei/pugo
title = Github
i18n = github

; icon class, used as <i>'s class
icon = "fa fa-github"

; open new tab or window to visit
blank = true

```

### Comment

```ini

[comment.disqus]
; site name of disqus
site = fuxiaohei

```

### Authors

```ini

; author data
[author]
-:pugo

[author.pugo]
name = pugo-robot
email =
url = http://pugo.io
avatar = @media/author.png
bio = the robot of pugo, who generates all default contents.

```

In post or page, you can use `author = pugo-robot` to reference this author.

### Build Settings

```ini

; ignore files to build or copy
[build.ignore]
-:CNAME
-:.git

```