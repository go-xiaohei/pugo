-----ini
title:About Gogs
slug:about-gogs
date:2015-11-10
author:傅小黑
author_email:fuxiaohei@vip.qq.com
tags:Gogs,Go,Git

-----markdown
# 什么是 Gogs?

Gogs (Go Git Service) 是一款极易搭建的自助 Git 服务。

## 开发目的

Gogs 的目标是打造一个最简单、最快速和最轻松的方式搭建自助 Git 服务。使用 Go 语言开发使得 Gogs 能够通过独立的二进制分发，并且支持 Go 语言支持的 **所有平台**，包括 Linux、Mac OS X、Windows 以及 ARM 平台。

## 开源组件

- Web 框架：[Macaron](http://go-macaron.com)
- UI 组件：
    - [Semantic UI](http://semantic-ui.com/)
    - [GitHub Octicons](https://octicons.github.com/)
    - [Font Awesome](http://fontawesome.io/)
- 前端插件：
    - [DropzoneJS](http://www.dropzonejs.com/)
    - [highlight.js](https://highlightjs.org/)
    - [clipboard.js](https://zenorocha.github.io/clipboard.js/)
    - [emojify.js](https://github.com/Ranks/emojify.js)
    - [jQuery Date Time Picker](https://github.com/xdan/datetimepicker)
    - [jQuery MiniColors](https://github.com/claviska/jquery-minicolors)
- ORM：[Xorm](https://github.com/go-xorm/xorm)
- 数据库驱动：
    - [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
    - [github.com/lib/pq](https://github.com/lib/pq)
    - [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
- 以及其它所有 Go 语言的第三方包依赖。