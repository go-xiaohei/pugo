```toml
title = "AWS S3 and Qiniu"
date = "2016-02-04 15:00:00"
slug = "en/docs/deploy/cloud"
hover = "docs"
lang = "en"
template = "docs.html"
```

`PuGo` can upload static files to cloud storage service. Now it supports `AWS S3` and `Qiniu` cloud storage.

```bash
pugo deploy qiniu --local="dest" --ak="ak" --sk="sk" --bucket="bucket"
pugo deploy aws-s3 --local="dest" --ak="ak" --sk="sk" --bucket="bucket" --region="region"
```

`--local` set local directory that saving static files after building.

`--ak` and `--sk` set access-key and secret-key from cloud service.

`--bucket` set uploading bucket name.

In AWS S3, `--region` set the region of bucket, such as `us-east-1`.

When use cloud storage, please make sure your bucket is public. Then you can set domain alias to the bucket as static website host.

#### Warning

When deploying to storage service, `PuGo` uploads compiled files to remote bucket. If there are some files in same name, overwrite them. But if not, PuGo won't change them. So after uploaded, you may need to clean some old or useless files by yourself.