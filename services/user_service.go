package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go-user/models"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Cambia la dirección si es necesario
	})
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User

	var cursor uint64
	iter := rdb.Scan(ctx, cursor, "user:*", 100).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		if key == "user:id" {
			continue
		}
		userJSON, err := rdb.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		var user models.User
		if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func CreateUser(user models.User) (models.User, error) {
	user.ID = int(rdb.Incr(ctx, "user:id").Val())

	userKey := fmt.Sprintf("user:%d", user.ID)
	userJSON, err := json.Marshal(user)
	if err != nil {
		return models.User{}, err
	}

	if err := rdb.Set(ctx, userKey, userJSON, 0).Err(); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func GetUserByID(id int) (models.User, error) {
	userKey := fmt.Sprintf("user:%d", id)
	userJSON, err := rdb.Get(ctx, userKey).Result()
	if err == redis.Nil {
		return models.User{}, fmt.Errorf("user not found")
	} else if err != nil {
		return models.User{}, err
	}

	var user models.User
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func UpdateUser(newUser models.User) (models.User, error) {
	// actualiza el usuario en redis utilizando la clave newUser.ID
	userJSON, err := json.Marshal(newUser)
	if err != nil {
		return models.User{}, err
	}
	if ok, _ := rdb.SetXX(ctx, fmt.Sprintf("user:%d", newUser.ID), userJSON, 0).Result(); ok {
		return newUser, nil
	} else {
		return models.User{}, fmt.Errorf("user not found")
	}
}
