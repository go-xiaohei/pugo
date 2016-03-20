```toml
title = "数据格式"
date = "2016-02-04 15:00:00"
slug = "zh/docs/cnt/format"
hover = "docs"
lang = "zh"
template = "docs.html"
```

`PuGo` 支持 `toml` 和 `ini` 作为 meta 和文章或页面的 front-matter。 在 [入门](/zh/guide) 中,使用 `toml` 进行介绍。 这里对别的数据格式使用时的区别进行说明。

### INI

`ini` 类似 `toml` 但更简单，因而你需要更多的信息描述内容：

##### 导航

在 `meta.ini` 文件中：

```ini
[nav]
- : home
- : about

[nav.home]
; navigator to the link
link = "/"
; link title to fill text href element
title = "Home"
; i18n key if load i18n translation
i18n = "home"
; hover class to test whether is active of this navigation item
hover = "home"
; if blank is true, it forces browser to open new tab to display the linked page
blank = true

[nav.about]
......
```

##### 作者

和”导航”的数据格式一样：

```ini
[author]
- : pugo

[author.pugo]
; author'name, must be unique
name = "pugo"
; author's email, please be private as possible
email = ""
; author's link
url = "http://pugo.io"
; author's avatar, optional. If empty, generate Gravatar image by email
avatar = ""
; author's profile 
bio = "the robot of pugo, who generates all default contents."
```
