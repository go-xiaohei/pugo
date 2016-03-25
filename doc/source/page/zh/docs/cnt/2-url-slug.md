```toml
title = "URL 和 Slug"
date = "2016-02-04 15:00:00"
slug = "zh/docs/cnt/url-slug"
hover = "docs"
lang = "zh"
template = "docs.html"
```

`PuGo` 依赖 `slug` 生成 URL。

### 文章

文章的 URL 按照格式 `/[year]/[month]/[day]/[slug].html` 生成。 如果没有设置 `slug` 使用文章标题生成 URL：

```toml
title = "URL and Slug"
# The URL is:
# /2016/3/15/URL-and-Slug.html
```

### 页面

页面的 URL 按照 slug 或文件的相对位置（没有设置slug时）生成。

```toml
# file : source/page/about/me.md
title = "URL and Slug"
# The URL is:
# /about/me.html
```
