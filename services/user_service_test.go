package services

import (
	"context"
	"fmt"
	"testing"

	"go-user/models"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

func setupMockRedis() (*miniredis.Miniredis, *redis.Client) {
	// Start a mock Redis server
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	// Create a new Redis client that connects to the mock Redis server
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Replace the global Redis client with the mock client
	rdb = client

	ctx = context.Background()

	return mr, client
}

func TestGetAllUsers(t *testing.T) {
	// Setup the mock Redis server
	mr, _ := setupMockRedis()
	defer mr.Close() // Ensure the mock server is closed after the test

	t.Parallel()

	t.Run("returns an empty list of users", func(t *testing.T) {
		// when
		users, _ := GetAllUsers()

		// then
		if len(users) != 0 {
			t.Errorf("Expected an empty list of users, got %d", len(users))
		}
	})

	t.Run("returns all users", func(t *testing.T) {
		// given 2 users were created
		users := []models.User{
			{Name: "Mathew", Age: 30},
			{Name: "Jane", Age: 25},
		}
		for _, user := range users {
			var u, _ = CreateUser(user)
			fmt.Printf("user %d created\n", u.ID)
		}

		// when
		result, error := GetAllUsers()

		if error != nil {
			t.Errorf("Error getting all users: %v", error)
		}

		// then all users are returned
		if len(result) != len(users) {
			t.Errorf("Expected %d users, got %d", len(users), len(result))
		}
	})
}

func TestGetUserByID(t *testing.T) {
	// Setup the mock Redis server
	mr, _ := setupMockRedis()
	defer mr.Close()

	t.Parallel()

	t.Run("returns a single user", func(t *testing.T) {
		// given
		user := models.User{Name: "John", Age: 30}
		var newUser, _ = CreateUser(user)

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
}

func TestUpdateUser(t *testing.T) {
	// Setup the mock Redis server
	mr, _ := setupMockRedis()
	defer mr.Close()

	t.Parallel()

	t.Run("updates a user", func(t *testing.T) {
		// given
		user := models.User{Name: "John", Age: 30}
		var newUser, _ = CreateUser(user)

		// when
		newUser.Name = "Jane"
		newUser.Age = 25
		updatedUser, _ := UpdateUser(newUser)

		// then
		if updatedUser.Name != "Jane" {
			t.Errorf("Expected user name 'Jane', got %s", updatedUser.Name)
		}
	})

	t.Run("updates an unknown user", func(t *testing.T) {
		// given an unknown user
		unknownUser := models.User{ID: -999, Name: "John", Age: 30}

		// when
		updatedUser, err := UpdateUser(unknownUser)

		// then
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if updatedUser != (models.User{}) {
			t.Errorf("Expected void user, got %v", updatedUser)
		}
	})
}
