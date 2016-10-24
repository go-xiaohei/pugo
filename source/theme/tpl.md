### template documentation

PuGo provides three themes in this directory.

use specific theme:

  $ ./pugo build --theme="theme/default"

##### theme structure

**including** files:

- meta.html : title, meta, style , etc in `<header>`
- header.html : title, navigator at the top of page
- footer.html : script, copyright at the bottom of page
- comment.html : comment list and form in `post.html` and `page.html`

use `go template syntax` ---- `{{template "meta.html" .}}` to import it.

**major page** files:

- post.html - single post page
- posts.html - posts list page
- archive.html - posts archive page
- page.html - single page template as default. page can set template in meta with `template=xxx.html`
- index.html - homepage template, if not exist, use `posts.html`
