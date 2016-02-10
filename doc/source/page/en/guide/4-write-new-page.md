```toml
title = "Write New Page"
date = "2016-02-05 15:00:00"
slug = "en/guide/write-new-page"
hover = "guide"
lang = "en"
template = "guide.html"
```

To create a new page, run command:

```bash
pugo new page <title>
```
The new post file is created in `source/page/<title>.md`. For example:

```bash
pugo new page "this is a new page"
```

Result is `source/page/this-is-a-new-post.md`.

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
# Optional, it uses relative path that based on [source/page]
slug = "this-is-a-new-post"

# post's creating date
date = "2016-02-05 15:00:00"

# post's updating date, optional. If not set ,use creating date
update_date = "2016-02-05 16:00:00"

# author's name. It finds the proper author by name in meta.toml
author = "pugo"

# language environment. Use correct i18n tool in this page
lang = "en"

# use specific template file
# Otherwise, use `page.html` as default
template = "page.html"

# hove class set global hover class to render this page,
# read following "Hover" section
hover = "hover"

# meta set extra data
# read following "Meta" section
[meta]
data1 = "data1"
data2 = "data2" 
```

### Content

The post's content is `markdown` format. So just write it after `toml` block.

```md

# title

the post content

[link](/)

```

### URL

```toml
slug = "this-is-a-new-post"
```

`slug` is optional item. If you don's set `slug` value, use relative path to create link. 

For example, a page `source/page/new/this-is-page.md` without `slug` item. It will be built at `http://{domain}/new/this-is-page.html`.


### Hover

`hover = "hover"` set global `Hover` var in template when rendering this page. So you can set hover status of navigation, such as :

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

You can add extra string data in a page (**only support string value**). So you can do some custom changes for page view. Such as :

```html
<aside>
    {{range $k,$v := .Page.Meta}}
    <p><span>{{.$k}}:</span>{{.$v}}</p>
    {{end}}
</aside>
```


### Theme

All pages use template `theme/default/page.html`. If you need custom style, set `template = "my_page.html"` and modify `my_page.html` to display page data.