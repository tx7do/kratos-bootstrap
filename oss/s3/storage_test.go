package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestNewStorage(t *testing.T) {
	t.Parallel()

	cfg := &conf.OSS{
		S3: &conf.OSS_S3{
			Endpoint:       "localhost:9000",
			Region:         "us-east-1",
			Bucket:         "test-bucket",
			AccessKey:      "access-key",
			SecretKey:      "secret-key",
			UseSsl:         false,
			ForcePathStyle: true,
		},
	}

	storage := NewStorage(cfg)
	if storage == nil {
		t.Fatal("NewStorage() returned nil")
	}
	if storage.SDK() == nil {
		t.Fatal("NewStorage().SDK() returned nil")
	}
	if storage.Bucket() != "test-bucket" {
		t.Fatalf("NewStorage().Bucket() = %q, want %q", storage.Bucket(), "test-bucket")
	}
}

func TestStoragePutObjectValidation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name        string
		storage     *Storage
		key         string
		body        func() *bytes.Buffer
		contentType string
		wantErr     error
	}{
		{name: "nil storage", storage: nil, key: "test.txt", body: func() *bytes.Buffer { return bytes.NewBufferString("hello") }, contentType: "text/plain", wantErr: ErrNilClient},
		{name: "nil client", storage: &Storage{bucket: "bucket"}, key: "test.txt", body: func() *bytes.Buffer { return bytes.NewBufferString("hello") }, contentType: "text/plain", wantErr: ErrNilClient},
		{name: "empty bucket", storage: &Storage{client: &awss3.Client{}}, key: "test.txt", body: func() *bytes.Buffer { return bytes.NewBufferString("hello") }, contentType: "text/plain", wantErr: ErrEmptyBucket},
		{name: "empty key", storage: &Storage{client: &awss3.Client{}, bucket: "bucket"}, key: "", body: func() *bytes.Buffer { return bytes.NewBufferString("hello") }, contentType: "text/plain", wantErr: ErrEmptyObjectKey},
		{name: "nil body", storage: &Storage{client: &awss3.Client{}, bucket: "bucket"}, key: "test.txt", body: nil, contentType: "text/plain", wantErr: ErrNilObjectBody},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var body *bytes.Buffer
			if tt.body != nil {
				body = tt.body()
			}

			_, err := tt.storage.PutObject(ctx, tt.key, body, tt.contentType)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("PutObject() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorageGetObjectValidation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		storage *Storage
		key     string
		wantErr error
	}{
		{name: "nil storage", storage: nil, key: "test.txt", wantErr: ErrNilClient},
		{name: "nil client", storage: &Storage{bucket: "bucket"}, key: "test.txt", wantErr: ErrNilClient},
		{name: "empty bucket", storage: &Storage{client: &awss3.Client{}}, key: "test.txt", wantErr: ErrEmptyBucket},
		{name: "empty key", storage: &Storage{client: &awss3.Client{}, bucket: "bucket"}, key: "", wantErr: ErrEmptyObjectKey},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.storage.GetObject(ctx, tt.key)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("GetObject() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestStoragePutAndGetObject(t *testing.T) {
	t.Parallel()

	const (
		bucket      = "test-bucket"
		objectKey   = "demo/hello.txt"
		contentType = "text/plain"
		payload     = "hello world"
	)

	var putCalled bool
	var getCalled bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			putCalled = true
			if r.URL.Path != "/"+bucket+"/"+objectKey {
				t.Fatalf("unexpected PUT path: %s", r.URL.Path)
			}
			if got := r.Header.Get("Content-Type"); got != contentType {
				t.Fatalf("unexpected Content-Type: %s", got)
			}
			data, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("failed reading PUT body: %v", err)
			}
			if string(data) != payload {
				t.Fatalf("unexpected PUT body: %s", string(data))
			}
			w.Header().Set("ETag", "\"test-etag\"")
			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
			getCalled = true
			if r.URL.Path != "/"+bucket+"/"+objectKey {
				t.Fatalf("unexpected GET path: %s", r.URL.Path)
			}
			w.Header().Set("Content-Type", contentType)
			_, _ = w.Write([]byte(payload))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()

	endpoint := strings.TrimPrefix(server.URL, "http://")
	cfg := &conf.OSS{
		S3: &conf.OSS_S3{
			Endpoint:       endpoint,
			Region:         "us-east-1",
			Bucket:         bucket,
			AccessKey:      "access-key",
			SecretKey:      "secret-key",
			UseSsl:         false,
			ForcePathStyle: true,
		},
	}

	storage := NewStorage(cfg)
	if storage == nil {
		t.Fatal("NewStorage() returned nil")
	}

	putOutput, err := storage.PutObject(context.Background(), objectKey, strings.NewReader(payload), contentType)
	if err != nil {
		t.Fatalf("PutObject() error = %v", err)
	}
	if putOutput == nil || putOutput.ETag == nil || *putOutput.ETag != "\"test-etag\"" {
		t.Fatalf("PutObject() unexpected output: %#v", putOutput)
	}

	getOutput, err := storage.GetObject(context.Background(), objectKey)
	if err != nil {
		t.Fatalf("GetObject() error = %v", err)
	}
	defer getOutput.Body.Close()

	data, err := io.ReadAll(getOutput.Body)
	if err != nil {
		t.Fatalf("failed reading GetObject() body: %v", err)
	}
	if string(data) != payload {
		t.Fatalf("GetObject() body = %q, want %q", string(data), payload)
	}
	if !putCalled || !getCalled {
		t.Fatalf("expected both PUT and GET requests to be called, put=%v get=%v", putCalled, getCalled)
	}
}
