// Package handlers implements application route handlers
package handlers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alaleks/geospace/internal/server/app/authentication"
	"github.com/alaleks/geospace/internal/server/database"
	"github.com/gofiber/fiber/v2"
)

// typical errors
var (
	ErrMissingRequiredField  = errors.New("missing are no required field")
	ErrUserNotExists         = errors.New("user with current email does not exist")
	ErrInvalidPassword       = errors.New("password is invalid")
	ErrInvalidAuthentication = errors.New("permission denied, user are not authorization")
)

// messages
var (
	MsgCreateUser = fmt.Sprintf("successfully registration")
	MsgAuth       = fmt.Sprintf("successfully authentication")
	MsgLogout     = fmt.Sprintf("successfully exiting")
)

// Hdls represents the handlers and includes db instance.
type Hdls struct {
	db   *database.DB
	auth *authentication.Auth
}

// New creates a new pointer Hdls instance.
func New(db *database.DB, auth *authentication.Auth) *Hdls {
	return &Hdls{
		db:   db,
		auth: auth,
	}
}

// Register provides registration user.
func (h *Hdls) Register(c *fiber.Ctx) error {
	var user struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.BodyParser(&user); err != nil {
		return h.errorBadRequest(c, err)
	}

	if user.Email == "" ||
		user.Name == "" ||
		user.Password == "" ||
		!strings.Contains(user.Email, "@") {
		return h.errorBadRequest(c, ErrMissingRequiredField)
	}

	uid, err := h.db.CreateUser(user.Name, user.Email, h.auth.EncryptPass(user.Password))
	if err != nil {
		return h.errorBadRequest(c, err)
	}

	token, err := h.auth.GetTokenJWT(uid)
	if err != nil {
		return h.errorBadRequest(c, err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(authentication.Expiration),
		Secure:   false,
		HTTPOnly: true,
	})

	return h.sendOK(c, MsgCreateUser)
}

// Login provides authentification user.
func (h *Hdls) Login(c *fiber.Ctx) error {
	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&user); err != nil {
		return h.errorBadRequest(c, err)
	}

	if user.Email == "" ||
		user.Password == "" ||
		!strings.Contains(user.Email, "@") {
		return h.errorBadRequest(c, ErrMissingRequiredField)
	}

	userDB, err := h.db.GetUser(user.Email)
	if err != nil {
		return h.errorBadRequest(c, ErrUserNotExists)
	}

	if !h.auth.CheckPass(user.Password, userDB.Password) {
		return h.errorBadRequest(c, ErrInvalidPassword)
	}

	token, err := h.auth.GetTokenJWT(userDB.UID)
	if err != nil {
		return h.errorBadRequest(c, err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(authentication.Expiration),
		Secure:   false,
		HTTPOnly: true,
	})

	return h.sendOK(c, MsgAuth)
}

// Logout performs exit user.
func (h *Hdls) Logout(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: expired,
	})
	return h.sendOK(c, MsgLogout)
}

// CheckAuthentication checks token validity.
// Token can be provided in Cookie access_token
// or in Header Authorization as Bearer token.
func (h *Hdls) CheckAuthentication(c *fiber.Ctx) error {
	var token string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		token = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("access_token") != "" {
		token = c.Cookies("access_token")
	}

	if strings.TrimSpace(token) == "" {
		return h.errorAuth(c, ErrInvalidAuthentication)
	}

	return c.Next()
}

// errorBadRequest performs send status code 400 and error.
func (h *Hdls) errorBadRequest(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).SendString(err.Error())
}

// errorAuth performs send status code 401 and error.
func (h *Hdls) errorAuth(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
}

// sendOK performs send status 200 and message.
func (h *Hdls) sendOK(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusOK).SendString(msg)
}
