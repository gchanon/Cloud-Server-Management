package behavior

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golf/cloudmgmt/services/cloudMgmt/model"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	users         map[string]*model.UserModel
	counterUserId int64
}

func NewUserBehavior() *User {
	return &User{
		users: make(map[string]*model.UserModel),
	}
}

func (user *User) GenSeedUser() error {

	hashPassword, errHashing := bcrypt.GenerateFromPassword([]byte("not-so-secure-password"), bcrypt.DefaultCost)
	if errHashing != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}

	// gen the seed data responding to task 2 req.
	user.users["123123123"] = &model.UserModel{
		UserId:       "123123123",
		Email:        "john.smith@gmail.com",
		PasswordHash: string(hashPassword),
	}

	return nil
}

func (user *User) GetByEmail(email string) (*model.UserModel, error) {
	for _, u := range user.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
}
