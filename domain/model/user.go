package model

type User struct {
	ID    string `json:"id" firestore:"-"`
	Name  string `json:"name" firestore:"name" validate:"required" updateAllowed:"true"`
	Email string `json:"email" firestore:"email" validate:"required,email" updateAllowed:"true"`
	Role  string `json:"role" firestore:"role" validate:"required,oneof=worker admin" updateAllowed:"true"`

	// Deleted is used for logical deletion
	Deleted bool `json:"-" firestore:"deleted"`
}

func (u *User) UserClaims() map[string]any {
	var claims = make(map[string]any)
	claims["role"] = u.Role

	return claims
}
