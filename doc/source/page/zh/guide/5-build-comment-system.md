```toml
title = "添加评论功能"
date = "2016-02-05 15:00:00"
slug = "zh/guide/add-comment-system"
hover = "guide"
lang = "zh"
template = "guide.html"
```

`PuGo` 支持 [Disqus](#) 和 [多说](#) 评论系统，需要在 `meta.toml` 添加一下配置

```toml
[comment]
# disqus.com 评论
disqus = "test"

# duoshuo.com 评论
duoshuo = "test"
```

### 模板说明

嵌入评论的默认模板在 `embed/comment.html`。您可以参考它，实现更多评论平台的支持。

##### Disqus.com

`Disqus` 的通用代码如下：

```html
<div id="disqus_thread"></div>
<script>
  var disqus_shortname = 'site-name';
  var disqus_url = 'page-url';
  var disqus_title = "page-title";
  (function(){
    var dsq = document.createElement('script'); dsq.type = 'text/javascript'; dsq.async = true;
    dsq.src = 'https://go.disqus.com/embed.js';
    (document.getElementsByTagName('head')[0] || document.getElementsByTagName('body')[0]).appendChild(dsq);
  })();
</script>
```

只需要修改 `disqus = "site-name"`，就可以使用。

