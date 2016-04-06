```toml
title = "AWS S3 和 七牛云"
date = "2016-02-04 15:00:00"
slug = "zh/docs/deploy/cloud"
hover = "docs"
lang = "zh"
template = "docs.html"
```

`PuGo` 可以向 `AWS S3` 和 `七牛云存储` 上传内容发布。

```bash
pugo deploy qiniu --local="dest" --ak="ak" --sk="sk" --bucket="bucket"
pugo deploy aws-s3 --local="dest" --ak="ak" --sk="sk" --bucket="bucket" --region="region"
```

`--local` 设置本地编译好的内容的文件夹。

`--ak` 和 `--sk` 设置 access-key 和 secret-key。

`--bucket` 设置云服务的 bucket 名称。

AWS S3 还需要 `--region` 设置 bucket 所在 region，如 `us-east-1`.

使用云存储时，你需要确认 bucket 是公开的。然后你可以查询相关文档绑定域名，设置主页等。

#### 注意

`PuGo` 上传编译好的文件到云服务时，会将重名的文件覆盖，其他文件不会变化。因此可能你需要自己检查 bucket 中的无用文件，并自己删除。