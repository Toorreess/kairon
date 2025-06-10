package usecases

import (
	"context"
	"fmt"
	db "kairon/adapters/database"
	"kairon/cmd/api/infrastructure"
	"kairon/domain/model"
	"kairon/repositories"
	"kairon/utils"
	"time"
)

type OrderUsecase interface {
	Read(id string) (model.Order, error)
	Create(cm model.Order) (model.Order, error)
	Delete(id string) error
	List(queryOpts infrastructure.QueryOpts) ([]model.Order, error)
	Pay(id string) (model.Order, error)
	Cancel(id string) (model.Order, error)
}

type OrderUsecaseImp struct {
	orderRepository   repositories.OrderRepository
	productRepository repositories.ProductRepository
}

func NewOrderUsecase(dr repositories.OrderRepository, pr repositories.ProductRepository) OrderUsecase {
	return &OrderUsecaseImp{
		orderRepository:   dr,
		productRepository: pr,
	}
}

func (cu *OrderUsecaseImp) Read(id string) (model.Order, error) {
	return cu.orderRepository.Read(id)
}

func (cu *OrderUsecaseImp) Create(cm model.Order) (model.Order, error) {
	var order model.Order

	txErr := cu.orderRepository.RunTransaction(context.Background(), func(tx db.DBTransaction) error {
		cm.CreatedAt = time.Now().Unix()
		cm.Status = "pending"

		for _, product := range cm.SelectedProducts {
			// Get current product data
			productData, err := tx.Get(cu.productRepository.Index(), product.ID, model.Product{})
			if err != nil {
				return fmt.Errorf("error getting product %s: %v", product.ID, err)
			}

			var currentProduct model.Product
			if err := utils.Map2Struct(productData, &currentProduct); err != nil {
				return fmt.Errorf("error parsing product data: %v", err)
			}

			// Check if product has infinite stock
			if !currentProduct.InfiniteStock {
				// Check if enough stock is available
				if currentProduct.Stock < product.Quantity {
					return fmt.Errorf("not enough stock for product %s", product.ID)
				}

				// Update stock
				updates := map[string]any{
					"stock": currentProduct.Stock - product.Quantity,
				}

				if err := tx.Update(cu.productRepository.Index(), product.ID, model.Product{}, updates); err != nil {
					return fmt.Errorf("error updating product stock: %v", err)
				}
			}
		}

		// Create the order
		orderMap, err := tx.Create(cu.orderRepository.Index(), cm)
		if err != nil {
			return fmt.Errorf("error creating order: %v", err)
		}

		return utils.Map2Struct(orderMap, &order)
	})

	if txErr != nil {
		return model.Order{}, txErr
	}

	return order, nil
}

func (cu *OrderUsecaseImp) Delete(id string) error {
	return cu.orderRepository.Delete(id)
}

func (cu *OrderUsecaseImp) List(queryOpts infrastructure.QueryOpts) ([]model.Order, error) {
	return cu.orderRepository.List(queryOpts)
}

func (cu *OrderUsecaseImp) Pay(id string) (model.Order, error) {
	order, err := cu.orderRepository.Read(id)
	if err != nil {
		return model.Order{}, err
	}

	if order.Status != "pending" {
		return model.Order{}, fmt.Errorf("not valid status: %s", id)
	}

	changes := map[string]any{
		"status": "paid",
	}
	return cu.orderRepository.Update(id, changes)
}

func (cu *OrderUsecaseImp) Cancel(id string) (model.Order, error) {
	order, err := cu.orderRepository.Read(id)
	if err != nil {
		return model.Order{}, err
	}

	if order.Status != "pending" {
		return model.Order{}, fmt.Errorf("not valid status: %s", id)
	}

	changes := map[string]any{
		"status": "cancelled",
	}
	return cu.orderRepository.Update(id, changes)
}
