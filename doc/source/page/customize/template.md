```ini
title = Theme
slug = docs/templates
date = 2015-11-14
update_date = 2015-12-30
author = fuxiaohei
author_url = http://fuxiaohei.me/
hover = docs
template =

[meta]
Source = "https://github.com/go-xiaohei/pugo/blob/master/doc/source/page/customize/template.md"
Version = "0.9.0"
```

The default theme is in `/template/default`. You can create a new directory such as `xxx` to `/template/xxx`.

then use it via :

```bash
$./pugo build --theme=xxx
$./pugo server --theme=xxx
```

the files list:

- *index.html* - index page template
- *posts.html* - posts list template, used to `/posts/[page_number]` and `/` if index.html is missing
- *post.html* - single post template
- *page.html* - single page template
- *archive.html* - archive template, used to `/archive`

embedded files:

- *comment.html* - comment template, embedded in post.html or page.html
- *meta.html* - meta template, embedded in other templates
- *header.html* - header template, embedded in other templates
- *footer.html* - footer template, embedded in other templates

#### Syntax

All templates use `Go Template` as default template engine, so you need learn something about it from `Go` language.
