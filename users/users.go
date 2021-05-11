package users

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type Users struct {
	db *gorm.DB
}

func New(db *gorm.DB) Users {
	return Users{db}
}

func (u Users) Find() []User {
	var users []User
	result := u.db.Find(&users)
	fmt.Println("RowsAffected:", result.RowsAffected, "ERROR", result.Error)
	return users
}
