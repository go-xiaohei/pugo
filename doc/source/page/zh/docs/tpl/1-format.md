```toml
title = "基本语法"
date = "2016-02-04 15:00:00"
slug = "zh/docs/tpl/syntax"
hover = "docs"
lang = "zh"
template = "docs.html"
```

`PuGo` 使用 `Go` 的模板语法。 你可以阅读 [text/template](#) 和 [html/template](#) 学习 Go 的模板语法。这里写一些需使用注意：

**作用域**

使用 `{{if}}`,`{{range}}`,`{{with}}` 时, 作用域会切换到当前变量，而不是全局。所以你需要 `{{$}}` 来读取当前作用域外的数据。

```html
<ul>{{range .List}}
    <li>item:{{.}} - from {{$.ListName}}</li>
{{end}}</ul>
```