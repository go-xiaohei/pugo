```toml
title = "Helper Function"
date = "2016-02-04 15:00:00"
slug = "en/docs/tpl/helper"
hover = "docs"
lang = "en"
template = "docs.html"
```

`PuGo` adds some helper functions to easy usage.

`{{HTML "<p>html code</p>"}}` print html code.

`{{Include "file.html" .Data}}` embed included template with data.

`{{url "link"}}` print url with base path, as '[base]/link`.

`{{fullUrl "link"}}` print url with domain, as `http://[domain]/[base]/link`.