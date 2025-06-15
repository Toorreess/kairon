package model

type Activity struct {
	ID          string `json:"id" firestore:"-"`
	Name        string `json:"name" firestore:"name" validate:"required" updateAllowed:"true"`
	Duration    int    `json:"duration" firestore:"duration" updateAllowed:"true"`
	MaxCapacity int    `json:"max_capacity" firestore:"max_capacity"`
	IsActive    bool   `json:"is_active" firestore:"is_active" updateAllowed:"true"`

	// Deleted is used for logical deletion
	Deleted bool `json:"-" firestore:"deleted"`
}

type ActivityReserveRequest struct {
	MemberID     string   `json:"member_id" validate:"required"`
	ActivityList []string `json:"activity_list" validate:"required"`
}
