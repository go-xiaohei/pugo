package deploy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/codegangsta/cli"
	"gopkg.in/inconshreveable/log15.v2"
	"time"
)

type AwsS3 struct {
	Local     string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
}

func (a *AwsS3) Command() cli.Command {
	return cli.Command{
		Name:  "aws-s3",
		Usage: "deploy via aws-sdk to Amazon Storage Service",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "local", Value: "dest", Usage: "local website directory"},
			cli.StringFlag{Name: "ak", Usage: "accesss key"},
			cli.StringFlag{Name: "sk", Usage: "secret key"},
			cli.StringFlag{Name: "bucket", Usage: "storage bucket name"},
			cli.StringFlag{Name: "region", Usage: "storage bucket region"},
		},
		Action: func(ctx *cli.Context) {
			t := time.Now()
			a2, err := a.Create(ctx)
			if err != nil {
				log15.Error("AWS|Fail|%s", err.Error())
				return
			}
			if err = a2.Do(); err != nil {
				log15.Error("AWS|Fail|%s", err.Error())
				return
			}
			log15.Info("AWS|Finish|%s", time.Since(t))
		},
	}
}

func (a *AwsS3) String() string {
	return "AWS-S3"
}

func (a *AwsS3) Create(ctx *cli.Context) (Method, error) {
	a2 := &AwsS3{
		Local:     ctx.String("local"),
		AccessKey: ctx.String("ak"),
		SecretKey: ctx.String("sk"),
		Bucket:    ctx.String("bucket"),
		Region:    ctx.String("region"),
	}
	if !com.IsDir(a2.Local) {
		return nil, fmt.Errorf("directory '%s' is not existed", a2.Local)
	}
	if a2.AccessKey == "" || a2.SecretKey == "" {
		return nil, fmt.Errorf("S3's accessKey or secretKey is empty")
	}
	if a2.Bucket == "" || a2.Region == "" {
		return nil, fmt.Errorf("S3's Bucket or Region is not setted")
	}
	return a2, nil
}

func (a *AwsS3) Do() error {
	creds := credentials.NewStaticCredentials(a.AccessKey, a.SecretKey, "")
	_, err := creds.Get()
	if err != nil {
		return err
	}

	cfg := aws.NewConfig().WithRegion(a.Region).WithCredentials(creds)
	s3client := s3.New(session.New(), cfg)

	log15.Info("AWS|Bucket|%s", a.Region)

	err = filepath.Walk(a.Local, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(a.Local, path)
		rel = filepath.ToSlash(rel)

		fileData, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		size := info.Size()
		fileType := http.DetectContentType(fileData)
		params := &s3.PutObjectInput{
			Bucket:        aws.String(a.Bucket), // required
			Key:           aws.String(rel),      // required
			ACL:           aws.String("public-read"),
			Body:          bytes.NewReader(fileData),
			ContentLength: aws.Int64(size),
			ContentType:   aws.String(fileType),
			Metadata: map[string]*string{
				"Key": aws.String(filepath.Base(rel)), //required
			},
		}
		if _, err = s3client.PutObject(params); err != nil {
			return err
		}
		log15.Info("AWS|Upload|%s", rel)
		return nil
	})
	return err
}
