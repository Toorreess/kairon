package usecases

import (
	"kairon/cmd/api/infrastructure"
	"kairon/domain/model"
	"kairon/repositories"
)

type ProductUsecase interface {
	Read(id string) (model.Product, error)
	Create(cm model.Product) (model.Product, error)
	Update(id string, changes map[string]any) (model.Product, error)
	Delete(id string) error
	List(queryOpts infrastructure.QueryOpts) ([]model.Product, error)
}

type ProductUsecaseImp struct {
	productRepository repositories.ProductRepository
}

func NewProductUsecase(dr repositories.ProductRepository) ProductUsecase {
	return &ProductUsecaseImp{
		productRepository: dr,
	}
}

func (cu *ProductUsecaseImp) Read(id string) (model.Product, error) {
	return cu.productRepository.Read(id)
}

func (cu *ProductUsecaseImp) Create(cm model.Product) (model.Product, error) {
	return cu.productRepository.Create(cm)
}

func (cu *ProductUsecaseImp) Update(id string, changes map[string]any) (model.Product, error) {
	return cu.productRepository.Update(id, changes)
}

func (cu *ProductUsecaseImp) Delete(id string) error {
	return cu.productRepository.Delete(id)
}

func (cu *ProductUsecaseImp) List(queryOpts infrastructure.QueryOpts) ([]model.Product, error) {
	return cu.productRepository.List(queryOpts)
}
