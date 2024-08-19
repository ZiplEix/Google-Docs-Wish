package database

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	Email    string
	Password string
	ID       string
}

func NewUser(data map[string]interface{}, id ...string) *User {
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
		"id":       user.ID,
	}
}

func GetUserFromId(id string) (*User, error) {
	state, err := FirestoreClient.Doc("users/" + id).Get(context.Background())
	if err != nil {
		return nil, err
	}

	user := NewUser(state.Data(), state.Ref.ID)

	return user, nil
}

func (user *User) CreateInDb() (string, error) {
	docRef, wr, err := FirestoreClient.Collection("users").Add(context.Background(), user.ToMap())
	if err != nil {
		return "", fmt.Errorf("error creating user: %v", err)
	}

	fmt.Printf("User created: %v\n", wr)

	user.ID = docRef.ID

	_, err = docRef.Set(context.Background(), user.ToMap())
	if err != nil {
		return "", fmt.Errorf("error updating user: %v", err)
	}

	return docRef.ID, nil
}

func GetUserFromCookie(c *fiber.Ctx) (User, error) {
	userID := c.Locals("userID").(string)

	state, err := FirestoreClient.Doc("users/" + userID).Get(c.Context())
	if err != nil {
		return User{}, err
	}

	user := NewUser(state.Data(), state.Ref.ID)

	return *user, nil
}
