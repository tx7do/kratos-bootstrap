package polaris

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/polarismesh/polaris-go"

	"github.com/go-kratos/kratos/v2/config"
)

type source struct {
	client  polaris.ConfigAPI
	options *options
}

func New(client polaris.ConfigAPI, opts ...Option) (config.Source, error) {
	o := &options{
		namespace: "default",
		fileGroup: "",
		fileName:  "",
	}

	for _, opt := range opts {
		opt(o)
	}

	if o.fileGroup == "" {
		return nil, errors.New("fileGroup invalid")
	}

	if o.fileName == "" {
		return nil, errors.New("fileName invalid")
	}

	return &source{
		client:  client,
		options: o,
	}, nil
}

// Load return the config values
func (s *source) Load() ([]*config.KeyValue, error) {
	configFile, err := s.client.GetConfigFile(s.options.namespace, s.options.fileGroup, s.options.fileName)
	if err != nil {
		fmt.Println("fail to get config.", err)
		return nil, err
	}

	content := configFile.GetContent()
	k := s.options.fileName

	s.options.configFile = configFile

	return []*config.KeyValue{
		{
			Key:    k,
			Value:  []byte(content),
			Format: strings.TrimPrefix(filepath.Ext(k), "."),
		},
	}, nil
}

// Watch return the watcher
func (s *source) Watch() (config.Watcher, error) {
	return newWatcher(s.options.configFile), nil
}
