package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func NewClient(cfg *conf.OSS) *awss3.Client {
	if cfg == nil || cfg.GetS3() == nil {
		log.Fatal("missing s3 configuration")
		return nil
	}

	s3Cfg := cfg.GetS3()

	region := s3Cfg.GetRegion()
	if region == "" {
		region = "us-east-1"
	}

	loadOpts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(region),
	}

	if s3Cfg.GetAccessKey() != "" || s3Cfg.GetSecretKey() != "" || s3Cfg.GetToken() != "" {
		loadOpts = append(loadOpts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(s3Cfg.GetAccessKey(), s3Cfg.GetSecretKey(), s3Cfg.GetToken()),
		))
	}

	endpoint := normalizeEndpoint(s3Cfg.GetEndpoint(), s3Cfg.GetUseSsl())

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), loadOpts...)
	if err != nil {
		log.Fatal("failed loading aws s3 config", err)
		return nil
	}

	return awss3.NewFromConfig(awsCfg, func(o *awss3.Options) {
		o.UsePathStyle = s3Cfg.GetForcePathStyle()
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})
}

func normalizeEndpoint(endpoint string, useSSL bool) string {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return ""
	}

	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return endpoint
	}

	scheme := "http"
	if useSSL {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s", scheme, endpoint)
}

func isNilReader(body io.Reader) bool {
	if body == nil {
		return true
	}

	v := reflect.ValueOf(body)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func prepareBody(body io.Reader) (io.Reader, int64, error) {
	if rs, ok := body.(io.ReadSeeker); ok {
		size, err := readerSize(rs)
		if err != nil {
			return nil, 0, err
		}
		return rs, size, nil
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, 0, err
	}

	return bytes.NewReader(data), int64(len(data)), nil
}

func readerSize(rs io.ReadSeeker) (int64, error) {
	current, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	end, err := rs.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	_, err = rs.Seek(current, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return end - current, nil
}
