```ini
title = Meta, Navigation & Comment
slug = docs/meta
date = 2015-11-11
update_date = 2015-12-20
author = fuxiaohei
author_url = http://fuxiaohei.me/
hover = docs
template =

[meta]
Source = "https://github.com/go-xiaohei/pugo-io/blob/master/source/page/writing/meta.md"
Version = ">=0.8.5"
```

The site information are saved in file `meta.md` in source directory including site meta, navigation.

### Meta

```ini
; must be in [meta] block
[meta]
title = Pugo.Static
subtitle = generator
; print in html <meta>
keyword = pugo,golang,static,site,generator
; print in html <meta>
desc = pugo is a simple static site generator
; build links for feed, sitemap
domain = pugo.io
; root path for site, if empty, build as http://{domain}/
root = http://pugo.io/
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
link = https://github.com/go-xiaohei/pugo-static
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
