```ini
title = Theme Function
slug = docs/template-func
date = 2015-12-30
update_date = 2015-12-30
author = fuxiaohei
author_url = http://fuxiaohei.me/
hover = docs
template =

[meta]
Source = "https://github.com/go-xiaohei/pugo/blob/master/doc/source/page/customize/template_func.md"
Version = "0.9.0"
```

There are some built-in functions in `PuGo` theme.

### HTML Content

Go will escape any string in template. So it need add-on functions to print html.

```html

<p>{{HTML "<span>html</span>"}}</p>

<!-- use bytes -->
<p>{{HTMLByte .BytesData}}</p>

```

### URL Builder

If you need print full url path:

```html

<a href="{{url .Post.Url}}">{{.Post.Title}}</a>
<!-- /base/{{.Post.Url}} -->

<a href="{{fullUrl .Post.Url}}">{{.Post.Title}}</a>
<!-- http://domain/base/{{.Post.Url}} -->

```

### International

- **{{.I18n}}** `*helper.I18n`  global I18n support

```html
<!-- read more or 阅读更多 -->
read more : {{.I18n.Tr "post.readmore"}}

<!-- read more about %s -->
read more : {{.I18n.Trf "post.readmore" .Post.Title}}

<!-- use html content -->
<p>{{.I18n.TrHTML "post.readmore"}}</p>
<!-- <p><a href="#">title</a></p> -->
```
