package entity

type User struct {
	ID       string `json:"id"`
	Role     string `json:"role"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
