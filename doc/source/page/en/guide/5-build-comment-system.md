```toml
title = "Add Comment System"
date = "2016-02-05 15:00:00"
slug = "en/guide/add-comment-system"
hover = "guide"
lang = "en"
template = "guide.html"
```

`PuGo` supports [Disqus](#) and [Duoshuo](#) comment systems. Just config it in `meta.toml`.

```toml
[comment]
# disqus.com comment system
disqus = "test"

# duoshuo.com comment system
duoshuo = "test"
```

### Theme

the default comment embedded template are `embed/comment.html`. If you use different comment system, modify the theme file.

##### Disqus.com

`Disqus`'s universal comment html is like:

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

You should write `disqus = "site-name"` to make correct.

