package model

type Product struct {
	ID            string  `json:"id" firestore:"-"`
	Name          string  `json:"name" firestore:"name" validate:"required" updateAllowed:"true"`
	Price         float64 `json:"price" firestore:"price" updateAllowed:"true"`
	Available     bool    `json:"available" firestore:"available" updateAllowed:"true"`
	Stock         int     `json:"stock" firestore:"stock" updateAllowed:"true"`
	InfiniteStock bool    `json:"infinite_stock" firestore:"infinite_stock" updateAllowed:"true"`

	// Deleted is used for logical deletion
	Deleted bool `json:"-" firestore:"deleted"`
}
