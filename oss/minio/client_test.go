package minio

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mio "github.com/minio/minio-go/v7"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	cfg := &conf.OSS{
		Minio: &conf.OSS_MinIO{
			Endpoint:  "localhost:9000",
			AccessKey: "minioadmin",
			SecretKey: "minioadmin",
			UseSsl:    false,
		},
	}

	client := NewClient(cfg)
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
}

func TestNewStorage(t *testing.T) {
	t.Parallel()

	cfg := &conf.OSS{
		Minio: &conf.OSS_MinIO{
			Endpoint:  "localhost:9000",
			AccessKey: "minioadmin",
			SecretKey: "minioadmin",
			UseSsl:    false,
		},
	}

	storage := NewStorage(cfg)
	if storage == nil {
		t.Fatal("NewStorage() returned nil")
	}
	if storage.SDK() == nil {
		t.Fatal("NewStorage().SDK() returned nil")
	}
}

func TestStoragePutObjectValidation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name        string
		storage     *Storage
		bucket      string
		key         string
		body        func() *bytes.Buffer
		contentType string
		wantErr     error
	}{
		{name: "nil storage", storage: nil, bucket: "bucket", key: "test.txt", body: func() *bytes.Buffer { return bytes.NewBufferString("hello") }, contentType: "text/plain", wantErr: ErrNilClient},
		{name: "nil client", storage: &Storage{}, bucket: "bucket", key: "test.txt", body: func() *bytes.Buffer { return bytes.NewBufferString("hello") }, contentType: "text/plain", wantErr: ErrNilClient},
		{name: "empty bucket", storage: &Storage{client: &mio.Client{}}, bucket: "", key: "test.txt", body: func() *bytes.Buffer { return bytes.NewBufferString("hello") }, contentType: "text/plain", wantErr: ErrEmptyBucket},
		{name: "empty key", storage: &Storage{client: &mio.Client{}}, bucket: "bucket", key: "", body: func() *bytes.Buffer { return bytes.NewBufferString("hello") }, contentType: "text/plain", wantErr: ErrEmptyObjectKey},
		{name: "nil body", storage: &Storage{client: &mio.Client{}}, bucket: "bucket", key: "test.txt", body: nil, contentType: "text/plain", wantErr: ErrNilObjectBody},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var body *bytes.Buffer
			if tt.body != nil {
				body = tt.body()
			}

			_, err := tt.storage.PutObject(ctx, tt.bucket, tt.key, body, tt.contentType)
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
		bucket  string
		key     string
		wantErr error
	}{
		{name: "nil storage", storage: nil, bucket: "bucket", key: "test.txt", wantErr: ErrNilClient},
		{name: "nil client", storage: &Storage{}, bucket: "bucket", key: "test.txt", wantErr: ErrNilClient},
		{name: "empty bucket", storage: &Storage{client: &mio.Client{}}, bucket: "", key: "test.txt", wantErr: ErrEmptyBucket},
		{name: "empty key", storage: &Storage{client: &mio.Client{}}, bucket: "bucket", key: "", wantErr: ErrEmptyObjectKey},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.storage.GetObject(ctx, tt.bucket, tt.key)
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
		if r.Method == http.MethodGet && r.URL.Path == "/"+bucket+"/" && r.URL.Query().Has("location") {
			w.Header().Set("Content-Type", "application/xml")
			_, _ = w.Write([]byte("<LocationConstraint xmlns=\"http://s3.amazonaws.com/doc/2006-03-01/\">us-east-1</LocationConstraint>"))
			return
		}

		switch r.Method {
		case http.MethodPut:
			putCalled = true
			if r.URL.Path != "/"+bucket+"/"+objectKey {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if got := r.Header.Get("Content-Type"); got != contentType && got != "application/octet-stream" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			data, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !strings.Contains(string(data), payload) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Header().Set("ETag", "\"test-etag\"")
			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
			getCalled = true
			if r.URL.Path != "/"+bucket+"/"+objectKey {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", contentType)
			w.Header().Set("Last-Modified", "Mon, 02 Jan 2006 15:04:05 GMT")
			w.Header().Set("ETag", "\"test-etag\"")
			_, _ = w.Write([]byte(payload))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()

	endpoint := strings.TrimPrefix(server.URL, "http://")
	cfg := &conf.OSS{
		Minio: &conf.OSS_MinIO{
			Endpoint:  endpoint,
			AccessKey: "minioadmin",
			SecretKey: "minioadmin",
			UseSsl:    false,
		},
	}

	storage := NewStorage(cfg)
	if storage == nil {
		t.Fatal("NewStorage() returned nil")
	}

	putOutput, err := storage.PutObject(context.Background(), bucket, objectKey, strings.NewReader(payload), contentType)
	if err != nil {
		t.Fatalf("PutObject() error = %v", err)
	}
	if putOutput.ETag != "test-etag" {
		t.Fatalf("PutObject() unexpected etag: %s", putOutput.ETag)
	}

	obj, err := storage.GetObject(context.Background(), bucket, objectKey)
	if err != nil {
		t.Fatalf("GetObject() error = %v", err)
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
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

func TestPrepareBody(t *testing.T) {
	t.Parallel()

	t.Run("read seeker", func(t *testing.T) {
		t.Parallel()

		body, size, err := prepareBody(strings.NewReader("hello"))
		if err != nil {
			t.Fatalf("prepareBody() error = %v", err)
		}
		if size != 5 {
			t.Fatalf("prepareBody() size = %d, want %d", size, 5)
		}

		data, err := io.ReadAll(body)
		if err != nil {
			t.Fatalf("failed reading prepared body: %v", err)
		}
		if string(data) != "hello" {
			t.Fatalf("prepared body = %q, want %q", string(data), "hello")
		}
	})

	t.Run("plain reader", func(t *testing.T) {
		t.Parallel()

		body, size, err := prepareBody(bytes.NewBufferString("world"))
		if err != nil {
			t.Fatalf("prepareBody() error = %v", err)
		}
		if size != 5 {
			t.Fatalf("prepareBody() size = %d, want %d", size, 5)
		}

		data, err := io.ReadAll(body)
		if err != nil {
			t.Fatalf("failed reading prepared body: %v", err)
		}
		if string(data) != "world" {
			t.Fatalf("prepared body = %q, want %q", string(data), "world")
		}
	})
}
