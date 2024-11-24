package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/knadh/koanf/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type UploadConfig struct {
	Endpoint     string
	Region       string
	AccessKey    string
	SecretKey    string
	MainBucket   string
	PublicBucket string
	Secure       bool
}

func NewUploadServiceFromApplicationConfig(cfg *koanf.Koanf) *UploadService {
	uploadConfig := UploadConfig{
		Endpoint:     cfg.String("minio.endpoint"),
		Region:       cfg.String("minio.region"),
		AccessKey:    cfg.String("minio.access-key"),
		SecretKey:    cfg.String("minio.secret-key"),
		MainBucket:   cfg.String("minio.bucket"),
		PublicBucket: cfg.String("minio.public-bucket"),
		Secure:       cfg.Bool("minio.secure"),
	}
	return NewUploadService(uploadConfig)
}

type UploadService struct {
	Client       *minio.Client
	MainBucket   string
	PublicBucket string
	Region       string
	Secure       bool
	Endpoint     string
}

func NewUploadService(cfg UploadConfig) *UploadService {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.Secure,
	})
	if err != nil {
		slog.Error("failed to create minion client instance", "endpoint", cfg.Endpoint, "err", err.Error())
		panic(err)
	}

	service := new(UploadService)
	service.Client = client
	service.MainBucket = cfg.MainBucket
	service.Region = cfg.Region
	service.PublicBucket = cfg.PublicBucket
	service.Secure = cfg.Secure
	service.Endpoint = cfg.Endpoint

	return service
}

func (s *UploadService) GetPublicURL(path string) string {
	var (
		schema string
	)

	if s.Secure {
		schema = "https"
	} else {
		schema = "http"
	}

	return fmt.Sprintf("%s://%s/%s/%s", schema, s.Endpoint, s.PublicBucket, path)
}

func (s *UploadService) InitBuckets(ctx context.Context) error {

	err := s.makeMainBucket(ctx)
	if err != nil {
		return err
	}

	err = s.makePublicBucket(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *UploadService) makeMainBucket(ctx context.Context) error {
	exists, err := s.Client.BucketExists(ctx, s.MainBucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	err = s.Client.MakeBucket(ctx, s.MainBucket, minio.MakeBucketOptions{
		Region:        s.Region,
		ObjectLocking: false,
	})
	s.Client.SetBucketPolicy(ctx, s.PublicBucket, "")
	return err
}

func (s *UploadService) makePublicBucket(ctx context.Context) error {
	exists, err := s.Client.BucketExists(ctx, s.PublicBucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	err = s.Client.MakeBucket(ctx, s.PublicBucket, minio.MakeBucketOptions{
		Region:        s.Region,
		ObjectLocking: false,
	})
	if err != nil {
		return err
	}
	err = s.Client.SetBucketPolicy(ctx, s.PublicBucket, publicBucketPolicy)
	return err
}

const (
	publicBucketPolicy = `{
   "Version":"2012-10-17",
   "Statement":[
      {
         "Effect":"Allow",
         "Principal":{
            "AWS":[
               "*"
            ]
         },
         "Action":[
            "s3:GetBucketLocation",
            "s3:ListBucket"
         ],
         "Resource":[
            "arn:aws:s3:::openlibrary-public"
         ]
      },
      {
         "Effect":"Allow",
         "Principal":{
            "AWS":[
               "*"
            ]
         },
         "Action":[
            "s3:GetObject"
         ],
         "Resource":[
            "arn:aws:s3:::openlibrary-public/*"
         ]
      }
   ]
}`
)
