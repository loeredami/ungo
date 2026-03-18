package ungo

import "os"

type Service interface {
	Init() error
	Name() string
	Shutdown() error
}

type ServiceRegistry struct {
	services *LinkedList[Service]
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: NewLinkedList[Service](),
	}
}

func (sr *ServiceRegistry) Add(service Service) {
	if err := service.Init(); err != nil {
		os.Exit(ExitConfigurationError)
	}
	sr.services.Add(service)

	OnShutdown(func() {
		service.Shutdown()
	})
}
