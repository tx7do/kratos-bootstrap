package minio

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"reflect"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	tlsUtils "github.com/tx7do/go-utils/tls"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

var (
	ErrNilClient      = errors.New("minio client is nil")
	ErrEmptyBucket    = errors.New("minio bucket is empty")
	ErrEmptyObjectKey = errors.New("minio object key is empty")
	ErrNilObjectBody  = errors.New("minio object body is nil")
)

type Storage struct {
	client *minio.Client
}

func NewClient(cfg *conf.OSS) *minio.Client {
	if cfg == nil || cfg.GetMinio() == nil {
		log.Errorf("missing minio configuration")
		return nil
	}

	minioCfg := cfg.GetMinio()

	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(minioCfg.GetAccessKey(), minioCfg.GetSecretKey(), minioCfg.GetToken()),
		Secure: minioCfg.GetUseSsl(),
		Region: minioCfg.GetRegion(),
	}
	if minioCfg.GetForcePathStyle() {
		opts.BucketLookup = minio.BucketLookupPath
	}
	if tlsConf := minioCfg.GetTls(); tlsConf != nil {
		tlsCfg, err := loadClientTlsConfig(tlsConf)
		if err != nil {
			log.Errorf("failed load minio tls config: %v", err)
			return nil
		}
		if tlsCfg != nil {
			opts.Transport = &http.Transport{
				TLSClientConfig: tlsCfg,
			}
		}
	}

	impl, err := minio.New(minioCfg.GetEndpoint(), opts)
	if err != nil {
		log.Errorf("failed opening connection to minio: %v", err)
		return nil
	}

	return impl
}

func NewStorage(cfg *conf.OSS) *Storage {
	return &Storage{client: NewClient(cfg)}
}

func (s *Storage) SDK() *minio.Client {
	if s == nil {
		return nil
	}

	return s.client
}

func (s *Storage) PutObject(ctx context.Context, bucket, key string, body io.Reader, contentType string) (minio.UploadInfo, error) {
	if s == nil || s.client == nil {
		return minio.UploadInfo{}, ErrNilClient
	}
	if bucket == "" {
		return minio.UploadInfo{}, ErrEmptyBucket
	}
	if key == "" {
		return minio.UploadInfo{}, ErrEmptyObjectKey
	}
	if isNilReader(body) {
		return minio.UploadInfo{}, ErrNilObjectBody
	}

	preparedBody, size, err := prepareBody(body)
	if err != nil {
		return minio.UploadInfo{}, err
	}

	return s.client.PutObject(ctx, bucket, key, preparedBody, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
}

func (s *Storage) GetObject(ctx context.Context, bucket, key string) (*minio.Object, error) {
	if s == nil || s.client == nil {
		return nil, ErrNilClient
	}
	if bucket == "" {
		return nil, ErrEmptyBucket
	}
	if key == "" {
		return nil, ErrEmptyObjectKey
	}

	return s.client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
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

func loadClientTlsConfig(cfg *conf.TLS) (*tls.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	var tlsCfg *tls.Config
	var err error

	if cfg.File != nil {
		if tlsCfg, err = tlsUtils.LoadClientTlsConfigFile(
			cfg.File.GetKeyPath(),
			cfg.File.GetCertPath(),
			cfg.File.GetCaPath(),
		); err != nil {
			return nil, err
		}
	} else if cfg.Config != nil {
		if tlsCfg, err = tlsUtils.LoadClientTlsConfigString(
			cfg.Config.GetKeyPem(),
			cfg.Config.GetCertPem(),
			cfg.Config.GetCaPem(),
		); err != nil {
			return nil, err
		}
	}

	if tlsCfg != nil && cfg.GetInsecureSkipVerify() {
		tlsCfg.InsecureSkipVerify = true
	}

	return tlsCfg, nil
}
