package registry

import (
	"github.com/go-kratos/kratos/v2/registry"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

type RegistrarCreator func(c *conf.Registry) registry.Registrar
type DiscoveryCreator func(c *conf.Registry) registry.Discovery

type RegistrarCreatorMap map[string]RegistrarCreator
type DiscoveryCreatorMap map[string]DiscoveryCreator

var registrars RegistrarCreatorMap
var discoveries DiscoveryCreatorMap

func RegisterRegistrarCreator(name string, reg RegistrarCreator) {
	if registrars == nil {
		registrars = make(RegistrarCreatorMap)
	}
	registrars[name] = reg
}

func GetRegistrarCreator(name string) RegistrarCreator {
	if registrars == nil {
		return nil
	}
	return registrars[name]
}

func RegisterDiscoveryCreator(name string, dis DiscoveryCreator) {
	if discoveries == nil {
		discoveries = make(DiscoveryCreatorMap)
	}
	discoveries[name] = dis
}

func GetDiscoveryCreator(name string) DiscoveryCreator {
	if discoveries == nil {
		return nil
	}
	return discoveries[name]
}
