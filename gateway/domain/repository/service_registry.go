package repository

import (
	"context"
	"snake-game/gateway/domain/entity"
)

type ServiceRegistry interface {
	RegisterService(ctx context.Context, service *entity.ServiceInfo) error
	GetService(ctx context.Context, name string) (*entity.ServiceInfo, error)
	GetAllServices(ctx context.Context) ([]*entity.ServiceInfo, error)
	UpdateServiceHealth(ctx context.Context, name string, health bool) error
}