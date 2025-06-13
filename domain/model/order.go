package model

type OrderRequest struct {
	Amount           float64           `json:"amount" firestore:"amount" validate:"required" updateAllowed:"true"`
	SelectedProducts []SelectedProduct `json:"products" firestore:"products" validate:"required"`
	MemberID         string            `json:"member_id" firestore:"member_id" validate:"required"`
}

type Order struct {
	ID               string            `json:"id" firestore:"-"`
	Created          int64             `json:"created" firestore:"created"`
	Amount           float64           `json:"amount" firestore:"amount" validate:"required" updateAllowed:"true"`
	SelectedProducts []SelectedProduct `json:"products" firestore:"products" validate:"required"`
	Status           string            `json:"status" firestore:"status" validate:"oneof=pending paid cancelled" updateAllowed:"true"`
	MemberID         string            `json:"member_id" firestore:"member_id" validate:"required"`

	// Deleted is used for logical deletion
	Deleted bool `json:"-" firestore:"deleted"`
}

type SelectedProduct struct {
	ID       string  `json:"id" firestore:"id" validate:"required"`
	Quantity int     `json:"quantity" firestore:"quantity" validate:"required"`
	Price    float64 `json:"price" firestore:"price" validate:"required"`
}
