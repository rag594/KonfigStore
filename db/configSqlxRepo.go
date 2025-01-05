package db

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/rag594/konfigStore/model"
)

type ConfigRepo[T model.TenantId, V any] struct {
	Conn *sqlx.DB
}

func RegisterConfigForDbOps[T model.TenantId, V any](conn *sqlx.DB) *ConfigRepo[T, V] {
	return &ConfigRepo[T, V]{
		Conn: conn,
	}
}

func (c *ConfigRepo[T, V]) SaveConfig(ctx context.Context, config *model.Config[T, V]) (int64, error) {
	insertQuery := "INSERT INTO Config (entityId, configKey, value) VALUES (?, ?, ?)"
	res, err := c.Conn.ExecContext(ctx, insertQuery, config.EntityId, config.Key, config.Value)
	if err != nil {
		return -1, err
	}
	id, _ := res.LastInsertId()

	return id, nil

}

func (c *ConfigRepo[T, V]) GetConfigByKeyForEntity(ctx context.Context, key string, entityId T) (*V, error) {
	d := &model.ConfigValue[V]{}
	err := c.Conn.GetContext(ctx, d, "select value from Config where configKey = ? and entityId = ?", key, entityId)
	if err != nil {
		return nil, err
	}
	return d.Val, nil
}
