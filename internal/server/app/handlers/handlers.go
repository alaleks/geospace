// Package handlers implements application route handlers
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alaleks/geospace/internal/server/app/authentication"
	"github.com/alaleks/geospace/internal/server/chatgpt"
	"github.com/alaleks/geospace/internal/server/database"
	"github.com/gofiber/fiber/v2"
)

const (
	timeout        = 500 * time.Millisecond
	counterDoneJob = 2
)

// typical errors
var (
	ErrUserNotExists         = errors.New("user with current email does not exist")
	ErrInvalidPassword       = errors.New("password is invalid")
	ErrInvalidAuthentication = errors.New("permission denied, user are not authorization")
	ErrEmptyParam            = errors.New("parameter cannot be empty")
	ErrFindCity              = errors.New("city in not found")
	ErrNotAvailable          = errors.New("no access to service")
	ErrEmptyResults          = errors.New("was get empty results")
)

// messages
var (
	MsgLogout = "successfully exiting"
	MsgPing   = "all systems work properly :-)"
)

// Hdls represents the handlers and includes db instance.
type Hdls struct {
	agent   *fiber.Agent
	auth    *authentication.Auth
	chatGPT *chatgpt.ChatGPT
	db      *database.DB
}

// New creates a new pointer Hdls instance.
func New(db *database.DB, auth *authentication.Auth) *Hdls {
	return &Hdls{
		db:      db,
		auth:    auth,
		agent:   fiber.AcquireAgent(),
		chatGPT: chatgpt.New(),
	}
}

// SignUp provides registration a new user.
func (h *Hdls) SignUp(c *fiber.Ctx) error {
	var user struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.BodyParser(&user); err != nil {
		return h.errorBadRequest(c, err)
	}

	switch {
	case user.Email == "":
		return h.errorBadRequest(c, fmt.Errorf("email cannot be empty"))
	case !strings.Contains(user.Email, "@"):
		return h.errorBadRequest(c, fmt.Errorf("email has invalid format"))
	case user.Password == "":
		return h.errorBadRequest(c, fmt.Errorf("password cannot be empty"))
	}

	uid, err := h.db.CreateUser(user.Name, user.Email, h.auth.EncryptPass(user.Password))
	if err != nil {
		return h.errorBadRequest(c, err)
	}

	token, err := h.auth.GetTokenJWT(uid)
	if err != nil {
		return h.errorBadRequest(c, err)
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	return c.JSON(response)
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

	switch {
	case user.Email == "":
		return h.errorBadRequest(c, fmt.Errorf("email cannot be empty"))
	case !strings.Contains(user.Email, "@"):
		return h.errorBadRequest(c, fmt.Errorf("email has invalid format"))
	case user.Password == "":
		return h.errorBadRequest(c, fmt.Errorf("password cannot be empty"))
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

	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	return c.JSON(response)
}

// GetCountry returns list country with country code.
func (h *Hdls) GetCountry(c *fiber.Ctx) error {
	var countries []string

	err := h.db.SQLX.Select(&countries,
		`SELECT DISTINCT CONCAT(country_code, ": ", country) AS country 
		FROM cities WHERE country <> "" ORDER BY country;`)
	if err != nil {
		return h.errorBadRequest(c, err)
	}

	return c.SendString(strings.Join(countries, ","))
}

// Logout performs exit user.
func (h *Hdls) Logout(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: expired,
	})
	return c.SendString(MsgLogout)
}

// CheckAuthentication checks token validity.
// Token can be provided in Cookie access_token
// or in Header Authorization as Bearer token.
func (h *Hdls) CheckAuthentication(c *fiber.Ctx) error {
	var token string

	if strings.HasPrefix(c.Get("Authorization"), "Bearer ") {
		token = strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	} else if c.Cookies("access_token") != "" {
		token = c.Cookies("access_token")
	}

	if strings.TrimSpace(token) == "" {
		return h.errorAuth(c, ErrInvalidAuthentication)
	}

	if err := h.auth.CheckToken(token); err != nil {
		return h.errorAuth(c, ErrInvalidAuthentication)
	}

	return c.Next()
}

// Ping performs check work server.
func (h *Hdls) Ping(c *fiber.Ctx) error {
	if err := h.db.SQLX.Ping(); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Errorf("database is down: %v", err).Error())
	}

	return c.SendString(MsgPing)
}

// errorBadRequest performs send status code 400 and error.
func (h *Hdls) errorBadRequest(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).SendString(err.Error())
}

// errorAuth performs send status code 401 and error.
func (h *Hdls) errorAuth(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
}

// getDistancebyRoad getting distance between two points by road using api OpenStreetMap.
func (h *Hdls) getDistancebyRoad(lon1, lat1, lon2, lat2 float64) (int, error) {
	url := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=false",
		lon1, lat1, lon2, lat2)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var (
		chErr = make(chan error, 1)
		dist  = make(chan float64, 1)
	)

	go func() {
		req := h.agent.Request()
		req.Header.SetMethod(fiber.MethodGet)
		req.SetRequestURI(url)

		if err := h.agent.Parse(); err != nil {
			chErr <- err
			return
		}

		code, body, _ := h.agent.Bytes()
		if code != 200 {
			chErr <- ErrNotAvailable
			return
		}

		var response struct {
			Code   string `json:"code"`
			Routes []struct {
				Distance float64 `json:"distance"`
			}
		}

		err := json.Unmarshal(body, &response)
		if err != nil {
			chErr <- err
			return
		}

		if response.Code != "Ok" {
			chErr <- err
			return
		}

		if len(response.Routes) == 0 {
			chErr <- ErrEmptyResults
			return
		}

		dist <- response.Routes[0].Distance
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case err := <-chErr:
		return 0, err
	case result := <-dist:
		return int(result / 1000), nil
	}
}
