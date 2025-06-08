package repositories

import (
	"context"
	db "kairon/adapters/database"
	"kairon/cmd/api/infrastructure"
	"kairon/domain/model"
	"kairon/utils"
)

var orderIndex string = "Order"

type OrderRepository interface {
	RunTransaction(ctx context.Context, f func(tx db.DBTransaction) error) error
	Read(id string) (model.Order, error)
	Create(cm model.Order) (model.Order, error)
	Update(id string, changes map[string]any) (model.Order, error)
	Delete(id string) error
	List(queryOpts infrastructure.QueryOpts) ([]model.Order, error)
	Index() string
}

type OrderRepositoryImp struct {
	DB *db.Connection
}

func NewOrderRepository(dbConn *db.Connection) OrderRepository {
	return &OrderRepositoryImp{
		DB: dbConn,
	}
}

func (cs *OrderRepositoryImp) RunTransaction(ctx context.Context, f func(tx db.DBTransaction) error) error {
	return cs.DB.RunTransaction(ctx, f)
}

func (cs *OrderRepositoryImp) Index() string {
	return orderIndex
}

func (cs *OrderRepositoryImp) Read(id string) (model.Order, error) {
	order := model.Order{}
	resMap, err := cs.DB.Read(orderIndex, id, model.Order{})
	if err != nil {
		return order, err
	}

	err = utils.Map2Struct(resMap, &order)
	return order, err
}

func (cs *OrderRepositoryImp) Create(cm model.Order) (model.Order, error) {
	order := model.Order{}
	resMap, err := cs.DB.Create(orderIndex, cm)
	if err != nil {
		return order, err
	}

	err = utils.Map2Struct(resMap, &order)
	return order, err
}

func (cs *OrderRepositoryImp) Update(id string, changes map[string]any) (model.Order, error) {
	order := model.Order{}
	resMap, err := cs.DB.Update(orderIndex, id, model.Order{}, changes)
	if err != nil {
		return order, err
	}

	err = utils.Map2Struct(resMap, &order)
	return order, err
}

func (cs *OrderRepositoryImp) Delete(id string) error {
	return cs.DB.Delete(orderIndex, id, model.Order{})
}

func (cs *OrderRepositoryImp) List(queryOpts infrastructure.QueryOpts) ([]model.Order, error) {
	orders := []model.Order{}
	res, err := cs.DB.List(orderIndex, model.Order{}, queryOpts)
	if err != nil {
		return nil, err
	}

	for _, v := range res {
		order := model.Order{}
		err = utils.Map2Struct(v, &order)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}
