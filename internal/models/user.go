package models

type User struct {
	ID         int    `json:"id" example:"1"`
	UserName   string `json:"username" example:"admin"`
	FirstName  string `json:"firstName" example:"John"`
	LastName   string `json:"lastName" example:"Wick"`
	Email      string `json:"email" example:"wick@continental.com"`
	Password   string `json:"password" example:"admin"`
	Phone      string `json:"phone" example:"8-999-666-99-66"`
	UserStatus int    `json:"userStatus" example:"1"`
}
