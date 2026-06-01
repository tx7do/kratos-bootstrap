package s3

import (
	"bytes"
	"io"
	"strings"
	"testing"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestNormalizeEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		endpoint string
		useSSL   bool
		want     string
	}{
		{name: "empty", endpoint: "", useSSL: true, want: ""},
		{name: "trim and https", endpoint: " s3.amazonaws.com ", useSSL: true, want: "https://s3.amazonaws.com"},
		{name: "http", endpoint: "127.0.0.1:9000", useSSL: false, want: "http://127.0.0.1:9000"},
		{name: "keep existing https", endpoint: "https://s3.us-east-1.amazonaws.com", useSSL: false, want: "https://s3.us-east-1.amazonaws.com"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := normalizeEndpoint(tt.endpoint, tt.useSSL)
			if got != tt.want {
				t.Fatalf("normalizeEndpoint() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	cfg := &conf.OSS{
		S3: &conf.OSS_S3{
			Endpoint:       "localhost:9000",
			Region:         "ap-southeast-1",
			AccessKey:      "access-key",
			SecretKey:      "secret-key",
			Token:          "token",
			UseSsl:         false,
			ForcePathStyle: true,
			Bucket:         "test-bucket",
		},
	}

	client := NewClient(cfg)
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
}

func TestNewClient_DefaultRegion(t *testing.T) {
	t.Parallel()

	cfg := &conf.OSS{
		S3: &conf.OSS_S3{
			Endpoint:  "s3.amazonaws.com",
			AccessKey: "access-key",
			SecretKey: "secret-key",
			UseSsl:    true,
		},
	}

	client := NewClient(cfg)
	if client == nil {
		t.Fatal("NewClient() returned nil when region is empty")
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
