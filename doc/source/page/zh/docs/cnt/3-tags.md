```toml
title = "文章标签"
date = "2016-02-04 15:00:00"
slug = "zh/docs/cnt/tags"
hover = "docs"
lang = "zh"
template = "docs.html"
```

在文章的 front-matter 中可以设置标签：

```toml
tags = ["pugo","golang"]
```

或 ini 格式：

```ini
; 使用逗号分割
tags = "pugo,golang"
```

`PuGo` 按照标签生成对应的文章列表页面，按照时间新旧排序，不做分页处理。

