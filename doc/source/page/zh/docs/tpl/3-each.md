```toml
title = "页面的独立数据"
date = "2016-02-04 15:00:00"
slug = "zh/docs/tpl/each"
hover = "docs"
lang = "zh"
template = "docs.html"
```

不同的页面有一些特有的数据：

### post.html

`{{.Post}}` 是当前的文章数据。

```html
<article class="article">
    <div class="row">
        <div class="col-md-10 col-md-offset-1 panel panel-default">
            <header class="header">
                <div class="meta">
                    <span class="date">
                        <span class="month">{{printf "%d" .Post.Created.Month}}</span>
                        <span class="day">{{.Post.Created.Day}}</span>
                    </span>
                </div>
                <h3 class="title">
                    <a href="{{.Post.URL}}">{{.Post.Title}}</a>
                </h3>
            </header>
            <aside class="aside clearfix">
                {{range .Post.Tags}}
                <a class="tag label label-info" href="{{.URL}}">{{.Name}}</a>
                {{end}}
                {{if .Post.Author}}
                <a class="stat label label-default pull-right"{{if .Post.Author.URL}} href="{{.Post.Author.URL}}" target="_blank"{{end}}>{{.Post.Author.Name}}</a>
                {{end}}
            </aside>
            <section class="brief">{{.Post.ContentHTML}}</section>
        </div>
    </div>
</article>
```

### posts.html

`{{.Posts}}` 是文章列表数据， `{{.Pager}}` 是分页的工具。 如果当前页面时标签下的文章列表， `{{.Tag}}` 是对应标签的数据。

```html
{{range .Posts}}
<article class="article">
......
</article>
{{end}}
<div class="article-pager text-center">
    {{if .Pager.Prev}}<a class="btn btn-lg btn-info" href="{{.Pager.PrevURL}}">{{.I18n.Tr "pager.prev"}}</a>{{end}}
    {{if .Pager.Next}}<a class="btn btn-lg btn-info" href="{{.Pager.NextURL}}">{{.I18n.Tr "pager.next"}}</a>{{end}}
</div>
```

### page.html

`{{.Page}}` 是当前页面的数据。

```html
<article class="article">
    <div class="row">
        <div class="col-md-10 col-md-offset-1 panel panel-default">
            <header class="header">
                <div class="meta">
                    <span class="date">
                        <span class="month">{{printf "%d" .Page.Created.Month}}</span>
                        <span class="day">{{.Page.Created.Day}}</span>
                    </span>
                </div>
                <h3 class="title">
                    <a href="{{.Page.URL}}">{{.Page.Title}}</a>
                </h3>
            </header>
            {{if .Page.Author}}
            <aside class="aside clearfix">
                <a class="stat label label-default pull-right"{{if .Page.Author.URL}} href="{{.Page.Author.URL}}" target="_blank"{{end}}>{{.Page.Author.Name}}</a>
            </aside>
            {{end}}
            <section class="brief">{{.Page.ContentHTML}}</section>
        </div>
    </div>
</article>
```

