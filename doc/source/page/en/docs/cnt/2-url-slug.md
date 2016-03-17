```toml
title = "URL and Slug"
date = "2016-02-04 15:00:00"
slug = "en/docs/cnt/url-slug"
hover = "docs"
lang = "en"
template = "docs.html"
```

`PuGo` generates URL by `slug` in default.

### Post

Post's URL is created by certain layout as `/[year]/[month]/[day]/[slug].html`. If no `slug` for a post, `PuGo` uses title to generate URL:

```toml
title = "URL and Slug"
# The URL is:
# /2016/3/15/URL-and-Slug.html
```

The `Permalink` of a post is from its URL without suffix string:

    URL         : /2016/3/15/URL-and-Slug.html
    Permalink   : /2016/3/15/URL-and-Slug

### Page

Page's URL is created by relative path or slug. When slug is empty, use relative path.

```toml
# file : source/page/about/me.md
title = "URL and Slug"
# The URL is:
# /about/me.html
```

The `Permalink` of a page is from its URL without suffix string:

    URL         : /about/me.html
    Permalink   : /about/me
