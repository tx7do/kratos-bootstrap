package registry

type Type string

const (
	Consul      Type = "consul"
	Etcd        Type = "etcd"
	ZooKeeper   Type = "zookeeper"
	Nacos       Type = "nacos"
	Kubernetes  Type = "kubernetes"
	Eureka      Type = "eureka"
	Polaris     Type = "polaris"
	Servicecomb Type = "servicecomb"
)
