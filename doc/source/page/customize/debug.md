```ini
title = Debug Mode
slug = docs/debug
date = 2015-11-11
update_date = 2015-11-14
author = fuxiaohei
author_email = fuxiaohei@vip.qq.com
hover = docs
template =

[meta]
Source = "https://github.com/go-xiaohei/pugo-io/blob/master/source/page/customize/debug.md"
Version = ">=0.7.0"
```

There is a special flag when run `pugo`:

    $./pugo --debug

It watches `source` directory as default. But when `--debug`, it watches `template` too. So when template changes, it rebuilds either. That's useful when you creating new theme.

Meanwhile, in debug mode, the log lines are more, to record all running information.
