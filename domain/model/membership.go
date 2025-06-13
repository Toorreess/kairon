package model

type Membership struct {
	ID        string  `json:"id" firestore:"-"`
	Name      string  `json:"name" firestore:"name" validate:"required" updateAllowed:"true"`
	Price     float64 `json:"price" firestore:"price" updateAllowed:"true"`
	Available bool    `json:"available" firestore:"available" updateAllowed:"true"`

	// Deleted is used for logical deletion
	Deleted bool `json:"-" firestore:"deleted"`
}
