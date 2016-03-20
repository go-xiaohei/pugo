```toml
title = "FTP 和 SFTP"
date = "2016-02-04 15:00:00"
slug = "zh/docs/deploy/ftp-sftp"
hover = "docs"
lang = "zh"
template = "docs.html"
```

`PuGo` 可以使用 FTP and SFTP 账号发布，目前只支持 **用户名** 和 **密码** 登陆的方式。

```bash
pugo deploy ftp --local="dest" --user="user" --password="xxx" --host="127.0.0.1:21" --directory="pugo"
pugo deploy sftp --local="dest" --user="user" --password="xxx" --host="127.0.0.1:22" --directory="pugo"
```

`--local` 设置本地编译好的内容的文件夹。

`--user` 和 `--password` 设置连接的账号和密码。 SFTP 还不支持使用 `.ssh/keys` 登陆。

`--host` 设置连接的地址和端口。

`--directory` 设置线上保存的目录。SFTP 中 `~/` 代表用户目录。