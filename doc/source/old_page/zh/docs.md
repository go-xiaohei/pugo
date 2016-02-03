```ini
title = 文档
slug = docs
date = 2015-11-11
update_date = 2015-12-20
author = fuxiaohei
author_url ＝ http://fuxiaohei.me/
hover = docs
lang = zh
template =

[meta]
Source = "https://github.com/go-xiaohei/pugo-io/blob/master/doc/source/page/docs.md"
Version = "0.9.0"
```

`PuGo` is a simple static site generator by [Golang](https://golang.org). It compiled [markdown](https://help.github.com/articles/markdown-basics/) content to site pages with beautiful theme. No dependencies, cross platform and very fast.

### Why create it

I wanted to write a dynamic blog engine with golang. But because of the lack of golang, it can't make blog extensible and scalable, such as plugin system. And site generators are growing in trend, for example, `Jeklly`, `Hexo`. So I try to write a site generator.

After days working, `PuGo` has completed major features:

- build and serve markdown file as posts or pages
- hot re-build when source file changes
- design and write a beautiful default theme
- basic third-party comment support

### Development

`PuGo` keep developing now, but it's available to use in production with following tips:

- please use release version, not master branch
- if upgrade, read release-note carefully to migrate data if needed

### Cases

- this site, hahaha
- [fuxiaohei.me](http://fuxiaohei.me) - the author
- [wuwen.org](http://wuwen.org/) - the leader of `Gogs`
- [lunny.info](http://lunny.info) - the author of `xorm`

