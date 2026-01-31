package repository

import (
	"context"
	"sync"

	"snake-game/gateway/domain/entity"
)

type serviceRegistryImpl struct {
	services map[string]*entity.ServiceInfo
	mutex    sync.RWMutex
}

func NewServiceRegistry() *serviceRegistryImpl {
	return &serviceRegistryImpl{
		services: make(map[string]*entity.ServiceInfo),
	}
}

func (r *serviceRegistryImpl) RegisterService(ctx context.Context, service *entity.ServiceInfo) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.services[service.Name] = service
	return nil
}

func (r *serviceRegistryImpl) GetService(ctx context.Context, name string) (*entity.ServiceInfo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	service, exists := r.services[name]
	if !exists {
		return nil, nil
	}
	
	return service, nil
}

func (r *serviceRegistryImpl) GetAllServices(ctx context.Context) ([]*entity.ServiceInfo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	services := make([]*entity.ServiceInfo, 0, len(r.services))
	for _, service := range r.services {
		services = append(services, service)
	}
	
	return services, nil
}

func (r *serviceRegistryImpl) UpdateServiceHealth(ctx context.Context, name string, health bool) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	service, exists := r.services[name]
	if !exists {
		return nil
	}
	
	service.Health = health
	r.services[name] = service
	
	return nil
}