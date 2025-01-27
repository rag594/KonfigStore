package readPolicy

import (
	"context"
	"github.com/rag594/konfigStore/config"
)

type IReadPolicy[T config.TenantId, V any] interface {
	GetConfig(ctx context.Context, cacheKey string, entityId T) (*V, error)
}
