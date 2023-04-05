// Package handlers implements application route handlers
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alaleks/geospace/internal/server/app/authentication"
	"github.com/alaleks/geospace/internal/server/database"
	"github.com/alaleks/geospace/internal/server/database/models"
	"github.com/alaleks/geospace/pkg/distance"
	"github.com/gofiber/fiber/v2"
)

// typical errors
var (
	ErrUserNotExists                 = errors.New("user with current email does not exist")
	ErrInvalidPassword               = errors.New("password is invalid")
	ErrInvalidAuthentication         = errors.New("permission denied, user are not authorization")
	ErrEmptyDataForCalculateDisnance = errors.New("departure and destination cannot be empty")
	ErrFindCity                      = errors.New("city in not found")
	ErrNotAvailable                  = errors.New("no access to service")
)

// messages
var (
	MsgLogout = "successfully exiting"
	MsgPing   = "all systems work properly :-)"
)

// Hdls represents the handlers and includes db instance.
type Hdls struct {
	db    *database.DB
	auth  *authentication.Auth
	agent *fiber.Agent
}

// New creates a new pointer Hdls instance.
func New(db *database.DB, auth *authentication.Auth) *Hdls {
	return &Hdls{
		db:    db,
		auth:  auth,
		agent: fiber.AcquireAgent(),
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

	var response = struct {
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

	var response = struct {
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
	return h.sendOK(c, MsgLogout)
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

// CalculateDistance performs a distance between two cities.
func (h *Hdls) CalculateDistance(c *fiber.Ctx) error {
	var (
		departure   = c.Query("departure")
		destination = c.Query("destination")
		msg         string
	)

	if strings.TrimSpace(departure) == "" ||
		strings.TrimSpace(destination) == "" {
		return h.errorAuth(c, ErrEmptyDataForCalculateDisnance)
	}

	var (
		cityDeparture   models.City
		cityDestination models.City
		err             error
	)

	cityDeparture, err = h.db.FindCity(departure)
	if err != nil {
		return h.errorBadRequest(c,
			fmt.Errorf("%s : %v", departure, ErrFindCity))
	}

	cityDestination, err = h.db.FindCity(destination)
	if err != nil {
		return h.errorBadRequest(c,
			fmt.Errorf("%s : %v", destination, ErrFindCity))
	}

	distStraight := distance.CalcGreatCircle(cityDeparture.Latitude, cityDeparture.Longitude,
		cityDestination.Latitude, cityDestination.Longitude)

	distanceRoad, err := h.getDistancebyRoad(cityDeparture.Longitude, cityDeparture.Latitude,
		cityDestination.Longitude, cityDestination.Latitude)

	if err != nil || distanceRoad == 0 {
		msg = fmt.Sprintf("distance between %s, %s and %s, %s in a straight line %d km",
			cityDeparture.Name, cityDeparture.Country,
			cityDestination.Name, cityDestination.Country, int(distStraight))

		return h.sendOK(c, msg)
	}

	msg = fmt.Sprintf(`distance between %s, %s and %s, %s:
	- in a straight line %d km
	- by road %d km`,
		cityDeparture.Name, cityDeparture.Country,
		cityDestination.Name, cityDestination.Country,
		int(distStraight), distanceRoad)

	return h.sendOK(c, msg)
}

// Ping performs check work server.
func (h *Hdls) Ping(c *fiber.Ctx) error {
	if err := h.db.SQLX.Ping(); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Errorf("database is down: %v", err).Error())
	}

	return h.sendOK(c, MsgPing)
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

// getDistancebyRoad getting distance between two points by road using api OpenStreetMap.
func (h *Hdls) getDistancebyRoad(lon1, lat1, lon2, lat2 float64) (int, error) {
	url := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=false",
		lon1, lat1, lon2, lat2)

	req := h.agent.Request()
	req.SetTimeout(500 * time.Millisecond)
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(url)

	if err := h.agent.Parse(); err != nil {
		return 0, err
	}

	code, body, _ := h.agent.Bytes()
	if code != 200 {
		return 0, ErrNotAvailable
	}

	var response struct {
		Code   string `json:"code"`
		Routes []struct {
			Distance float64 `json:"distance"`
		}
	}

	err := json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	if response.Code != "Ok" {
		return 0, err
	}

	distance := int(response.Routes[0].Distance / 1000)

	return distance, nil
}
