```ini
title = Theme Variables
slug = docs/template-vars
date = 2015-11-11
update_date = 2015-12-20
author = fuxiaohei
author_url = http://fuxiaohei.me/
hover = docs
template =

[meta]
Source = "https://github.com/go-xiaohei/pugo-io/blob/master/source/page/customize/template_var.md"
Version = ">=0.8.5"
```

### Global Vars

Global vars are registered in all pages:

- **{{.Nav}}**  `model.Navs`    Navigation items at the top of page, used in `header.html`.
- **{{.Meta}}** `*model.Meta`   Meta info from `meta.md`, used in `meta.html`.
- **{{.Title}}**    `string`  Title of this page, changing with current data. If a post ,use post title. If a page, use page title.
- **{{.Desc}}** `string` Description of this page, changing with current data. Same to title.
- **{{.Root}}** `string`    Root path for global url. If you set sub directory, need fix by this var.
- **{{.Version}}** `string` Pugo's version string
- **{{.Permalink}}** `string` Permalink of the page, existing in whole each pages.
- **{{.Comment}}**  `*model.Comment`    Comment settings from `comment.md`, used in `comment.html`.

### Post Vars

In Single post template `post.html`:

- **{{.Post}}** `*model.Post` the post data of this page.
- **{{.Post.Title}}** `string` the title of this post.
- **{{.Post.Slug}}** `string` the slug of this post, unique identifier.
- **{{.Post.Permalink}}** `string` the permalink of this post, same value to `{{.Permalink}}`, maybe it's not visitable sometimes.
- **{{.Post.Url}}** `string` the visitable of this post, maybe not same to permalink.
- **{{.Post.Created}}** `model.Time` the created time of this post.
- **{{.Post.Updated}}** `model.Time` the update time of this post. If no updates, be equal to created time.
- **{{.Post.Tags}}** `[]model.Tag` the tags of this post.
- **{{.Post.Author}}** `model.Author` the author of this post.
- **{{.Post.Raw}}** `[]byte` the raw content bytes of this post.
- **{{.Post.ContentHTML}}** `template.HTML` the content html of this post. **use this one to render post html**.
- **{{.Post.PreviewHTML}}** `template.HTML` the preview content html of this post in posts list page. **use this one to render post preview html**.

The tags can be ranged (sample code in `post.html`):

```html
{{range .Post.Tags}}
    <a href="{{.Url}}">{{.Name}}</a>
{{end}}
```

In post list template `posts.html`:

- **{{.Posts}}**    `[]*model.Post` the listed posts.
- **{{.Pager}}**    `*helper.Pager` the pager object, contains current, prev and next page values.

usage example:

```html
{{range .Posts}}
    <a href="{{.Url}}">{{.Title}}</a>
{{end}}
{{if .Pager.Prev}}<a href="{{.Pager.PrevUrl}}">prev</a>{{end}}
{{if .Pager.Next}}<a href="{{.Pager.NextUrl}}">next</a>{{end}}
```

**The variables are some to `index.html`**.

### Page Vars

In page template `page.html`. Major values are same to post, but `Tags`,`PreviewHTML` are dropped, `Meta` is added:

- **{{.Page}}** `*model.Page`   the page content.
- **{{.Page.Meta}}**    `map[string]string` the addition key values to the page.

### Archive Vars

In template `archive.html`:

- **{{.Archives}}** `*model.Archive`    the archive list divided by year.
