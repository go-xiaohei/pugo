```ini

; post title, required
title = Welcome

; post slug, use to build permalink and url, required
slug = welcome

; post description, show in header meta
desc = welcome to try pugo site generator

; post created time, support
; 2015-11-28, 2015-11-28 12:28, 2015-11-28 12:28:38
date = 2015-12-20 12:20:20

; post updated time, optional
; if null, use created time
update_date = 2015-12-20 12:30:30

; author identifier, reference to meta.md [author.pugo], required
author = pugo

; thumbnails to the post
thumb = @media/golang.png

; tags, optional
tags = pugo
```

When you read the blog, `PuGo` is running successfully. Then enjoy your blog journey.

This blog is generated from file `source/welcome.md`. You can learn it and try to write your own blog article with following guide together.

#### blog meta

Blog's meta data, such as title, author, are created by first `ini` section with block **\`\`\`ini ..... \`\`\`**:

```ini
; post title, required
title = Welcome

; post slug, use to build permalink and url, required
slug = welcome

; post created time, support
; 2015-11-28, 2015-11-28 12:28, 2015-11-28 12:28:38
date = 2015-11-28 11:28

; post updated time, optional
; if null, use created time
update_date = 2015-11-28 12:28

; post description, show in header meta
desc = welcome to try pugo.static site generator

; author identifier, reference to meta.md [author.pugo], required
author = pugo-robot

; tags, optional
tags = pugo,welcome
```

#### blog content

Content are from the second section as `markdown` format:

```markdown
When you read the blog, `pugo` is running successfully. Then enjoy your blog journey.

This blog is generated from file `source/welcome.md`. You can learn it and try to write your own blog article with following guide together.

...... (markdown content)
```

Just write content after blog meta, all words will be parsed as markdown content.
