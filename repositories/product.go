package repositories

import (
	"context"
	db "kairon/adapters/database"
	"kairon/cmd/api/infrastructure"
	"kairon/domain/model"
	"kairon/utils"
)

var productIndex string = "Product"

type ProductRepository interface {
	RunTransaction(ctx context.Context, f func(tx db.DBTransaction) error) error
	Read(id string) (model.Product, error)
	Create(cm model.Product) (model.Product, error)
	Update(id string, changes map[string]any) (model.Product, error)
	Delete(id string) error
	List(queryOpts infrastructure.QueryOpts) ([]model.Product, error)
	Index() string
}

type ProductRepositoryImp struct {
	DB *db.Connection
}

func NewProductRepository(dbConn *db.Connection) ProductRepository {
	return &ProductRepositoryImp{
		DB: dbConn,
	}
}

func (cs *ProductRepositoryImp) RunTransaction(ctx context.Context, f func(tx db.DBTransaction) error) error {
	return cs.DB.RunTransaction(ctx, f)
}

func (cs *ProductRepositoryImp) Index() string {
	return productIndex
}

func (cs *ProductRepositoryImp) Read(id string) (model.Product, error) {
	product := model.Product{}
	resMap, err := cs.DB.Read(productIndex, id, model.Product{})
	if err != nil {
		return product, err
	}

	err = utils.Map2Struct(resMap, &product)
	return product, err
}

func (cs *ProductRepositoryImp) Create(cm model.Product) (model.Product, error) {
	product := model.Product{}
	resMap, err := cs.DB.Create(productIndex, cm)
	if err != nil {
		return product, err
	}

	err = utils.Map2Struct(resMap, &product)
	return product, err
}

func (cs *ProductRepositoryImp) Update(id string, changes map[string]any) (model.Product, error) {
	product := model.Product{}
	resMap, err := cs.DB.Update(productIndex, id, model.Product{}, changes)
	if err != nil {
		return product, err
	}

	err = utils.Map2Struct(resMap, &product)
	return product, err
}

func (cs *ProductRepositoryImp) Delete(id string) error {
	return cs.DB.Delete(productIndex, id, model.Product{})
}

func (cs *ProductRepositoryImp) List(queryOpts infrastructure.QueryOpts) ([]model.Product, error) {
	products := []model.Product{}
	res, err := cs.DB.List(productIndex, model.Product{}, queryOpts)
	if err != nil {
		return nil, err
	}

	for _, v := range res {
		product := model.Product{}
		err = utils.Map2Struct(v, &product)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}
