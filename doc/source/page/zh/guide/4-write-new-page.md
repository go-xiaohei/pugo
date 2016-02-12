```toml
title = "新建页面"
date = "2016-02-05 15:00:00"
slug = "zh/guide/write-new-page"
hover = "guide"
lang = "zh"
template = "guide.html"
```

添加新的页面：

```bash
pugo new page <title>
```
新的页面会创建在 `source/page/<title>.md` 文件，就像：

```bash
pugo new page "this is a new page"
```

文件是 `source/page/this-is-a-new-post.md`。

### Front-Matter

`Front-Matter` 描述页面的相关内容，可以用 markdown 中的 `toml` 代码块定义：

    ```toml
    title = "this is a new post"
    ```

页面的完整字段：

```toml
# 页面的标题
title = "this is a new post"

# 页面的链接，用于生成 URL
# 如果不填，将使用基于 [source/page] 目录的相对文件地址作为链接
slug = "this-is-a-new-post"

# 页面的创建时间
date = "2016-02-05 15:00:00"

# 页面的更新时间；如果不填，等于创建时间
update_date = "2016-02-05 16:00:00"

# 页面的作者名；会在 meta.toml 中寻找对应的作者
author = "pugo"

# 语言定义，使用对应的国际化内容渲染这个页面
lang = "en"

# 使用特定的模板文件
# 如果不填，用 `page.html`
template = "page.html"

# 给这个页面定义选定状态
# 用于响应导航的选定状态
hover = "hover"

# 添加额外数据
# 阅读以下 "Meta" 的内容
[meta]
data1 = "data1"
data2 = "data2" 
```

### 页面内容

页面内容采用 `markdown` 格式，直接写在 `toml` 代码块之后：

```md

# title

the post content

[link](/)

```

### URL

```toml
slug = "this-is-a-new-post"
```

`slug` 是选填字段。如果不填，会使用基于文件目录的相对地址生成 URL.

比如页面 `source/page/new/this-is-page.md` 不填 `slug`，生成的 URL 会是 `http://{domain}/new/this-is-page.html`.


### Hover

`hover = "hover"` 为整个页面设置了 `Hover` 变量。编译页面时可以根据 `Hover` 设置导航的状态 :

```html
{{if eq .Hover .Nav.Hover}}
    <a href="#">current link</a>
{{end}}
```

### Meta

```toml
[meta]
data1 = "data1"
data2 = "data2" 
```

您可以添加一些额外的一对对的字符串数据给页面，在模板中使用以满足一些自定义的功能:

```html
<aside>
    {{range $k,$v := .Page.Meta}}
    <p><span>{{.$k}}:</span>{{.$v}}</p>
    {{end}}
</aside>
```


### 主题模板

所有页面的主题模板是 `theme/default/page.html`. 使用 `template = "my_page.html"` 设置特定模板。