```toml
title = "公共变量"
date = "2016-02-04 15:00:00"
slug = "zh/docs/tpl/global"
hover = "docs"
lang = "zh"
template = "docs.html"
```

公共变量是所有页面都可以读取的数据。

### 简单数据

`{{.Version}}` PuGo 的版本号。

`{{.Root}}` 打印站点完整地址，如 `http://pugo.io`，没有最后的斜杠。

`{{.Base}}` 打印站点的 base 地址，比如如果完整地址是 `http://pugo.io/blog`, base 地址就是 `/blog`，没有最后的斜杠。

`{{.Hover}}` 是当前的 hover 值。 你可以用于和导航项目对比，判断导航是否是 hover 的。

`{{.Title}}` 当前页面的标题，是用在 `<title>{{.Title}}</title>`. 默认是 meta 的站点名称，或文章或页面的标题。

`{{.Desc}}` 是站点的介绍，用在 `<meta>`.

`{{.PostType}}` 是当前页面的类型，有 "index","post","page","archive","post-list","post-tag" 和 "tag"。

`{{.URL}}` 和 `{{.PermaKey}}` 是当前页面的 URL 和 永久关键字（不是永久连接）。

### 结构数据

`{{.Nav}}` 是全局的导航数据。

```html
<ul id="nav-list">{{range .Nav}}
    <li class="{{if eq .Hover $.Hover}} hover{{end}}"><a href="{{.Link}}" class="link">{{.Title}}</a></li>{{end}}
</ul>
<ul id="nav-list-with-i18n">{{range .Nav}}
    <li class="{{if eq .Hover $.Hover}} hover{{end}}"><a href="{{.TrLink $.I18n}}" class="link">{{.Tr $.I18n}}</a></li>{{end}}
</ul>
```

`{{.Meta}}` 是站点的基本数据，包括 Title, Subtitle, Keyword, Desc, Cover(cover image) 和 Language。

`{{.Comment}}` 是评论设置，包括 Disqus 和 Duoshuo。

`{{.Analytics}}` 是第三方统计设置， 包括 Google and Baidu。

`{{.I18n}}` 是 i18n 工具，用于打印不同语言的数值。

```html
{{.I18n.Tr "nav.item"}}
{{.I18n.Tr "post.readmore"}}
```