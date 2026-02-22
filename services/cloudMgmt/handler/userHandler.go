package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golf/cloudmgmt/services/cloudMgmt/behavior"
	"github.com/golf/cloudmgmt/services/cloudMgmt/handler/middleware"
	externalrepo "github.com/golf/cloudmgmt/services/cloudMgmt/repo/externalRepo"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	behavior     *behavior.User
	jwtSecret    string
	jwtExpire    int
	cookieDomain string
}

func NewUserHandler(behavior *behavior.User, jwtSecret string, jwtExpire int, cookieDomain string) *UserHandler {
	return &UserHandler{
		behavior:     behavior,
		jwtSecret:    jwtSecret,
		jwtExpire:    jwtExpire,
		cookieDomain: cookieDomain,
	}
}

func (handler *UserHandler) Login(c *fiber.Ctx) error {
	var req externalrepo.LoginRequest
	if errParseBody := c.BodyParser(&req); errParseBody != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "email and password are required")
	}

	user, errGetEmail := handler.behavior.GetByEmail(req.Email)
	if errGetEmail != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid email or password")
	}

	token, err := middleware.GenerateToken(user.UserId, user.Email, handler.jwtSecret, handler.jwtExpire)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate token")
	}

	middleware.SetTokenCookie(c, token, handler.jwtExpire, handler.cookieDomain)

	return c.Status(fiber.StatusOK).JSON(externalrepo.LoginResponse{
		LoginStatus: "success",
	})
}
