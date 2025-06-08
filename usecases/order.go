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
	Update(id string, changes map[string]any) (model.Order, error)
	Delete(id string) error
	List(queryOpts infrastructure.QueryOpts) ([]model.Order, error)
}

type OrderUsecaseImp struct {
	orderRepository repositories.OrderRepository
}

func NewOrderUsecase(dr repositories.OrderRepository) OrderUsecase {
	return &OrderUsecaseImp{
		orderRepository: dr,
	}
}

func (cu *OrderUsecaseImp) Read(id string) (model.Order, error) {
	return cu.orderRepository.Read(id)
}

func (cu *OrderUsecaseImp) Create(cm model.Order) (model.Order, error) {
	var order model.Order

	txErr := cu.orderRepository.RunTransaction(context.Background(), func(tx db.DBTransaction) error {
		cm.CreatedAt = time.Now().Unix()

		for _, product := range cm.SelectedProducts {
			// Get current product data
			productData, err := tx.Get("Product", product.ID, model.Product{})
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

				if err := tx.Update("Product", product.ID, model.Product{}, updates); err != nil {
					return fmt.Errorf("error updating product stock: %v", err)
				}
			}
		}

		// Create the order
		orderMap, err := tx.Create("Order", cm)
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

// Needed?
func (cu *OrderUsecaseImp) Update(id string, changes map[string]any) (model.Order, error) {
	return cu.orderRepository.Update(id, changes)
}

func (cu *OrderUsecaseImp) Delete(id string) error {
	return cu.orderRepository.Delete(id)
}

func (cu *OrderUsecaseImp) List(queryOpts infrastructure.QueryOpts) ([]model.Order, error) {
	return cu.orderRepository.List(queryOpts)
}
