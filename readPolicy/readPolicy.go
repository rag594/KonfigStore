package readPolicy

import (
	"context"
	"github.com/rag594/konfigStore/model"
)

type IReadPolicy[T model.TenantId, V any] interface {
	GetConfig(ctx context.Context, cacheKey string, entityId T) (*V, error)
}
