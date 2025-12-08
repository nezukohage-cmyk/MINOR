package services

import (
	"errors"

	"lexxi/models"
	"lexxi/utils"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func RegisterUser(username, email, phone, password string) (string, error) {
	// check username uniqueness
	count, _ := mgm.Coll(&models.User{}).CountDocuments(mgm.Ctx(), bson.M{"username": username})
	if count > 0 {
		return "", errors.New("username already taken")
	}

	// check email uniqueness if provided
	if email != "" {
		count, _ = mgm.Coll(&models.User{}).CountDocuments(mgm.Ctx(), bson.M{"email": email})
		if count > 0 {
			return "", errors.New("email already in use")
		}
	}

	// check phone uniqueness if provided
	if phone != "" {
		count, _ = mgm.Coll(&models.User{}).CountDocuments(mgm.Ctx(), bson.M{"phone": phone})
		if count > 0 {
			return "", errors.New("phone already in use")
		}
	}

	// hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", errors.New("failed hashing password")
	}

	// create user model
	user := &models.User{
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: hashedPassword,
	}

	err = mgm.Coll(user).Create(user)
	if err != nil {
		return "", errors.New("failed to save user")
	}

	// generate token
	token, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
