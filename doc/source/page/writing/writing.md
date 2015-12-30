```ini
title = Writing
slug = docs/writing
date = 2015-11-11
update_date = 2015-12-30
author = fuxiaohei
author_url = http://fuxiaohei.me/
hover = docs
template =

[meta]
Source = "https://github.com/go-xiaohei/pugo/blob/master/doc/source/page/writing/writing.md"
Version = "0.9.0"
```

The format is from `github` markdown preview page. If in this format, you can read the `.md` file more clear in github.There must be two sections in this format: **meta info** and **body**

### Meta

Meta info contains all values that same to `meta info` in first format, but block quoted by markdown block syntax "\`\`\`ini ...... ```".

```ini
; the title
title = Welcome to PuGo

; unique link
slug = welcome-pugo

; write a sentence to describe the post
desc = welcome to pugo

; created date,
; support minute time string
date = 2015-11-08 22:11
; support date time string, too
date = 2015-11-08

; if not set update_date, update time is same to date above
update_date = 2015-12-30

; author name, email and link
; if author value can be find meta.ini, use that author's data
author = pugo-robot
author_email =
author_url =

; thumbnails to the post or page
; @media to dest media url
; @static to dest static url
thumb = @media/golang.png

; ===== use in post
; post tags
tags = pugo,content,write

; ===== use in page
; set template to page
template = page.html

; set hover item in navigator in this page
hover = write-post

; add meta data to the post
[meta]
version = "0.9.0"
```

### Body

When ini block ends, extra contents are body for this markdown content totally,

You can visit this [example](https://github.com/go-xiaohei/pugo/blob/master/doc/source/page/writing/write.md) to learn the format.

-----

You can create an empty post or page sample by :

    $ ./pugo new post
    $ ./pugo new page

