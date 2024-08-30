package services

import (
	"errors"
	"sync"

	"go-user/models"
)

var (
	users  = make(map[int]models.User)
	nextID = 1
	mutex  sync.Mutex
)

func GetAllUsers() []models.User {
	mutex.Lock()
	defer mutex.Unlock()

	var userList []models.User
	for _, user := range users {
		userList = append(userList, user)
	}
	return userList
}

func CreateUser(user models.User) models.User {
	mutex.Lock()
	defer mutex.Unlock()

	user.ID = nextID
	users[nextID] = user
	nextID++
	return user
}

// GetUserByID retrieves a user by ID.
//
// The returned error is nil if the user is found, otherwise it is an error
// with the message "user not found".
func GetUserByID(id int) (models.User, error) {
	mutex.Lock()
	defer mutex.Unlock()

	user, exists := users[id]
	if !exists {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func UpdateUser(newUser models.User) (models.User, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := users[newUser.ID]; !exists {
		return models.User{}, errors.New("user not found")
	}

	users[newUser.ID] = newUser

	return newUser, nil
}
