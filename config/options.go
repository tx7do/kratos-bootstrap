package config

const remoteConfigSourceConfigFile = "remote.yaml"

type Type string

const (
	LocalFile  Type = "file"
	Nacos      Type = "nacos"
	Consul     Type = "consul"
	Etcd       Type = "etcd"
	Apollo     Type = "apollo"
	Kubernetes Type = "kubernetes"
	Polaris    Type = "polaris"
)
