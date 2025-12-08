package services

import (
	"errors"

	"lexxi/models"
	"lexxi/utils"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func Login(identifier, password string) (string, error) {
	var user models.User

	// Try matching username, email, or phone
	filter := bson.M{
		"$or": []bson.M{
			{"username": identifier},
			{"email": identifier},
			{"phone": identifier},
		},
	}

	err := mgm.Coll(&models.User{}).First(filter, &user)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Compare password
	if !utils.CheckPassword(user.Password, password) {
		return "", errors.New("invalid password")
	}

	// Generate new token
	token, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
