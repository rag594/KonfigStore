package config

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/iancoleman/strcase"
	"reflect"
	"strings"
)

type TenantId interface {
	int | string
}

type ConfigValue[V any] struct {
	Val *V `json:"value" db:"value"`
}

type Config[T TenantId, V any] struct {
	EntityId T              `json:"entity_id" db:"entityId"`
	Key      string         `json:"key" db:"configKey"`
	Value    ConfigValue[V] `json:"value" db:"-"`
}

func (c *Config[T, V]) GetKey() string {
	return c.Key
}

func GetGenericType(x any) string {
	// Get type as a string
	s := reflect.TypeOf(x).String()

	return strcase.ToScreamingSnake(strings.Split(s, ".")[1])
}

func NewConfig[T TenantId, V any](value *V) *Config[T, V] {
	return &Config[T, V]{
		Value: ConfigValue[V]{Val: value},
		Key:   GetGenericType(value),
	}
}

func (c *ConfigValue[V]) Scan(src interface{}) error {
	var source []byte
	_m1 := new(V)

	switch src.(type) {
	case []uint8:
		source = src.([]uint8)
	case nil:
		return nil
	default:
		return errors.New("incompatible type for ConfigValue")
	}
	err := json.Unmarshal(source, &_m1)
	if err != nil {
		return err
	}
	_m2 := &ConfigValue[V]{Val: _m1}
	*c = *_m2
	return nil
}

func (c ConfigValue[V]) Value() (driver.Value, error) {
	j, err := json.Marshal(c.Val)
	if err != nil {
		return nil, err
	}
	return driver.Value(j), nil
}
