package apollo

import (
	"strings"

	"github.com/apolloconfig/agollo/v4"
	apolloconfig "github.com/apolloconfig/agollo/v4/env/config"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
)

type apollo struct {
	client agollo.Client
	opt    *options
}

func NewSource(opts ...Option) config.Source {
	op := options{}
	for _, o := range opts {
		o(&op)
	}
	client, err := agollo.StartWithConfig(func() (*apolloconfig.AppConfig, error) {
		return &apolloconfig.AppConfig{
			AppID:            op.appid,
			Cluster:          op.cluster,
			NamespaceName:    op.namespace,
			IP:               op.endpoint,
			IsBackupConfig:   op.isBackupConfig,
			Secret:           op.secret,
			BackupConfigPath: op.backupPath,
		}, nil
	})
	if err != nil {
		panic(err)
	}
	return &apollo{client: client, opt: &op}
}

func (e *apollo) load() []*config.KeyValue {
	kvs := make([]*config.KeyValue, 0)
	namespaces := strings.Split(e.opt.namespace, ",")

	for _, ns := range namespaces {
		if !e.opt.originConfig {
			kv, err := e.getConfig(ns)
			if err != nil {
				log.Errorf("apollo get config failed，err:%v", err)
				continue
			}
			kvs = append(kvs, kv)
			continue
		}
		if strings.Contains(ns, ".") && !strings.HasSuffix(ns, "."+properties) &&
			(format(ns) == yaml || format(ns) == yml || format(ns) == json) {
			kv, err := e.getOriginConfig(ns)
			if err != nil {
				log.Errorf("apollo get config failed，err:%v", err)
				continue
			}
			kvs = append(kvs, kv)
			continue
		}
		kv, err := e.getConfig(ns)
		if err != nil {
			log.Errorf("apollo get config failed，err:%v", err)
			continue
		}
		kvs = append(kvs, kv)
	}
	return kvs
}

func (e *apollo) getConfig(ns string) (*config.KeyValue, error) {
	next := map[string]any{}
	e.client.GetConfigCache(ns).Range(func(key, value any) bool {
		// all values are out properties format
		resolve(genKey(ns, key.(string)), value, next)
		return true
	})
	f := format(ns)
	codec := encoding.GetCodec(f)
	val, err := codec.Marshal(next)
	if err != nil {
		return nil, err
	}
	return &config.KeyValue{
		Key:    ns,
		Value:  val,
		Format: f,
	}, nil
}

func (e *apollo) getOriginConfig(ns string) (*config.KeyValue, error) {
	value, err := e.client.GetConfigCache(ns).Get("content")
	if err != nil {
		return nil, err
	}
	// serialize the namespace content KeyValue into bytes.
	return &config.KeyValue{
		Key:    ns,
		Value:  []byte(value.(string)),
		Format: format(ns),
	}, nil
}

func (e *apollo) Load() (kv []*config.KeyValue, err error) {
	return e.load(), nil
}

func (e *apollo) Watch() (config.Watcher, error) {
	w, err := newWatcher(e)
	if err != nil {
		return nil, err
	}
	return w, nil
}
