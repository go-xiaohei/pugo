```toml
title = "Global Variables"
date = "2016-02-04 15:00:00"
slug = "en/docs/tpl/global"
hover = "docs"
lang = "en"
template = "docs.html"
```

`Global Variables` are assigned in all templates and can be used directly.

### Simple Data

`{{.Version}}` print PuGo's version number.

`{{.Root}}` print full url path, such as `http://pugo.io`, no ending slash.

`{{.Base}}` print base path. For example, when root path is `http://pugo.io/blog`, base path is `/blog` , no ending slash.

`{{.Hover}}` is current hover class in this page. You can use it to compare with hover class in navigation.

`{{.Title}}` is the title of this page, use in `<title>{{.Title}}</title>`. Default value is website name. In post or page, use its title.

`{{.Desc}}` is  the description of this page, use in `<meta>`.

`{{.PostType}}` is the type of this page, including "index","post","page","archive","post-list","post-tag" and "tag".

`{{.URL}}` and `{{.PermaKey}}` are the URL and perma keyword of this page.

### Global Object

`{{.Nav}}` is navigation from Meta, use to render global navbar.

```html
<ul id="nav-list">{{range .Nav}}
    <li class="{{if eq .Hover $.Hover}} hover{{end}}"><a href="{{.Link}}" class="link">{{.Title}}</a></li>{{end}}
</ul>
<ul id="nav-list-with-i18n">{{range .Nav}}
    <li class="{{if eq .Hover $.Hover}} hover{{end}}"><a href="{{.TrLink $.I18n}}" class="link">{{.Tr $.I18n}}</a></li>{{end}}
</ul>
```

`{{.Meta}}` is basic info from Meta, including Title, Subtitle, Keyword, Desc, Cover(cover image) and Language.

`{{.Comment}}` is comment option, including Disqus and Duoshuo.

`{{.Analytics}}` is analytics option, including Google and Baidu.

`{{.I18n}}` is i18n tool, use to render value to i18n value.

```html
{{.I18n.Tr "nav.item"}}
{{.I18n.Tr "post.readmore"}}
```