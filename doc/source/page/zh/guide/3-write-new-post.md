```toml
title = "新建文章"
date = "2016-02-05 15:00:00"
slug = "zh/guide/write-new-post"
hover = "guide"
lang = "zh"
template = "guide.html"
```

创建新的文章：

```bash
pugo new post <title>
```

新的文章会创建在文件 `source/post/[year]/<title>.md`，比如:

```bash
pugo new post "this is a new post"
```

文件是 `source/post/2016/this-is-a-new-post.md`.

### Front-Matter

`Front-Matter` 描述文章的相关内容，可以用 markdown 中的 `toml` 代码块定义：

    ```toml
    title = "this is a new post"
    ```

所有的字段有：

```toml
# 文章的标题
title = "this is a new post"

# 文章的固定连接，必须填，用于生成访问的 URL
slug = "this-is-a-new-post"

# 文章的创建时间
date = "2016-02-05 15:00:00"

# 文章的更新时间；如果不填，等于创建时间
update_date = "2016-02-05 16:00:00"

# 文章的作者名。对应作者的详细信息会在 meta.toml 寻找。
author = "pugo"

# 文章的标签，可不填
tags = ["tag1","tag2"]

# 文章缩略图
# @media 会生成基于 source/media 的 URL
thumb = "@media/post-1.png"
```

`Front-Matter` 也支持`ini`格式的内容。

### 文章内容

文章内容采用 `markdown` 格式，直接写在 `toml` 代码块之后：

```md
# title

the post content

[link](/)
```


### 主题模板

文章使用模板 `theme/default/post.html`。您可以更新模板满足个性化需求。