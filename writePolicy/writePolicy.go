package writePolicy

import (
	"context"
	"github.com/rag594/konfigStore/model"
)

type IWritePolicy[T model.TenantId, V any] interface {
	SetConfig(ctx context.Context, cacheKey string, entityId T, value *V) error
}
