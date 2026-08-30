package nacos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestNewNacosRegistry(t *testing.T) {
	cfg := conf.Registry{
		Nacos: &conf.Registry_Nacos{
			Address: "127.0.0.1",
			Port:    8848,
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)
}

func TestNewNacosRegistryWithFullParams(t *testing.T) {
	cfg := conf.Registry{
		Nacos: &conf.Registry_Nacos{
			Address:             "127.0.0.1",
			Port:                8848,
			NamespaceId:         "ns",
			RegionId:            "cn-hangzhou",
			AppName:             "app",
			AppKey:              "key",
			Username:            "nacos",
			Password:            "nacos",
			Timeout:             durationpb.New(5 * time.Second),
			BeatInterval:        durationpb.New(3 * time.Second),
			UpdateThreadNum:     10,
			NotLoadCacheAtStart: true,
			Scheme:              "https",
			GrpcPort:            9948,
			AppendToStdout:      true,
			AsyncUpdateService:  true,
			Group:               "GROUP1",
			Cluster:             "CLUSTER1",
			Weight:              50,
			Kind:                "http",
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)

	assert.Equal(t, "GROUP1", reg.opts.group)
	assert.Equal(t, "CLUSTER1", reg.opts.cluster)
	assert.Equal(t, float64(50), reg.opts.weight)
	assert.Equal(t, "http", reg.opts.kind)
}
