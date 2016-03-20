```toml
title = "辅助函数"
date = "2016-02-04 15:00:00"
slug = "zh/docs/tpl/helper"
hover = "docs"
lang = "zh"
template = "docs.html"
```

`PuGo` 注册了一些模板的辅助函数：

`{{HTML "<p>html code</p>"}}` 直接打印 HTML 代码，避免转义。

`{{Include "file.html" .Data}}` 使用 Data 渲染特定的页面。

`{{url "link"}}` 使用 base 地址拼接 URL 如 '[base]/link`。

`{{fullUrl "link"}}` 使用完整地址拼接 URL 如 `http://[domain]/[base]/link`。