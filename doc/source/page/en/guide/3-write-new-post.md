```toml
title = "Write New Post"
date = "2016-02-05 15:00:00"
slug = "en/guide/write-new-post"
hover = "guide"
lang = "en"
template = "guide.html"
sort = 3
```

To create a new post, run command:

```bash
pugo new post <title>
```
The new post file is created in `source/post/[year]/<title>.md`. For example:

```bash
pugo new post "this is a new post"
```

Result is `source/post/2016/this-is-a-new-post.md`.

### Front-Matter

`Front-Matter` provides some details to describe your writting. Now it supports `toml` and writes as a code block in markdown.

    ```toml
    title = "this is a new post"
    ```

The full items:

```toml
# post's title
title = "this is a new post"

# post's slug link, use to generate visitable link
slug = "this-is-a-new-post"

# post's creating date
date = "2016-02-05 15:00:00"

# post's updating date, optional. If not set ,use creating date
update_date = "2016-02-05 16:00:00"

# author's name. It finds the proper author by name in meta.toml
author = "pugo"

# tags, optional, need be a array
tags = ["tag1","tag2"]

# thumbnail for this post, optional
# @media means finding the media in source/media directory
thumb = "@media/post-1.png"
```

**Front-Matter** support `ini` format.

### Content

The post's content is `markdown` format. So just write it after `toml` block.

```md
# title

the post content

[link](/)
```


### Theme

All posts use template `theme/default/post.html`. If you need custom style, try to modify it.