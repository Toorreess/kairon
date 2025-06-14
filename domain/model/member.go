package model

type Member struct {
	ID           string `json:"id" firestore:"-"`
	Name         string `json:"name" firestore:"name" validate:"required" updateAllowed:"true"`
	Identifier   string `json:"identifier" firesotre:"identifier" validate:"required" updateAllowed:"true"`
	Email        string `json:"email" firestore:"email" validate:"required,email" updateAllowed:"true"`
	Phone        string `json:"phone" firestore:"phone" updateAllowed:"true"`
	Status       string `json:"status" firestore:"status" validate:"oneof=active inactive" updateAllowed:"true"`
	MembershipID string `json:"membership_id" firestore:"membership_id" validate:"required" updateAllowed:"true"`

	ActivityList []string `json:"activity_list" firestore:"activity_list" updateAllowed:"true"`
	// Deleted is used for logical deletion
	Deleted bool `json:"-" firestore:"deleted"`
}
