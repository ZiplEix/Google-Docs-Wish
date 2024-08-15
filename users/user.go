package users

import (
	"context"
	"fmt"

	"github.com/ZiplEix/Google-Docs-Wish/database"
)

type User struct {
	Email    string
	Password string
	ID       string
}

func New(data map[string]interface{}, id ...string) *User {
	user := &User{
		Email:    "",
		Password: "",
		ID:       "",
	}

	if email, ok := data["email"].(string); ok {
		user.Email = email
	}

	if password, ok := data["password"].(string); ok {
		user.Password = password
	}

	if len(id) > 0 {
		user.ID = id[0]
	}

	return user
}

func (user *User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"email":    user.Email,
		"password": user.Password,
	}
}

func (user *User) CreateInDb() (string, error) {
	docRef, wr, err := database.FirestoreClient.Collection("users").Add(context.Background(), user.ToMap())
	if err != nil {
		return "", fmt.Errorf("error creating user: %v", err)
	}

	fmt.Printf("User created: %v\n", wr)

	user.ID = docRef.ID

	return docRef.ID, nil
}
