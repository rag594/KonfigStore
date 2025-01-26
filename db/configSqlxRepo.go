package db

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/rag594/konfigStore/model"
)

type ConfigRepo[T model.TenantId, V any] struct {
	ConfigKey string
	Conn      *sqlx.DB
}

func RegisterConfigForDbOps[T model.TenantId, V any](conn *sqlx.DB, configKey string) *ConfigRepo[T, V] {
	return &ConfigRepo[T, V]{
		Conn:      conn,
		ConfigKey: configKey,
	}
}

func (c *ConfigRepo[T, V]) SaveConfig(ctx context.Context, config *model.Config[T, V]) (int64, error) {
	insertQuery := "INSERT INTO Config (entityId, configKey, value) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE value = ?"
	res, err := c.Conn.ExecContext(ctx, insertQuery, config.EntityId, c.ConfigKey, config.Value, config.Value)
	if err != nil {
		return -1, err
	}
	id, _ := res.LastInsertId()

	return id, nil

}

func (c *ConfigRepo[T, V]) GetConfigByKeyForEntity(ctx context.Context, entityId T) (*V, error) {
	d := &model.ConfigValue[V]{}
	err := c.Conn.GetContext(ctx, d, "select value from Config where configKey = ? and entityId = ?", c.ConfigKey, entityId)
	if err != nil {
		return nil, err
	}
	return d.Val, nil
}
