package handlers

import (
	"strings"
	"sync"

	"github.com/alaleks/geospace/internal/server/database/models"
	"github.com/alaleks/geospace/pkg/distance"
	"github.com/gofiber/fiber/v2"
)

// CalculateDistanceAPI performs a distance between two cities.
func (h *Hdls) CalculateDistanceAPI(c *fiber.Ctx) error {
	var (
		departure   = c.Query("departure")
		destination = c.Query("destination")
	)

	if strings.TrimSpace(departure) == "" ||
		strings.TrimSpace(destination) == "" {
		return h.errorApiRequest(c, fiber.StatusBadRequest,
			ErrEmptyDataForCalculateDisnance)
	}

	var (
		response struct {
			Departure        models.City `json:"departure"`
			Destination      models.City `json:"destination"`
			DistanceStraight int         `json:"distance_straight"`
			DistanceRoad     int         `json:"distance_road,omitempty"`
		}
		err error
		wg  sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		response.Departure, err = h.db.FindCity(departure)
		if err != nil {
			return
		}
	}()

	go func() {
		defer wg.Done()
		response.Destination, err = h.db.FindCity(destination)
		if err != nil {
			return
		}
	}()

	wg.Wait()

	response.DistanceStraight = int(distance.CalcGreatCircle(
		response.Departure.Latitude, response.Departure.Longitude,
		response.Destination.Latitude, response.Destination.Longitude))

	response.DistanceRoad, _ = h.getDistancebyRoad(
		response.Departure.Longitude, response.Departure.Latitude,
		response.Destination.Longitude, response.Destination.Latitude)

	return c.JSON(response)
}

// errorApiRequest performs send status code and message error.
func (h *Hdls) errorApiRequest(c *fiber.Ctx, code int, err error) error {
	var errReq = struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    code,
		Message: err.Error(),
	}

	return c.Status(code).JSON(errReq)
}
