```toml
title = "新建站点"
date = "2016-02-05 15:00:00"
slug = "zh/guide/create-new-site"
hover = "guide"
lang = "zh"
template = "guide.html"
```

`PuGo` 已经将静态资源编码到代码中，你可以从执行文件直接创建站点，无需另外下载资源：

```bash
pugo new site
```

程序解压内容到当前目录，如下：

    - meta.toml // 站点的基本信息
    - source // 站点的所有内容
    --|-- post // 所有文章
    --|-- page // 所有页面
    --|-- lang // 语言文件，如果需要国际化支持
    --|-- media // 资源文件，比如图片，附件等
    - theme // 所有主题文件
    --|-- default // 主题 'default' 的相关文件
    ------|-- *.html // 主题文件都是html文件
    ------|-- embed // 嵌入用的局部模板文件
    ------|-- static // 主题需要的静态文件，如样式表，脚本，图片
    

@image2

## meta.toml

站点的基本信息都保存在 `meta.toml`. 它定义有站点相关的各种信息，你需要先填写完整。

### Meta

```toml
[meta]
# 站点名称
title = "PuGo"

# 站点副标题
# 用于生成页面的标题
subtitle = "Static Site Generator"

# 站点关键字
# 用在 <meta content="keyword">
keyword = "pugo,golang,static,site,generator"

# 站点简介
# 用在 <meta content="description">
desc = "pugo is a simple static site generator"

# domain 和 root 用于生成正确的URL地址
# root 可以设置为子目录地址，比如 http://domain/blog
domain = "pugo.io"
root = "http://pugo.io/"
```

### 导航

导航是一组按顺序定义在 `[[nav]]` 区块下的项目，完整的字段有：

```toml
[[nav]]
# 导航跳转的地址
link = "/"

# 导航的标题
# 显示在超链接中的文本
title = "Home"

# 国际化支持的字段
i18n = "home"

# 标记导航选定状态的全局变量
hover = "home"

# 如果true，让浏览器打开新的选项卡跳转页面
blank = true
```

### 作者们

您可以添加一些作者给站点。文章和页面可以选中其中一人作为作者。

```toml
[[author]]
# 作者的名称，需要唯一
name = "pugo"

# 作者的联系邮箱，请注意保密
email = ""

# 作者的联系网址
url = "http://pugo.io"

# 作者的头像，默认用Gravatar的肖像
avatar = ""

# 作者的个人介绍
bio = "the robot of pugo, who generates all default contents."
```

第一个作者被认为是站点的 **拥有者**。
