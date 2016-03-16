```toml
title = "Format"
date = "2016-02-04 15:00:00"
slug = "en/docs/cnt/format"
hover = "docs"
lang = "en"
template = "docs.html"
```

`PuGo` support `toml` and `ini` as meta and front-matter in post or page. In [Guide](#), it use `toml` as default. You can read them to write correct fields. This page explains differences in other formats.

### INI

`ini` is familiar to `toml` but more simple. So it need more words to describe some data.

##### Navigation

In `meta.ini`:

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

##### Author

Same to navigation data:

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
