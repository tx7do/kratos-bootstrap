package config

type Type string

const (
	// TypeLocalFile is the local file config type.
	TypeLocalFile Type = "file"

	// TypeApollo is the apollo remote config type.
	TypeApollo Type = "apollo"

	// TypeConsul is the consul remote config type.
	TypeConsul Type = "consul"

	// TypeEtcd is the etcd remote config type.
	TypeEtcd Type = "etcd"

	// TypeKubernetes is the kubernetes remote config type.
	TypeKubernetes Type = "kubernetes"

	// TypeNacos is the nacos remote config type.
	TypeNacos Type = "nacos"

	// TypePolaris is the polaris remote config type.
	TypePolaris Type = "polaris"
)
