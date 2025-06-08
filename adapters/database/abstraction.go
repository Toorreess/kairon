package db

import (
	"context"
	"errors"
	"kairon/cmd/api/infrastructure"
	"reflect"
)

type Connection struct {
	Client any
	Type   string
	Ctx    context.Context
}

type DBRepository interface {
	Get(index string, id string, hasDeleted bool) (map[string]any, error)
	Create(index string, cm any) (map[string]any, error)
	CreateWithID(index string, id string, cm any) (map[string]any, error)
	Update(index string, id string, hasDeleted bool, updates map[string]any) (map[string]any, error)
	Delete(index string, id string, hasDeleted bool) error
	List(index string, entity any, hasDeleted bool, queryOpts infrastructure.QueryOpts) ([]map[string]any, error)
	GetAll(index string, ids []string, entity any, hasDeleted bool) ([]map[string]any, []string, error)
	Close() error
	RunTransaction(ctx context.Context, f func(tx DBTransaction) error) error
}

type DBTransaction interface {
	Get(index string, id string, entity any) (map[string]any, error)
	GetAll(index string, ids []string, entity any) ([]map[string]any, []string, error)
	Create(index string, cm any) (map[string]any, error)
	CreateWithID(index, id string, cm any) (map[string]any, error)
	// Set(index string, id string, data any) error
	Update(index string, id string, entity any, updates map[string]any) error
	// Delete(index string, id string) error
	List(index string, entity any, queryOpts infrastructure.QueryOpts) ([]map[string]any, error)
}

func (conn *Connection) Close() error {
	return conn.Client.(DBRepository).Close()
}

// Ejemplo de modificación para delegar RunTransaction
func (conn *Connection) RunTransaction(ctx context.Context, f func(tx DBTransaction) error) error { // Ajustar firma si f devuelve valor
	if conn.Client == nil {
		return errors.New("Invalid client")
	}
	return conn.Client.(DBRepository).RunTransaction(ctx, f) // Llama al nuevo método en el cliente concreto
}

func (conn *Connection) Read(index string, id string, entity any) (map[string]any, error) {
	if conn.Client == nil {
		return nil, errors.New("Invalid client")
	}
	return conn.Client.(DBRepository).Get(index, id, entityHasDeleted(entity))
}

func (conn *Connection) Create(indexName string, cm any) (map[string]any, error) {
	if conn.Client == nil {
		return nil, errors.New("Invalid client")
	}
	return conn.Client.(DBRepository).Create(indexName, cm)
}

func (conn *Connection) CreateWithID(index string, id string, cm any) (map[string]any, error) {
	if conn.Client == nil {
		return nil, errors.New("Invalid client")
	}
	return conn.Client.(DBRepository).CreateWithID(index, id, cm)

}

func (conn *Connection) Update(index string, id string, entity any, changes map[string]any) (map[string]any, error) {
	if conn.Client == nil {
		return nil, errors.New("Invalid client")
	}
	return conn.Client.(DBRepository).Update(index, id, entityHasDeleted(entity), changes)
}

func (conn *Connection) Delete(index string, id string, entity any) error {
	if conn.Client == nil {
		return errors.New("Invalid client")
	}
	return conn.Client.(DBRepository).Delete(index, id, entityHasDeleted(entity))
}

func (conn *Connection) List(indexName string, entity any, queryOpts infrastructure.QueryOpts) ([]map[string]any, error) {
	if conn.Client == nil {
		return nil, errors.New("Invalid client")
	}
	return conn.Client.(DBRepository).List(indexName, entity, entityHasDeleted(entity), queryOpts)
}

func (conn *Connection) GetAll(indexName string, ids []string, entity any) ([]map[string]any, []string, error) {
	if conn.Client == nil {
		return nil, nil, errors.New("Invalid client")
	}
	return conn.Client.(DBRepository).GetAll(indexName, ids, entity, entityHasDeleted(entity))
}

func entityHasDeleted(entity any) bool {
	hasDeleted := false
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	e := val.Type()

	for i := 0; i < e.NumField(); i++ {
		fieldName := e.Field(i).Name
		if fieldName == "Deleted" {
			hasDeleted = true
		}
	}
	return hasDeleted
}
