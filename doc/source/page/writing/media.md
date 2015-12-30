```ini
title = Media
slug = docs/media
date = 2015-11-11
update_date = 2015-12-31
author = fuxiaohei
author_url = http://fuxiaohei.me/
hover = docs
template =

[meta]
Source = "https://github.com/go-xiaohei/pugo/blob/master/doc/source/page/writing/media.md"
Version = "0.9.0"
```

Media files are stored in `source/media` directory. They are copied into `{dest}/static/media` directory.

You can use placeholder `@static` and `@media` to fill correct url in post or page,

```markdown

this is a image, ![image](@media/image.png)

```

and some meta values.

```ini

; ----- in post
thumb = @media/post-thumbnail.png

; ----- in meta
[author.name]
avatar = @media/author-avatar.png

[meta]
cover = @media/cover.png

```