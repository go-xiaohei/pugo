```ini
[meta]
title = Pugo.Static
subtitle = site generator
; print in html <meta>
keyword = pugo,golang,static,site,generator
; print in html <meta>
desc = pugo is a simple static site generator
; build links for feed, sitemap
domain = localhost
; root path for site, if empty, build as http://{domain}/
root = http://localhost/

; nav data
[nav]
; reference to [home] block, same below.
-:home
-:archive
-:about
-:source

[nav.home]
link = /
title = Home
i18n = homepage
; set nav to active status
hover = home

[nav.archive]
link = /archive
title = Archive
i18n = archive
hover = archive

[nav.about]
link = /about
title = About
i18n = about
; browser open in new tab
blank = true
hover = about

; author data
[author]
-:pugo
-:fuxiaohei

[author.pugo]
name = pugo
email =
url = http://pugo.io
avatar =

[author.fuxiaohei]
name = fuxiaohei
nick = 傅小黑
email = fuxioahei@vip.qq.com
url = http://fuxiaohei.me
avatar =

; comment settings
[comment.disqus]
site = fuxiaohei

```
