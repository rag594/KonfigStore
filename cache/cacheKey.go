package cache

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/rag594/konfigStore/config"
	"strings"
)

const (
	defaultKeyPrefix = "KONFIG_STORE"
)

type CacheKey[T config.TenantId, V any] struct {
	EntityId  T
	Prefix    string
	ConfigKey string
}

type CacheKeyOptions[T config.TenantId, V any] func(*CacheKey[T, V])

func WithCacheKeyPrefix[T config.TenantId, V any](prefix string) CacheKeyOptions[T, V] {
	return func(c *CacheKey[T, V]) {
		c.Prefix = prefix
	}
}

func WithCacheConfigKey[T config.TenantId, V any](configKey string) CacheKeyOptions[T, V] {
	return func(c *CacheKey[T, V]) {
		c.ConfigKey = configKey
	}
}

func NewCacheKey[T config.TenantId, V any](entityId T, opts ...CacheKeyOptions[T, V]) *CacheKey[T, V] {
	c := &CacheKey[T, V]{
		EntityId:  entityId,
		Prefix:    defaultKeyPrefix,
		ConfigKey: strcase.ToScreamingSnake(strings.Split(fmt.Sprintf("%T", *new(V)), ".")[1]),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c

}

func (c *CacheKey[T, V]) DefaultValue() string {
	return fmt.Sprintf("%s_%v_%s", c.Prefix, c.EntityId, c.ConfigKey)
}
