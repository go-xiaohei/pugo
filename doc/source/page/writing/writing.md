```ini
title = Writing
slug = docs/writing
date = 2015-11-11
update_date = 2015-12-20
author = fuxiaohei
author_url = http://fuxiaohei.me/
hover = docs
template =

[meta]
Source = "https://github.com/go-xiaohei/pugo-io/blob/master/source/page/writing/writing.md"
Version = ">=0.8.5"
```

The format is from `github` markdown preview page. If in this format, you can read the `.md` file more clear in github.There must be two sections in this format: **meta info** and **body**

### Meta Info

Meta info contains all values that same to `Base Info` in first format, but block quoted by markdown block syntax "```".

    ```ini
    ; the title
    title = Welcome to Pugo.Static

    ; unique link
    slug = welcome-pugo-static

    ; write a sentence to describe the post
    desc = welcome to pugo.static

    ; created date,
    ; support minute time string
    date = 2015-11-08 22:11
    ; support date time string, too
    date = 2015-11-08

    ; if not set update_date, update time is same to date above
    update_date = 2015-11-11

    ; author name, email and link
    author = pugo-robot
    author_email =
    author_url =

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
    version = ">=0.8.5"
    ```

### Body

When ini block ends, extra contents are body for this markdown content totally,

You can visit this [example](https://github.com/go-xiaohei/pugo-io/blob/master/source/page/writing/write.md) to learn the format.

-----

You can create an empty post or page sample by :

    $ ./pugo new post
    $ ./pugo new page

