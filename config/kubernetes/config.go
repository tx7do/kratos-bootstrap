package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/go-kratos/kratos/v2/config"
)

type kube struct {
	opts   options
	client *kubernetes.Clientset
}

// NewSource new a kubernetes config source.
func NewSource(opts ...Option) config.Source {
	op := options{}
	for _, o := range opts {
		o(&op)
	}
	return &kube{
		opts: op,
	}
}

func (k *kube) init() (err error) {
	var c *rest.Config
	if k.opts.KubeConfig != "" {
		if c, err = clientcmd.BuildConfigFromFlags(k.opts.Master, k.opts.KubeConfig); err != nil {
			return err
		}
	} else {
		if c, err = rest.InClusterConfig(); err != nil {
			return err
		}
	}
	if k.client, err = kubernetes.NewForConfig(c); err != nil {
		return err
	}
	return nil
}

func (k *kube) load() (kvs []*config.KeyValue, err error) {
	cmList, err := k.client.
		CoreV1().
		ConfigMaps(k.opts.Namespace).
		List(context.Background(), metav1.ListOptions{
			LabelSelector: k.opts.LabelSelector,
			FieldSelector: k.opts.FieldSelector,
		})
	if err != nil {
		return nil, err
	}
	for _, cm := range cmList.Items {
		kvs = append(kvs, k.configMap(cm)...)
	}
	return kvs, nil
}

func (k *kube) configMap(cm v1.ConfigMap) (kvs []*config.KeyValue) {
	for name, val := range cm.Data {
		k := fmt.Sprintf("%s/%s/%s", k.opts.Namespace, cm.Name, name)

		kvs = append(kvs, &config.KeyValue{
			Key:    k,
			Value:  []byte(val),
			Format: strings.TrimPrefix(filepath.Ext(k), "."),
		})
	}
	return kvs
}

func (k *kube) Load() ([]*config.KeyValue, error) {
	if k.opts.Namespace == "" {
		return nil, errors.New("options namespace not full")
	}
	if err := k.init(); err != nil {
		return nil, err
	}
	return k.load()
}

func (k *kube) Watch() (config.Watcher, error) {
	return newWatcher(k)
}
