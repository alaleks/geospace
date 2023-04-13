package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/alaleks/geospace/internal/server/database/models"
	"github.com/alaleks/geospace/pkg/distance"
	"github.com/gofiber/fiber/v2"
)

// RespCity represents a data for response list near a citу,
// supplementing the data with the distance field.
type RespCity struct {
	models.City
	Distance int `json:"distance"`
}

// RespDistance represents a data for response
// for requesting calculation of distance.
type RespDistance struct {
	Departure        models.City `json:"departure"`
	Destination      models.City `json:"destination"`
	DistanceStraight int         `json:"distance_straight"`
	DistanceRoad     int         `json:"distance_road,omitempty"`
}

// CalculateDistanceAPI performs a distance between two cities.
func (h *Hdls) CalculateDistanceAPI(c *fiber.Ctx) error {
	departure := strings.TrimSpace(c.Query("departure"))
	destination := strings.TrimSpace(c.Query("destination"))

	if departure == "" {
		err := fmt.Errorf("departure %v", ErrEmptyParam)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	if destination == "" {
		err := fmt.Errorf("destination %v", ErrEmptyParam)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	var (
		response          RespDistance
		chErr             = make(chan error, 1)
		cityDepartureCh   = make(chan models.City, 1)
		cityDestinationCh = make(chan models.City, 1)
	)

	// run concurrently search in database
	go h.db.FindCityConc(departure, chErr, cityDepartureCh)
	go h.db.FindCityConc(destination, chErr, cityDestinationCh)

	// in cycle will be wait when search is done
	// at the error, immediately return the request error
	var n atomic.Int64

	for {
		if n.Load() == 2 {
			break
		}

		select {
		case err := <-chErr:
			err = fmt.Errorf("error find city in db: %v", err)
			return h.errorApiRequest(c, fiber.StatusBadRequest, err)
		case response.Departure = <-cityDepartureCh:
			n.Add(1)
		case response.Destination = <-cityDestinationCh:
			n.Add(1)
		case <-time.After(timeout):
			break
		}
	}

	response.DistanceStraight = int(distance.CalcGreatCircle(
		response.Departure.Latitude, response.Departure.Longitude,
		response.Destination.Latitude, response.Destination.Longitude))

	// here wait no more than 500 ms, if api OSM does not respond or responds long,
	// then return the distance in a straight line
	response.DistanceRoad, _ = h.getDistancebyRoad(
		response.Departure.Longitude, response.Departure.Latitude,
		response.Destination.Longitude, response.Destination.Latitude)

	return c.JSON(response)
}

// FindObjectsNearByNameAP performs search for all objects at a distance
// until n km from the object passed in the query.
func (h *Hdls) FindObjectsNearByNameAPI(c *fiber.Ctx) error {
	var (
		departure  = c.Query("departure")
		distanceTo = c.Query("distanceTo")
	)

	if strings.TrimSpace(departure) == "" {
		err := fmt.Errorf("departure %v", ErrEmptyParam)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	if strings.TrimSpace(distanceTo) == "" {
		err := fmt.Errorf("distanceTo %v", ErrEmptyParam)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	// convert distance to int
	dist, err := strconv.Atoi(distanceTo)
	if err != nil {
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	ciyDeparture, cities, err := h.db.FindObjectsNearByName(departure, dist)
	if err != nil {
		err = fmt.Errorf("error finding cities nearby in database: %v", err)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	respCities := make([]RespCity, 0, len(cities))
	for _, city := range cities {
		respCities = append(respCities, RespCity{
			city,
			int(distance.CalcGreatCircle(ciyDeparture.Latitude, ciyDeparture.Longitude,
				city.Latitude, city.Longitude)),
		})
	}

	response := struct {
		CitiesNearby []RespCity  `json:"cities_nearby"`
		Departure    models.City `json:"departure"`
		DistanceTo   int         `json:"distance_to"`
		QtyNearby    int         `json:"qty_nearby"`
	}{
		Departure:    ciyDeparture,
		DistanceTo:   dist,
		QtyNearby:    len(respCities),
		CitiesNearby: respCities,
	}

	return c.JSON(response)
}

// FindObjectsNearByCoordAPI performs search for all objects at a distance
// until n km from coordinates passed in the query.
func (h *Hdls) FindObjectsNearByCoordAPI(c *fiber.Ctx) error {
	switch {
	case strings.TrimSpace(c.Query("distanceTo")) == "":
		err := fmt.Errorf("distanceTo %v", ErrEmptyParam)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	case strings.TrimSpace(c.Query("lat")) == "":
		err := fmt.Errorf("lat %v", ErrEmptyParam)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	case strings.TrimSpace(c.Query("lon")) == "":
		err := fmt.Errorf("lon %v", ErrEmptyParam)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	dist, err := strconv.Atoi(c.Query("distanceTo"))
	if err != nil {
		err = fmt.Errorf("error convert distance to int: %v", err)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		err = fmt.Errorf("error convert lat to decimal number: %v", err)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	lon, err := strconv.ParseFloat(c.Query("lon"), 64)
	if err != nil {
		err = fmt.Errorf("error convert lon to decimal number: %v", err)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	cities, err := h.db.FindObjectsNearByCoord(lat, lon, dist)
	if err != nil {
		err = fmt.Errorf("error finding cities nearby in database: %v", err)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	respCities := make([]RespCity, 0, len(cities))
	for _, city := range cities {
		respCities = append(respCities, RespCity{
			city,
			int(distance.CalcGreatCircle(lat, lon, city.Latitude, city.Longitude)),
		})
	}

	response := struct {
		CitiesNearby []RespCity `json:"cities_nearby"`
		DistanceTo   int        `json:"distance_to"`
		QtyNearby    int        `json:"qty_nearby"`
	}{
		DistanceTo:   dist,
		QtyNearby:    len(respCities),
		CitiesNearby: respCities,
	}

	return c.JSON(response)
}

// errorApiRequest performs send status code and message error.
func (h *Hdls) errorApiRequest(c *fiber.Ctx, code int, err error) error {
	errReq := struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{
		Code:    code,
		Message: err.Error(),
	}

	return c.Status(code).JSON(errReq)
}
