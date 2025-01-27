package writePolicy

import (
	"context"
	"github.com/rag594/konfigStore/config"
)

type IWritePolicy[T config.TenantId, V any] interface {
	SetConfig(ctx context.Context, cacheKey string, entityId T, value *V) error
}
