package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) Validate() error {
	// Add validation logic for user fields if necessary
	return nil
}