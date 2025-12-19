package registry

type Type string

const (
	Consul      Type = "consul"
	Etcd        Type = "etcd"
	Eureka      Type = "eureka"
	Kubernetes  Type = "kubernetes"
	Nacos       Type = "nacos"
	Polaris     Type = "polaris"
	Servicecomb Type = "servicecomb"
	ZooKeeper   Type = "zookeeper"
)
