package model

type Order struct {
	ID               string            `json:"id" firestore:"-"`
	CreatedAt        int64             `json:"-" firestore:"created"`
	Reference        string            `json:"reference" firestore:"reference" validate:"required"`
	Amount           float64           `json:"amount" firestore:"amount" validate:"required" updateAllowed:"true"`
	SelectedProducts []SelectedProduct `json:"products" firestore:"products" validate:"required"`
	Status           string            `json:"status" firestore:"status" validate:"required,oneof=pending paid cancelled" updateAllowed:"true"`

	// Deleted is used for logical deletion
	Deleted bool `json:"-" firestore:"deleted"`
}

type SelectedProduct struct {
	ID       string  `json:"id" firestore:"id" validate:"required"`
	Quantity int     `json:"quantity" firestore:"quantity" validate:"required"`
	Price    float64 `json:"price" firestore:"price" validate:"required"`
	Tax      string  `json:"tax" firestore:"tax"`
}
