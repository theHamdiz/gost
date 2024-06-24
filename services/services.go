package services

type ServiceContainer struct {
	services map[string]interface{}
}

func NewServiceContainer() *ServiceContainer {
	return &ServiceContainer{
		services: make(map[string]interface{}),
	}
}

func (sc *ServiceContainer) RegisterService(name string, service interface{}) {
	sc.services[name] = service
}

func (sc *ServiceContainer) GetService(name string) interface{} {
	return sc.services[name]
}
