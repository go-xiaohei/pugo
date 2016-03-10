package deploy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"golang.org/x/net/context"
	"gopkg.in/inconshreveable/log15.v2"
	"qiniupkg.com/api.v7/kodo"
)

type Qiniu struct {
	Local     string
	AccessKey string
	SecretKey string
	Bucket    string
	// MaxAge int64
}

func (q *Qiniu) Command() cli.Command {
	return cli.Command{
		Name:  "qiniu",
		Usage: "deploy via qiniu-sdk to qiniu cloud storage",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "local", Value: "public", Usage: "local website directory"},
			cli.StringFlag{Name: "ak", Usage: "accesss key"},
			cli.StringFlag{Name: "sk", Usage: "secret key"},
			cli.StringFlag{Name: "bucket", Usage: "storage bucket name"},
		},
		Action: func(ctx *cli.Context) {
			q2, err := q.Create(ctx)
			if err != nil {
				log15.Error("Qiniu|Fail|%s", err.Error())
				return
			}
			if err = q2.Do(); err != nil {
				log15.Error("Qiniu|Fail|%s", err.Error())
				return
			}
			log15.Info("Qiniu|Finish")
		},
	}
}

func (q *Qiniu) String() string {
	return "Qiniu"
}

func (q *Qiniu) Create(ctx *cli.Context) (Method, error) {
	q2 := &Qiniu{
		Local:     ctx.String("local"),
		AccessKey: ctx.String("ak"),
		SecretKey: ctx.String("sk"),
		Bucket:    ctx.String("bucket"),
	}
	if !com.IsDir(q2.Local) {
		return nil, fmt.Errorf("directory '%s' is not existed", q2.Local)
	}
	if q2.AccessKey == "" || q2.SecretKey == "" {
		return nil, fmt.Errorf("Qiniu's accessKey or secretKey is empty")
	}
	if q2.Bucket == "" {
		return nil, fmt.Errorf("Qiniu's Bucket is not setted")
	}
	return q2, nil
}

func (q *Qiniu) Do() error {
	kodo.SetMac(q.AccessKey, q.SecretKey)
	client := kodo.New(0, nil)
	bucket := client.Bucket(q.Bucket)
	log15.Info("Qiniu|Bucket|%s", q.Bucket)
	err := filepath.Walk(q.Local, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(q.Local, path)
		rel = filepath.ToSlash(rel)

		ctx := context.Background()

		var ret interface{}
		if err = bucket.PutFile(ctx, ret, rel, path, nil); err != nil {
			return err
		}
		log15.Debug("Qiniu|Upload|%s", rel)
		return nil
	})
	return err
}
