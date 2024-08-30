package services

import (
	"testing"

	"go-user/models"
)

func TestGetAllUsers(t *testing.T) {
	t.Parallel()

	t.Run("returns an empty list of users", func(t *testing.T) {
		// given
		users := GetAllUsers()

		// when
		if len(users) != 0 {
			t.Errorf("Expected an empty list of users, got %d", len(users))
		}
	})

	t.Run("returns all users", func(t *testing.T) {
		// given 2 users where created
		users := []models.User{
			{Name: "Mathew", Age: 30},
			{Name: "Jane", Age: 25},
		}
		for _, user := range users {
			CreateUser(user)
		}

		// when
		result := GetAllUsers()

		// then all users are returned
		if len(result) != len(users) {
			t.Errorf("Expected %d users, got %d", len(users), len(result))
		}
	})

	t.Run("returns a single user", func(t *testing.T) {
		// given
		user := models.User{Name: "John", Age: 30}
		var newUser models.User = CreateUser(user)

		// when
		result, _ := GetUserByID(newUser.ID)

		// then
		if result.Name != "John" {
			t.Errorf("Expected user name 'John', got %s", result.Name)
		}

		if result.Age != 30 {
			t.Errorf("Expected user age 30, got %d", result.Age)
		}
	})

	t.Run("updates a user", func(t *testing.T) {

		// given
		user := models.User{Name: "John", Age: 30}
		var newUser models.User = CreateUser(user)

		// when
		newUser.Name = "Jane"
		newUser.Age = 25
		updatedUser, _ := UpdateUser(newUser)

		// then
		if updatedUser.Name != "Jane" {
			t.Errorf("Expected user name 'Jane', got %s", updatedUser.Name)
		}
	})

	// cant update unknown user
	t.Run("updates an unknown user", func(t *testing.T) {

		// given an unknown user
		unknownUser := models.User{ID: -999, Name: "John", Age: 30}

		// when
		updatedUser, error := UpdateUser(unknownUser)

		// then
		if error == nil {
			t.Errorf("Expected an error, got nil")
		}
		if updatedUser != (models.User{}) {
			t.Errorf("Expected void user, got %v", updatedUser)
		}
	})
}
