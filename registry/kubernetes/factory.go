package kubernetes

import (
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	k8s "k8s.io/client-go/kubernetes"
	k8sRest "k8s.io/client-go/rest"
	k8sTools "k8s.io/client-go/tools/clientcmd"
	k8sUtil "k8s.io/client-go/util/homedir"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	baseRegistry "github.com/tx7do/kratos-bootstrap/registry"
)

func init() {
	_ = baseRegistry.RegisterDiscoveryFactory(baseRegistry.Kubernetes, NewDiscovery)
	_ = baseRegistry.RegisterRegistrarFactory(baseRegistry.Kubernetes, NewRegistrar)
}

const serviceAccountNamespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

// NewRegistry 创建一个注册发现客户端 - Kubernetes
func NewRegistry(cfg *conf.Registry) (*Registry, error) {
	if cfg == nil || cfg.Kubernetes == nil {
		return nil, nil
	}

	in := cfg.Kubernetes

	var restConfig *k8sRest.Config
	var err error

	forceInCluster := in.GetInCluster() && in.InCluster != nil

	switch {
	case forceInCluster:
		restConfig, err = k8sRest.InClusterConfig()
	case in.GetKubeconfig() != "":
		restConfig, err = k8sTools.BuildConfigFromFlags("", in.GetKubeconfig())
	default:
		if restConfig, err = k8sRest.InClusterConfig(); err != nil {
			home := k8sUtil.HomeDir()
			kubeConfig := filepath.Join(home, ".kube", "config")
			restConfig, err = k8sTools.BuildConfigFromFlags("", kubeConfig)
		}
	}
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}

	clientSet, err := k8s.NewForConfig(restConfig)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}

	var namespace string
	if in.GetNamespace() != "" {
		namespace = in.GetNamespace()
	} else if data, rerr := os.ReadFile(serviceAccountNamespaceFile); rerr == nil && len(data) > 0 {
		namespace = string(data)
	}

	reg := New(clientSet, namespace)

	return reg, nil
}

func NewDiscovery(c *conf.Registry) (registry.Discovery, error) {
	return NewRegistry(c)
}

func NewRegistrar(c *conf.Registry) (registry.Registrar, error) {
	return NewRegistry(c)
}
