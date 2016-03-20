```toml
title = "Basic Syntax"
date = "2016-02-04 15:00:00"
slug = "en/docs/tpl/syntax"
hover = "docs"
lang = "en"
template = "docs.html"
```

`PuGo` template use `Go` template syntax. Read [text/template](https://golang.org/pkg/text/template/) and [html/template](https://golang.org/pkg/html/template/) to learn basic syntax. There are some tips to help:

**Global Scope**

When use `{{if}}`,`{{range}}`,`{{with}}` keyword, it changes scope to current variable, not global scope, so you need use `{{$}}` to read data out of current scope:

```html
<ul>{{range .List}}
    <li>item:{{.}} - from {{$.ListName}}</li>
{{end}}</ul>
```