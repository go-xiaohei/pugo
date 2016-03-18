```toml
title = "Each Template"
date = "2016-02-04 15:00:00"
slug = "en/docs/tpl/each"
hover = "docs"
lang = "en"
template = "docs.html"
```

There are some different global variable in different page.

### post.html

`{{.Post}}` is data of one post.

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

`{{.Posts}}` is data of list of posts. `{{.Pager}}` is pagination tool to generate paged URL. If the list in one tag, `{{.Tag}}` is the tag data.

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

`{{.Page}}` is data of this page.

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

